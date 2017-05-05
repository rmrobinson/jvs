package main

import (
	"flag"
	"faltung.ca/jvs/service/slackbotd"
)

var (
	token = flag.String("token", "", "The Slack API token this bot user will use")
	commanderAddr = flag.String("commanderAddr", "127.0.0.1:10004", "The IP:Port of the commander server to interface with")
)

func main() {
	flag.Parse()

	s := slackbotd.Server{}
	s.Run(*token, *commanderAddr)
}
