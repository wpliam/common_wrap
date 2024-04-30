package porm

import (
	"database/sql"
	"fmt"
	"github.com/wpliam/common_wrap/porm/constant"
	"github.com/wpliam/common_wrap/porm/filter"
	"github.com/wpliam/common_wrap/porm/util"
	"reflect"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ParseRowsProto(rows *sql.Rows, message proto.Message, filter filter.Filter) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	fdMap := util.FieldDescMapping(message)
	var dest []interface{}
	for _, c := range columns {
		fd, ok := fdMap[c]
		if !ok {
			return constant.ErrFieldNotExistProtobuf
		}
		value, ok := filter.Type(c)
		if ok && value != nil {
			dest = append(dest, value)
			continue
		}
		value, err := parseProtoInterface(fd)
		if err != nil {
			return err
		}
		dest = append(dest, value)
	}
	if err = rows.Scan(dest...); err != nil {
		return err
	}
	protoReflect := proto.MessageReflect(message)
	for i, c := range columns {
		value, err := parseSQLInterface(dest[i])
		if err != nil {
			return err
		}
		fd, ok := fdMap[c]
		if !ok {
			return constant.ErrFieldNotExistProtobuf
		}
		v, ok := filter.Value(c, value)
		if ok && v != nil {
			value = v
		}
		protoReflect.Set(fd, protoreflect.ValueOf(value))
	}
	return nil
}

func parseProtoInterface(fd protoreflect.FieldDescriptor) (interface{}, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return &sql.NullBool{}, nil
	case protoreflect.StringKind:
		return &sql.NullString{}, nil
	case protoreflect.Int32Kind, protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
		return &sql.NullInt32{}, nil
	case protoreflect.Int64Kind, protoreflect.Sfixed64Kind, protoreflect.Sint64Kind:
		return &sql.NullInt64{}, nil
	case protoreflect.Uint32Kind:
		return &NullUInt32{}, nil
	case protoreflect.Uint64Kind:
		return &NullUInt64{}, nil
	case protoreflect.FloatKind:
		return &NullFloat32{}, nil
	case protoreflect.DoubleKind:
		return &sql.NullFloat64{}, nil
	case protoreflect.BytesKind:
		return &NullBytes{}, nil
	case protoreflect.GroupKind, protoreflect.MessageKind:
		return nil, fmt.Errorf("not support groupKind or messageKind")
	default:
		return nil, fmt.Errorf("unknown kind[%T]", fd.Kind())
	}
}

func parseSQLInterface(val any) (any, error) {
	switch v := val.(type) {
	case *sql.NullBool:
		return v.Bool, nil
	case *sql.NullInt32:
		return v.Int32, nil
	case *NullUInt32:
		return v.UInt32, nil
	case *sql.NullInt64:
		return v.Int64, nil
	case *NullUInt64:
		return v.UInt64, nil
	case *NullFloat32:
		return v.Float32, nil
	case *sql.NullFloat64:
		return v.Float64, nil
	case *NullBytes:
		return v.Bytes, nil
	case *sql.NullString:
		return v.String, nil
	case *sql.NullTime:
		return v.Time, nil
	default:
		return val, nil
	}
}

// deref is Indirect for reflect.Types
func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func valueType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}
