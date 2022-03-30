#/bin/sh

#
go install github.com/tinylib/msgp
$GOPATH/bin/msgp -o msgp_gen.go -io=false -tests=false -file data.go

# https://github.com/protocolbuffers/protobuf/releases/tag/v3.11.1
brew install protoc
protoc --go_out=. protobuf.proto
protoc --gogofaster_out=.  -I. -I$GOPATH/src  mygogo.proto


# https://github.com/andyleap/gencode
go install github.com/andyleap/gencode
$GOPATH/bin/gencode go -schema=gencode.schema -package model
