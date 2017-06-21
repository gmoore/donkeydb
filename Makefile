all: donkeydb donkeyclient

donkeydb: donkeydb.go
	go build donkeydb.go
	mkdir -p bin
	mv donkeydb bin

donkeyclient: donkeyclient.go
	go build donkeyclient.go
	mkdir -p bin
	mv donkeyclient bin