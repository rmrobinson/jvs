package main

import (
	"flag"
	"fmt"
	"github.com/rmrobinson/jvs/service/device"
)

func main() {
	var (
		dbPath = flag.String("dbPath", "", "The path to the DB to inspect")
		name = flag.String("name", "", "The optional name to save")
	)

	flag.Parse()

	if len(*dbPath) < 1 {
		fmt.Printf("Database path must be specified\n")
		return
	}

	brdb := device.BottlerocketDB{}

	err := brdb.Open(*dbPath)

	if err != nil {
		fmt.Printf("Unable to open DB: %s\n", err.Error())
		return
	}

	bID, idErr := brdb.ID()
	bName, nameErr := brdb.Name()

	if idErr != nil {
		fmt.Printf("Error getting bridge ID: %s\n", idErr.Error())
		return
	} else if nameErr != nil {
		fmt.Printf("Error getting bridge name: %s\n", nameErr.Error())
		return
	}

	fmt.Printf("Bridge %s (%s)\n", bID, bName)

	devices, err := brdb.Devices()

	if err != nil {
		fmt.Printf("Error getting devices: %s\n", err.Error())
		return
	}

	for _, device := range devices {
		fmt.Printf("Device %+v\n", device)
	}

	if len(*name) > 0 {
		err = brdb.SetName(*name)
	}
}
