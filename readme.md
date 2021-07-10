# get-ca

A Go program to get the CA certificate of a website.

## Background
Why make this program at all, when most high-level HTTP clients are already able to automatically fetch CA certificates for HTTPS exchange?  
  
Precisely because **high-level** HTTP clients are able to. But what about small micro-controllers like ESP8266/ESP32 with limited processing power? They cannot get a CA certificate on their own because they do not have the capability of high-level CA cert fetching. This app tries to solve that issue.

## How it works
The general idea is to make the micro-controller (which do not have the capability of fetching a web site's CA cert on its own) to access this app via an insecure HTTP port, fetch the CA cert, and then use the CA cert to do secure HTTPS operation.  

The main goal is to ensure that CA certs are not hard-coded into the firmware/source code to avoid updating the firmware every few year(s), but can be fetched dynamically through an external source.  
  
To invoke a CA cert fetching, just run this app. But first, create an `.env` file which contains:

```sh
SERVER_PORT=8079 # or whatever server port you'd like to use
```

build the app and then run it in a linux server which has the `openssl` command:
```sh
go build && ./get-ca
```

then invoke the HTTP endpoint like this:
```
http://localhost:8079?url=platform.antares.id&port=8443
```

and voila, there's your CA cert.

This app runs `openssl` in the background each time the CA fetch is invoked.  
  
To run this app as a daemon, you might want to make a systemd service like this:
```sh
# /etc/systemd/system/get-ca.service
[Unit]
Description=Get CA Service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=/root/get-ca/get-ca
WorkingDirectory=/root/get-ca

[Install]
WantedBy=multi-user.target
```
then run:

```
sudo systemctl start get-ca &&\
sudo systemctl enable get-ca
```

## Cross-compiling
You can cross-compile this app to run in any linux distribution by using [xgo](https://github.com/karalabe/xgo). 

Be sure to set [GOPATH](https://golang.org/doc/gopath_code) and install [docker](https://www.docker.com/) first in your development environment. To compile for `amd64` linux, just run:

```sh
xgo --targets=linux/amd64 .
# or --targets=linux/384 for 32-bit systems. But who uses 32-bit systems nowadays anyway?
```

## Disclaimer
I do not know if this is the best practice to achieve non hard-coded CA cert for microcontroller HTTPS operations, but currently this is the best option I can think of. If you have any inputs on security vulnerabilities please contact my email at [valianmasdani@gmail.com](mailto:valianmasdani@gmail.com) or open a new issue in this repository, to ensure that I get notified.  
  
This software is MIT Licensed.