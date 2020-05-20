package main

import (
	"fmt"
	"io/ioutil"
	"log"
	complex "proto/src/complexDemo"
	enums "proto/src/enumdemo"
	simplepb "proto/src/simple"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func main() {
	sn := doSimple()
	readWriteDemo(sn)
	jsonDemo(sn)
	doEnum()
	doComplex()
}

func doComplex() {

	fmt.Println()
	fmt.Println("Doing Complex message")
	cm := complex.ComplexMessage{
		OneDummy: &complex.DummyMessage{
			Id:   1,
			Name: "First Name",
		},
		MultipleDummy: []*complex.DummyMessage{
			&complex.DummyMessage{
				Id:   1,
				Name: "First Name",
			},
			&complex.DummyMessage{
				Id:   2,
				Name: "Second Name",
			},
			&complex.DummyMessage{
				Id:   3,
				Name: "Third Name",
			},
		},
	}

	fmt.Println(cm)

}

func doEnum() {

	fmt.Println()
	fmt.Println("Enum Demo")
	ep := enums.Person{
		FirstName: "VEnugopal",
		EyeColour: enums.Person_EYE_BLUE,
	}
	fmt.Println(ep)

}

func jsonDemo(pb proto.Message) {
	smAstring := toJSON(pb)
	fmt.Println(smAstring)

	sm2 := &simplepb.SimpleMessage{}
	fromJSON(smAstring, sm2)
	fmt.Println("Succfully created proto struct:", sm2)
}

func toJSON(pb proto.Message) string {
	marshler := jsonpb.Marshaler{}
	out, err := marshler.MarshalToString(pb)
	if err != nil {
		log.Fatalln("Cant convert to JSON", err)
		return " "
	}
	return out
}

func fromJSON(in string, pb proto.Message) {
	err := jsonpb.UnmarshalString(in, pb)
	if err != nil {
		log.Fatalln("COuld not unmarshal to json")
	}
}

func readWriteDemo(sn proto.Message) {
	writeToFile("simple.bin", sn)

	sn2 := &simplepb.SimpleMessage{}

	readFromFile("simple.bin", sn2)
	fmt.Println("Read the content:", sn2)
}
func doSimple() *simplepb.SimpleMessage {
	sn := simplepb.SimpleMessage{
		Id:         12345,
		IsSimple:   true,
		Name:       "My simple message",
		SampleList: []int32{1, 4, 5, 6},
	}

	fmt.Println(sn)

	sn.Name = "I renamed you"
	fmt.Println(sn)

	fmt.Println("The ID is :", sn.GetId())

	return &sn
}

func writeToFile(fname string, pb proto.Message) error {
	out, err := proto.Marshal(pb)
	if err != nil {
		log.Fatalln("cnt serialize to bytes", err)
		return err
	}
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("cnt serialize to bytes", err)
		return err
	}
	fmt.Println("Data has been writen!!!!!!!!")
	return nil
}

func readFromFile(fname string, pb proto.Message) error {
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Something went wrong while reading")
		return err
	}
	err2 := proto.Unmarshal(in, pb)
	if err2 != nil {
		log.Fatalln("Could not put protobuffer to struct")
		return err2
	}
	return nil
}
