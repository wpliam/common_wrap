package builder

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/wpliam/common_wrap/porm/client"
	"github.com/wpliam/common_wrap/porm/constant"
	"github.com/wpliam/common_wrap/porm/pb"
	"github.com/wpliam/common_wrap/porm/util"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrNoAssignWhereCondition = fmt.Errorf("no assign where condition")
	ErrNoAssignColumn         = fmt.Errorf("no assign column")
)

// Builder 建造sql接口
type Builder interface {
	Build(message proto.Message, opts *client.Options) (string, []interface{}, error)
}

func Get(buildType string) Builder {
	switch buildType {
	case constant.SelectOne, constant.SelectList, constant.SelectCount:
		return NewSelectBuilder(buildType)
	case constant.Insert:
		return NewInsertBuilder()
	case constant.Update:
		return NewUpdateBuilder()
	default:
		return nil
	}
}

func buildColumns(fdMap map[string]protoreflect.FieldDescriptor, op pb.CanOp, opts *client.Options) ([]string, error) {
	var columns []string
	if len(opts.Fields) == 0 {
		for key := range fdMap {
			columns = append(columns, key)
		}
	} else {
		for _, field := range opts.Fields {
			if _, ok := fdMap[field]; !ok {
				return nil, constant.ErrFieldNotExistProtobuf
			}
			columns = append(columns, field)
		}
	}
	if len(opts.CanOpField) == 0 {
		return columns, nil
	}
	var canOpColumn []string
	for _, c := range columns {
		admitted, ok := opts.CanOpField[c]
		if !ok || admitted&op != op {
			continue
		}
		canOpColumn = append(canOpColumn, c)
	}
	return canOpColumn, nil
}

func buildColumnAndValue(message proto.Message, op pb.CanOp, opts *client.Options) ([]string, []interface{}, error) {
	protoReflect := proto.MessageReflect(message)
	fdMap := util.FieldDescMapping(message)
	columns, err := buildColumns(fdMap, op, opts)
	if err != nil {
		return nil, nil, err
	}
	customField := false
	if len(opts.Fields) > 0 {
		customField = true
	}
	var fields []string
	var args []interface{}
	for _, c := range columns {
		fd, ok := fdMap[c]
		if !ok {
			return nil, nil, constant.ErrFieldNotExistProtobuf
		}
		value := protoReflect.Get(fd)
		if !value.IsValid() {
			continue
		}
		v := value.Interface()
		// 非自定义字段默认值不进行写入
		if !customField && value.Equal(fd.Default()) {
			continue
		}
		if val, ook := opts.TimeFieldFilter.Value(string(fd.Name()), v); ook {
			v = val
		}
		fields = append(fields, c)
		args = append(args, v)
	}
	return fields, args, nil
}
