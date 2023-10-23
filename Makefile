.PHONY: all
export GO111MODULE=on

APP=diceDB

help:
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

run: # run's your go code from source
run:
	go run cmd/main.go -host=127.0.0.1 -port=7379