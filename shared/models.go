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
	_ = reflect.TypeOf(Speedtest{})
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
	Runtime   string
	OSArch    string
	OSName    string
	OSVersion string
	CPU       string
	Cores     uint32
	RAM       string
	GPU       string
	Drives    string
}

type Speedtest struct {
	IP       string
	Ping     string
	Download float64
	Upload   float64
	Country  string
}
