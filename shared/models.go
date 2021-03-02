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
	_ = reflect.TypeOf(Speedtest{})
	_ = reflect.TypeOf(Nmap{})
	_ = reflect.TypeOf(NmapLocal{})
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

type Speedtest struct {
	IP       string
	Ping     string
	Download float64
	Upload   float64
	Country  string
}

type Nmap struct {
	IP          string
	TimeElapsed float32
	Scan        string
}

type NmapLocal struct {
	Range       string
	TimeElapsed float32
	Hosts       int
	Scan        string
}
