package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/rmrobinson/hue-go"
	"github.com/rmrobinson/jvs/service/building/bridge"
)

// Full disclosure, the below logic is very gross.
// The hue-go library should support pairing with a bridge,
// that would clean the below logic up significantly.

func main() {
	var (
		dbPath = flag.String("dbPath", "", "The path to the DB to update")
		bridgeIP = flag.String("bridgeIP", "", "The Hue bridge IP to pair with")
	)

	flag.Parse()

	if len(*dbPath) < 1 {
		fmt.Printf("The Hue DB path must be specified\n")
		os.Exit(1)
	}

	db := &bridge.HueDB{}
	err := db.Open(*dbPath)
	if err != nil {
		fmt.Printf("Error opening Hue DB: %s\n", err.Error())
		os.Exit(1)
	}

	if len(*bridgeIP) < 1  {
		fmt.Printf("The bridge IP must be specified\n")
		os.Exit(1)
	}

	bridgeUrl, err := url.Parse("http://" + *bridgeIP + "/description.xml")
	if err != nil {
		fmt.Printf("Unable to parse supplied bridge address: %s\n", err.Error())
		return
	}

	var b hue_go.Bridge

	if err = b.Init(bridgeUrl); err != nil {
		fmt.Printf("Unable to initialize supplied bridge address: %s\n", err.Error())
		return
	}

	fmt.Printf("Please press the 'Pair' button on the Hue bridge. Once pressed, type 'Y' to proceed\n")
	fmt.Print("> ")

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("Unable to read input: %s\n", err.Error())
		return
	}
	if input != "Y\n" {
		fmt.Printf("Non 'Y' character supplied, cancelling the pairing process\n")
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Unable to read hostname: %s\n", err.Error())
		return
	}

	var reqBody struct {
		DeviceType string `json:"devicetype"`
	}
	type respEntry struct {
		Result struct {
			Key string `json:"username"`
		} `json:"success"`
	}
	var respBody []respEntry

	reqBody.DeviceType = fmt.Sprintf("%s#%s\n", "deviced", hostname[0:18])
	reqBytes := new(bytes.Buffer)
	err = json.NewEncoder(reqBytes).Encode(&reqBody)
	if err != nil {
		fmt.Printf("Error encoding request to Hue API: %s\n", err.Error())
		return
	}

	resp, err := http.Post("http://" + *bridgeIP + "/api", "application/json; charset=utf-8", reqBytes)
	if err != nil {
		fmt.Printf("Unable to POST to the Hue API: %s\n", err.Error())
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Printf("Error decoding response from Hue API: %s\n", err.Error())
		return
	}
	if len(respBody) < 1 {
		fmt.Printf("No entries in the results array\n")
		return
	}

	bridgeID := b.Id()

	fmt.Printf("Saving %s for bridge %s\n", respBody[0].Result.Key, bridgeID)
	db.SaveProfile(context.Background(), bridgeID, respBody[0].Result.Key)
}
