package server

import (
	"net"
	"time"

	"github.com/jinzhu/gorm"
)

type client struct {
	gorm.Model
	Hostname string
	Name     string
	Path     string
	IP       string
	Location string
	LastConn time.Time
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
	Conn     net.Conn
	Client   *client
}
