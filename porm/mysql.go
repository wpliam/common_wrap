package porm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/wpliam/common_wrap/porm/builder"
	"github.com/wpliam/common_wrap/porm/client"
	"github.com/wpliam/common_wrap/porm/constant"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlCli struct {
	db *sql.DB
}

func NewMysqlClient(dsn string) client.Client {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("open db err" + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("ping db err" + err.Error())
	}
	return &mysqlCli{db: db}
}

func (m *mysqlCli) First(ctx context.Context, message proto.Message, opt ...client.Option) error {
	opts := client.NewOptions()
	for _, o := range opt {
		o(opts)
	}
	query, args, err := builder.Get(constant.SelectOne).Build(message, opts)
	if err != nil {
		return err
	}

	next := func(rows *sql.Rows) error {
		return ParseRowsProto(rows, message, opts.TimeFieldFilter)
	}
	return m.Query(ctx, query, next, args...)
}

func (m *mysqlCli) List(ctx context.Context, dst interface{}, opt ...client.Option) error {
	opts := client.NewOptions()
	for _, o := range opt {
		o(opts)
	}
	value := reflect.ValueOf(dst)
	// 判断是否为指针
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("dst is not pointer")
	}
	direct := reflect.Indirect(value)
	slice, err := valueType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	baseType := deref(slice.Elem())
	message, ok := reflect.New(baseType).Interface().(proto.Message)
	if !ok {
		return fmt.Errorf("struct not proto.message")
	}
	query, args, err := builder.Get(constant.SelectList).Build(message, opts)
	if err != nil {
		return err
	}
	next := func(rows *sql.Rows) error {
		val := reflect.New(baseType)
		okMessage, ook := val.Interface().(proto.Message)
		if !ook {
			return fmt.Errorf("struct not proto.message")
		}
		if err = ParseRowsProto(rows, okMessage, opts.TimeFieldFilter); err != nil {
			return err
		}
		direct.Set(reflect.Append(direct, val))
		return nil
	}
	if err = m.Query(ctx, query, next, args...); err != nil {
		return err
	}
	return m.Count(ctx, opts)
}

func (m *mysqlCli) Insert(ctx context.Context, message proto.Message, opt ...client.Option) (int64, error) {
	opts := client.NewOptions()
	for _, o := range opt {
		o(opts)
	}
	query, args, err := builder.Get(constant.Insert).Build(message, opts)
	if err != nil {
		return 0, err
	}
	result, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *mysqlCli) Update(ctx context.Context, message proto.Message, opt ...client.Option) (int64, error) {
	opts := client.NewOptions()
	for _, o := range opt {
		o(opts)
	}
	query, args, err := builder.Get(constant.Update).Build(message, opts)
	if err != nil {
		return 0, err
	}
	result, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *mysqlCli) Query(ctx context.Context, query string, next client.NextFunc, args ...any) error {
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err = next(rows); err != nil {
			return err
		}
	}
	return nil
}

func (m *mysqlCli) Count(ctx context.Context, opts *client.Options) error {
	if opts.Page != nil && opts.Page.Offset > 0 && opts.Page.Limit > 0 {
		query, args, err := builder.Get(constant.SelectCount).Build(nil, opts)
		if err != nil {
			return err
		}
		return m.db.QueryRowContext(ctx, query, args...).Scan(&opts.Page.Total)
	}
	return nil
}

func (m *mysqlCli) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}
