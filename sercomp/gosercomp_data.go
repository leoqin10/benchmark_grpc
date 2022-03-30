package main

import (
	"github.com/golang/protobuf/proto"

	model "github/leo/model"
)

var group = model.ColorGroup{
	Id:     1,
	Name:   "Reds",
	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
}

var protobufGroup = model.ProtoColorGroup{
	Id:     proto.Int32(17),
	Name:   proto.String("Reds"),
	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
}

var gogoProtobufGroup = model.GogoProtoColorGroup{
	Id:     proto.Int32(1),
	Name:   proto.String("Reds"),
	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
}
