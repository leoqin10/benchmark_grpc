package main

import (
	"testing"

	goproto "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	model "github/leo/model"
)

func BenchmarkMarshalByMsgp(b *testing.B) {
	bb := make([]byte, 0, 1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bb, _ = group.MarshalMsg(nil)
	}
	b.ReportMetric(float64(len(bb)), "marshaledBytes")
}

func BenchmarkUnmarshalByMsgp(b *testing.B) {
	bytes, _ := group.MarshalMsg(nil)
	result := model.ColorGroup{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result.UnmarshalMsg(bytes)
	}
}

func BenchmarkMarshalByProtoBuf(b *testing.B) {
	bb := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb, _ = proto.Marshal(&protobufGroup)
	}
	b.ReportMetric(float64(len(bb)), "marshaledBytes")
}

func BenchmarkUnmarshalByProtoBuf(b *testing.B) {
	bytes, _ := proto.Marshal(&protobufGroup)
	result := model.ProtoColorGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proto.Unmarshal(bytes, &result)
	}
}

func BenchmarkMarshalByGogoProtoBuf(b *testing.B) {
	bb := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb, _ = goproto.Marshal(&gogoProtobufGroup)
	}

	b.ReportMetric(float64(len(bb)), "marshaledBytes")
}

func BenchmarkUnmarshalByGogoProtoBuf(b *testing.B) {
	bytes, _ := proto.Marshal(&gogoProtobufGroup)
	result := model.GogoProtoColorGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goproto.Unmarshal(bytes, &result)
	}
}

func BenchmarkMarshalByGencode(b *testing.B) {
	group := model.GencodeColorGroup{
		Id:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	bb := make([]byte, 0, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb, _ = group.Marshal(bb)
	}

	b.ReportMetric(float64(len(bb)), "marshaledBytes")
}

func BenchmarkUnmarshalByGencode(b *testing.B) {
	group := model.GencodeColorGroup{
		Id:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	buf, _ := group.Marshal(nil)
	var groupResult model.GencodeColorGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		groupResult.Unmarshal(buf)
	}
}
