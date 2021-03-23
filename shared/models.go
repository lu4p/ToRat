package shared

import (
	"os"
	"reflect"
)

// Make sure Models are never garbled.
var (
	_ = reflect.TypeOf(Void(0))
	_ = reflect.TypeOf(Cmd{})
	_ = reflect.TypeOf(Shred{})
	_ = reflect.TypeOf(File{})
	_ = reflect.TypeOf(Dir{})
	_ = reflect.TypeOf(EncAsym{})
	_ = reflect.TypeOf(Hardware{})
	_ = reflect.TypeOf(OSInfo{})
	_ = reflect.TypeOf(Speedtest{})
	_ = reflect.TypeOf(Gomap{})
)

type Void int

type Cmd struct {
	Cmd        string
	Powershell bool
}

type Shred struct {
	Path   string
	Times  int
	Zeros  bool
	Remove bool
}

type File struct {
	Path    string
	Fpath   string
	Perm    os.FileMode
	Content []byte
}

type Dir struct {
	Path  string
	Files []string
}

type EncAsym struct {
	EncAesKey []byte
	EncData   []byte
}

type Hardware struct {
	Runtime   string `json:"Runtime"`
	OSArch    string `json:"OSArch"`
	OSName    string `json:"OSName"`
	OSVersion string `json:"OSVersion"`
	CPU       string `json:"CPU"`
	Cores     uint32 `json:"Cores"`
	RAM       string `json:"RAM"`
	GPU       string `json:"GPU"`
	Drives    string `json:"Drives"`
}

type OSInfo struct {
	Runtime   string `json:"Runtime"`
	OSArch    string `json:"OSArch"`
	OSName    string `json:"OSName"`
	OSVersion string `json:"OSVersion"`
}

type Speedtest struct {
	IP       string `json:"IP"`
	Ping     string `json:"Ping"`
	Download float64 `json:"Download"`
	Upload   float64 `json:"Upload"`
	Country  string `json:"Country"`
}

type Gomap struct {
	Scan string `json:"Scan"`
}
