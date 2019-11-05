package server

import (
	"fmt"
)

func (c *activeClient) runCommand(command string, print bool) (string, error) {
	err := c.sendSt(command)
	if err != nil {
		if print {
			fmt.Println("Error while running command:", err)
		}
		return "", err
	}
	output, err := c.recvSt()
	if err != nil {
		fmt.Println("err")
		return "", err
	}
	if print {
		fmt.Println(output)
	}
	return output, nil
}

func (c *activeClient) runCommandByte(command string) ([]byte, error) {
	err := c.sendSt(command)
	if err != nil {
		return nil, err
	}
	b, err := c.recv()
	if err != nil {
		return nil, err
	}
	return b, nil
}
