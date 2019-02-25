# ToRat
A Remote Administration tool written in Go using Tor as its transport mechanism
currently only supporting Windows clients.

Work in Progress...

## Setup
[How to setup](https://github.com/lu4p/ToRAT/wiki/Setup)

## Repository
The important parts live in

Client https://github.com/lu4p/ToRat_client

Server https://github.com/lu4p/ToRat_server

TLS certificate generator https://github.com/lu4p/genCert

## Current Features
<<<<<<< HEAD
- the ToRAT_client communicates over TCP(with TLS) proxied through Tor with the ToRat_server (hidden service)
	- [x] anonymity of client and server
	- [x] end-to-end encryption
- Cross Platform reverse shell (Windows, Linux, Mac OS)
- Windows:
	- Multiple User Account Control Bypasses (Privilege escalation)
	- Multiple Persistence methods (User, Admin)
- optional transport without Tor
	- [x] smaller binary
	- [ ] anonymity of client and server
- embedded Tor
- Unique persistent ID for every client
	- give a client an Alias
	- all Downloads from client get saved to ./$ID/$filename

### Server Shell
- Supports multiple connections
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
**screen** | take a Screenshot of the client
**cat** | view Textfiles from the client
**alias** | give the client a custom alias
**down** | download a file from the client
**up** | upload a file to the client
**escape** | escape a command and run it in a native shell on the client
**reconnect** | tell the client to reconnect
**exit** | background current session an return to main shell
else  | the command will be executed in a native shell on the client

## Upcoming Features
- [ ] Persistence and privilege escalation for Linux and Mac OS
- [ ] Support for Android and iOS
- [ ] Cat with support for .docx .pptx .od* .pdf files
- [ ] File-less Persistence on Windows https://github.com/ewhitehats/InvisiblePersistence
- [ ] ASCII-Art Welcome Message in server shell
=======
- ToRAT communicates over reverse tcp with tls over tor with the server
- Nice Server Shell
  - Tabcomplete for commands, filenames, directories
  - Arrow Key selection of clients
  - colored
- Multiple User Account Control Bypasses (Privilege escalation)
- Multiple Persistence methods (User, Admin)
- reverse shell
- Screenshot
- Unique Hostname for every client
- give clients a custom alias
- Cat to view textfiles from client
- shred for destroying files

## Upcoming Features
- Cross Platform Client (Android, MacOs, Windows, Linux)
- Transport without Tor
- embedded Tor https://godoc.org/github.com/cretz/bine/process/embedded
- Cat with support for .docx .pptx .od* .pdf 
- Fileless Persistence https://github.com/ewhitehats/InvisiblePersistence
>>>>>>> 17b801d7c4fc5a49531b727bf60bc788064289b3

## DISCLAIMER
USE FOR EDUCATIONAL PURPOSES ONLY

## Credits
- Tor https://www.torproject.org/
- Tor controller libary https://github.com/cretz/bine
- Python Uacbypass and Persistence Techniques https://github.com/rootm0s/WinPwnage
- Modern Cli https://github.com/abiosoft/ishell
- Colored Prints https://github.com/fatih/color
- Screenshot libary https://github.com/vova616/screenshot
