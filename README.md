# <img src="./ToRat_Logo.png" width="180px"> 
[![License](https://img.shields.io/github/license/lu4p/ToRat.svg)](https://unlicense.org/)
[![CircleCI](https://circleci.com/gh/lu4p/ToRat.svg?style=svg)](https://circleci.com/gh/lu4p/ToRat)
[![Go Report Card](https://goreportcard.com/badge/github.com/lu4p/ToRat)](https://goreportcard.com/report/github.com/lu4p/ToRat)

A Cross Platform Remote Administration tool written in Go using Tor as its transport mechanism
currently supporting Windows, Linux, MacOS clients.

## Docker Images Build Status

Name | Status | Pulls
--- | --- | ---
torat-pre | ![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/lu4p/torat-pre) | ![Docker Pulls](https://img.shields.io/docker/pulls/lu4p/torat-pre)
torat | ![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/lu4p/torat) | This image should not be pulled!

## How to
[How to use Docker Image](https://github.com/lu4p/ToRat/wiki/How-to-use-Docker-Image)

## Current Features
- the ToRAT_client communicates over TCP(with TLS) proxied through Tor with the ToRat_server (hidden service)
	- [x] anonymity of client and server
	- [x] end-to-end encryption
- Cross Platform reverse shell (Windows, Linux, Mac OS)
- Windows:
	- Multiple User Account Control Bypasses (Privilege escalation)
	- Multiple Persistence methods (User, Admin)
- Linux:
	- Multiple Persistence methods (User, Admin)
- optional transport without Tor e.g. Use Tor2Web, a DNS Hostname or public/ local IP
	- [x] smaller binary
	- [ ] anonymity of client and server
- embedded Tor
- Unique persistent ID for every client
	- give a client an Alias
	- all Downloads from client get saved to ./$ID/$filename
- sqlite via gorm for storing information about the clients

### Server Shell
- Supports multiple connections
- Welcome Banner
- Colored Output
- Tab-Completion of:
  - Commands
  - Files/ Directories in the working directory of the server

Command | Info
--- | ---
**select** |  Select client to interact with
**list** |  list all connected clients
**alias** |  Select client to give an alias
**cd** |  change the working directory of the server
**exit** | exit the server

#### Shell after selection of a client
- Tab-Completion of:
  - Commands
  - Files/ Directories in the working directory of the client

Command | Info
--- | ---
**cd** | change the working directory of the client
**ls** | list the content of the working directory of the client
**shred** | delete files/ directories unrecoverable
**shredremove** | same as shred + removes the shredded files
**screen** | take a Screenshot of the client
**cat** | view Textfiles from the client including .docx, .rtf, .pdf, .odt
**alias** | give the client a custom alias
**down** | download a file from the client
**up** | upload a file to the client
**escape** | escape a command and run it in a native shell on the client
**reconnect** | tell the client to reconnect
**exit** | background current session an return to main shell
else  | the command will be executed in a native shell on the client

## Upcoming Features
- [ ] Privilege escalation for Linux
- [ ] Persistence and privilege escalation for Mac OS
- [ ] Support for Android and iOS needs fix of https://github.com/ipsn/go-libtor/issues/12
- [ ] [File-less Persistence on Windows](https://github.com/ewhitehats/InvisiblePersistence)


## DISCLAIMER
USE FOR EDUCATIONAL PURPOSES ONLY

## Contribution
All contributions are welcome you don't need to be an expert in Go to contribute.

## Credits
- [Tor](https://www.torproject.org/)
- [Tor controller libary](https://github.com/cretz/bine)
- [Python Uacbypass and Persistence Techniques](https://github.com/rootm0s/WinPwnage)
- [Modern Cli](https://github.com/abiosoft/ishell)
- [Colored Prints](https://github.com/fatih/color)
- [Screenshot libary](https://github.com/vova616/screenshot)
