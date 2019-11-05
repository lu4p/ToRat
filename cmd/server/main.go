package main

import (
	_ "github.com/dimiro1/banner/autoload"
	server "github.com/lu4p/ToRat/torat_server"
)

func main() {
	go server.Start()
	server.Shell()
}
