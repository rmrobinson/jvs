package slackbotd

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"time"
)

type message struct {
	Contents string
	User string
	Source string
	Timestamp string

	// This is kept around so we can give an appropriate reply.
	event *slack.MessageEvent
}


type slackClient struct {
	api   *slack.Client
	rtm   *slack.RTM
	token string

	broadcastChannelID string

	//team     slack.Team
	//users    []slack.User
	//groups   []slack.Group
	//channels []slack.Channel

}

func (c *slackClient) sendMessage(m message) {
	c.rtm.SendMessage(c.rtm.NewOutgoingMessage(m.Contents, c.broadcastChannelID))
}

func (c *slackClient) sendMessageReply(m message) {
	info := c.rtm.GetInfo()

	channel := info.GetChannelByID(m.event.Channel)
	im := info.GetIMByID(m.event.Channel)

	if channel != nil {
		c.rtm.SendMessage(c.rtm.NewOutgoingMessage(m.Contents, channel.ID))
	} else if im != nil {
		c.rtm.SendMessage(c.rtm.NewOutgoingMessage(m.Contents, im.ID))
	} else {
		c.sendMessage(m)
	}
}

func (c *slackClient) run(token string, broadcastChannelName string, messageSink chan message) {
	c.token = token
	c.api = slack.New(c.token)

	logger := log.New(os.Stdout, "slackbotd: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	//c.api.SetDebug(true)

	c.rtm = c.api.NewRTM()

	go c.rtm.ManageConnection()

	for msg := range c.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:

		case *slack.ConnectedEvent:
			fmt.Println("homebot connected")

			info := c.rtm.GetInfo()

			// Find the channel we will use for generic broadcast messages
			// Save the ID then notify that we're alive!
			for _, channel := range info.Channels {
				if channel.Name == broadcastChannelName {
					c.broadcastChannelID = channel.ID

					m := message{
						Contents: "homebot connected",
						User: "homebot",
						Timestamp: time.Now().Format(time.RFC3339),
					}
					c.sendMessage(m)
					break
				}
			}

		case *slack.MessageEvent:
			info := c.rtm.GetInfo()

			user := info.GetUserByID(ev.User)
			channel := info.GetChannelByID(ev.Channel)
			im := info.GetIMByID(ev.Channel)

			m := message{
				Contents: ev.Text,
				Timestamp: ev.Timestamp,
				event: ev,
			}

			if user != nil {
				m.User = user.RealName
			} else if channel != nil {
				m.Source = "Channel"
			} else if im != nil {
				m.Source = "DM"
			}

			if ev.SubType == "message_changed" {
				m.Contents = "Sorry, edits are not supported"
				c.sendMessageReply(m)
				break
			}

			messageSink <- m

		case *slack.PresenceChangeEvent:
			info := c.rtm.GetInfo()
			user := info.GetUserByID(ev.User)
			fmt.Printf("Presence Change: %s is %s\n", user.RealName, ev.Presence)

		case *slack.LatencyReport:
			// TODO: use this for retry/backoff decisions

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Print("Invalid credentials\n")
			return

		default:
			// Ignore others
		}
	}
}
