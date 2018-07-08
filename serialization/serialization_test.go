package main

import (
	"encoding/json"
	"testing"

	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"

	gogopb "github.com/tddhit/alexander/serialization/gogopb"
	pb "github.com/tddhit/alexander/serialization/pb"
)

const (
	Male = iota
	Female
)

var u = &pb.User{
	Name:    "小明",
	Age:     10,
	Sex:     Male,
	Career:  "学生",
	Profile: []byte("小明是人物的代称，因这个名字天生有冷笑话的气质，又在学校时以说冷笑话出名，因此成为很多笑话中的主角名字，常出现于小学初中作文、英语、数学题、物理题、化学题和笑话中。此名拥有极为悠久的历史内涵，在中国的各种“名著”上皆有出现，体现了中国丰富的华夏文明，甚至一度代表中国冲出亚洲，成为了”最有影响力的名字“之一。"),
}
var u2 = &gogopb.User{
	Name:    u.Name,
	Age:     u.Age,
	Sex:     u.Sex,
	Career:  u.Career,
	Profile: u.Profile,
}

func BenchmarkJsonMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(u)
	}
}

func BenchmarkJsoniterMarshal(b *testing.B) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		json.Marshal(u)
	}
}

func BenchmarkPBMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		proto.Marshal(u)
	}
}

func BenchmarkGOGOPBMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gogoproto.Marshal(u2)
	}
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	out, _ := json.Marshal(u)
	user := &pb.User{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(out, user)
	}
}

func BenchmarkJsoniterUnmarshal(b *testing.B) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	out, _ := json.Marshal(u)
	user := &pb.User{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(out, user)
	}
}

func BenchmarkPBUnmarshal(b *testing.B) {
	out, _ := proto.Marshal(u)
	user := &pb.User{}
	for i := 0; i < b.N; i++ {
		proto.Unmarshal(out, user)
	}
}

func BenchmarkGOGOPBUnmarshal(b *testing.B) {
	out, _ := gogoproto.Marshal(u2)
	user := &gogopb.User{}
	for i := 0; i < b.N; i++ {
		gogoproto.Unmarshal(out, user)
	}
}
