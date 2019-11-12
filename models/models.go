package models

import (
	"os"

	"github.com/lu4p/shred"
)

type Void int

type Cmd struct {
	Cmd        string
	Powershell bool
}

type Shred struct {
	Conf shred.Conf
	Path string
}

type File struct {
	Path    string
	Perm    os.FileMode
	Content []byte
}

type Dir struct {
	Path  string
	Files []string
}
