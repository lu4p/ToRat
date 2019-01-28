# ToRat
A Remote Administation tool written in Go using Tor as its transport mechanism
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

## DISCLAIMER
USE FOR EDUCATIONAL PURPOSES ONLY

## Credits
- Tor https://www.torproject.org/
- Tor controller libary https://github.com/cretz/bine 
- Python Uacbypass and Persistence Techniques https://github.com/rootm0s/WinPwnage 
- Modern Cli https://github.com/abiosoft/ishell 
- Colored Prints https://github.com/fatih/color 
- Screenshot libary https://github.com/vova616/screenshot
