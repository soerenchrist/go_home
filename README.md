# Go Home

This is a pretty basic home automation server written in Go. It is currently in heavy development and taylored to my personal needs.
It is not meant to be a competitor or a replacement for awesome systems like "Home Assistant" or "OpenHAB", but a more lightweight and simple solution.


Some technologies, that are used are:
- GoLang
- Gin
- GoView/GinView
- Sqlite
- Mqtt

## Features
The feature set is currently pretty limited:
- Create devices
- Attach sensors to devices, that are either listening to external data (via http calls) or can poll for values in regular intervals
- Attach commands to devices, that can send HTTP requests to arbitrary endpoints
- Create rules to automatically invoke commands, based on sensor values
- (WIP) Listen to sensor values via MQTT

## Why?
The main aim of this project was, to get hands on experience with Go. It's only reason for existence is to be pretty lightweight and offers just the most bare features.
It can run on basically any device that runs linux

## How to run?
- Clone the repository
- Install the go programming language and cli
- Run the project using `go run .`
- Build an executable for your machine with `go build .`