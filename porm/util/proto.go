package util

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func FieldDescMapping(message proto.Message) map[string]protoreflect.FieldDescriptor {
	fdMap := make(map[string]protoreflect.FieldDescriptor)
	messageReflect := proto.MessageReflect(message)
	descriptor := messageReflect.Descriptor()
	fields := descriptor.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		fdMap[string(fd.Name())] = fd
	}
	return fdMap
}
