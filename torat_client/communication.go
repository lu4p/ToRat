package client

import (
	"encoding/binary"
	"io/ioutil"
	"net"
	"os"
)

const buffsize = 4096

type connection struct {
	Conn    net.Conn
	Sysinfo string
}

func (c *connection) recv() ([]byte, error) {
	var size int64
	err := binary.Read(c.Conn, binary.LittleEndian, &size)
	if err != nil {
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

func (c *connection) recvSt() (string, error) {
	recv, err := c.recv()
	if err != nil {
		return "", err
	}

	return string(recv), nil
}

func (c *connection) send(data []byte) error {
	size := len(data)
	err := binary.Write(c.Conn, binary.LittleEndian, int64(size))
	if err != nil {
		return err

	}
	_, err = c.Conn.Write(data)
	return err

}

func (c *connection) sendSt(cmdout string) error {
	return c.send([]byte(cmdout))
}

func (c *connection) sendFile(fname string) error {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		content = []byte("err")
	}
	err = c.send(content)
	if err != nil {
		return err
	}
	return nil
}

func (c *connection) recvFile(filename string) error {
	data, err := c.recv()
	if err != nil {
		return err
	}
	newFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}
