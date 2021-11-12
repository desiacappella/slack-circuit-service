all: main

main: slackservice main.go
	go build

slackservice: slackservice/*.go
	cd slackservice
	go build
	cd ..
