package client

import (
	"github.com/wpliam/common_wrap/porm/filter"
	"github.com/wpliam/common_wrap/porm/pb"
)

func NewOptions() *Options {
	return &Options{
		TimeFieldFilter: filter.NewEmptyFieldFilter(),
	}
}

// Options 选项
type Options struct {
	Fields          []string   // 自定义select,update,insert字段
	CanOpField      CanOpField // 可以操作的字段
	Table           string     // 表名
	Where           string     // where条件
	Args            []interface{}
	OrderBy         []*pb.OrderBy // 排序
	Page            *pb.Page      // 分页
	Join            string
	TimeFieldFilter filter.Filter // 时间字段过滤器
}

type CanOpField map[string]pb.CanOp

// Option 选项
type Option func(options *Options)

func WithFields(fields []string) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

func WithCanOpField(canOpField CanOpField) Option {
	return func(o *Options) {
		o.CanOpField = canOpField
	}
}

func WithTable(table string) Option {
	return func(o *Options) {
		o.Table = table
	}
}

func WithPage(page *pb.Page) Option {
	return func(o *Options) {
		o.Page = page
	}
}

func WithWhereArgs(where string, args ...interface{}) Option {
	return func(o *Options) {
		o.Where = where
		o.Args = args
	}
}

func WithJoin(join string) Option {
	return func(o *Options) {
		o.Join = join
	}
}

func WithOrderBy(orderBy ...*pb.OrderBy) Option {
	return func(o *Options) {
		o.OrderBy = orderBy
	}
}

func WithTimeField(timeField []string) Option {
	return func(os *Options) {
		os.TimeFieldFilter = filter.NewTimeFieldFilter(timeField)
	}
}

func WithTimeFieldFilter(filter filter.Filter) Option {
	return func(o *Options) {
		o.TimeFieldFilter = filter
	}
}
