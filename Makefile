#include .env

PROJECTNAME=$(shell basename "$(PWD)")


# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

build: go-clean go-build

go-clean:
	@echo "  >  Cleaning build cache"
	go clean

go-build:
	@echo "  >  Building binary..."
	go build -o bin/$(PROJECTNAME) -ldflags "-s -w -X 'main.Version=v1.0.0' -X 'app/build.User=$(id -u -n)' -X 'app/build.Time=$(date)'" 
   



