package porm

import (
	"context"
	"fmt"
	"github.com/wpliam/common_wrap/porm/client"
	"github.com/wpliam/common_wrap/porm/pb"
	"testing"
)

var TestOpField = map[string]pb.CanOp{
	"id":          pb.CanOp_SELECT | pb.CanOp(0),
	"name":        pb.CanOp_SELECT | pb.CanOp_INSERT | pb.CanOp(0),
	"status":      pb.CanOp_SELECT | pb.CanOp_UPDATE | pb.CanOp(0),
	"enable":      pb.CanOp_SELECT | pb.CanOp_UPDATE | pb.CanOp_INSERT | pb.CanOp(0),
	"like":        pb.CanOp_SELECT | pb.CanOp_UPDATE | pb.CanOp(0),
	"content":     pb.CanOp_SELECT | pb.CanOp_UPDATE | pb.CanOp_INSERT | pb.CanOp(0),
	"score":       pb.CanOp_SELECT | pb.CanOp_UPDATE | pb.CanOp(0),
	"create_time": pb.CanOp_SELECT | pb.CanOp(0),
	"update_time": pb.CanOp_SELECT | pb.CanOp(0),
}

var dsn = "root:12345678@tcp(127.0.0.1:3306)/orm_test?timeout=3s&readTimeout=5s&charset=utf8mb4&parseTime=true"

func Test_mysqlCli_First(t *testing.T) {
	cli := NewMysqlClient(dsn)
	ctx := context.Background()
	data := &pb.TestData{}
	if err := cli.First(ctx, data,
		client.WithTable("t_test_data"),
		client.WithTimeField([]string{"create_time", "update_time"}),
		client.WithCanOpField(TestOpField),
		client.WithWhereArgs("id = ?", 1),
	); err != nil {
		fmt.Printf("First err:%v \n", err)
		return
	}
	fmt.Printf("data:%s \n", data.String())
}

func Test_mysqlCli_List(t *testing.T) {
	var dataList []*pb.TestData
	page := &pb.Page{
		Limit:  10,
		Offset: 1,
		Total:  0,
	}
	cli := NewMysqlClient(dsn)
	ctx := context.Background()
	if err := cli.List(ctx, &dataList,
		client.WithTable("t_test_data"),
		client.WithTimeField([]string{"create_time", "update_time"}),
		client.WithPage(page),
		client.WithCanOpField(TestOpField),
		client.WithOrderBy(&pb.OrderBy{Key: "id", Desc: true}),
	); err != nil {
		fmt.Printf("List err:%v \n", err)
		return
	}
	for _, data := range dataList {
		fmt.Printf("%s \n", data.String())
	}
	fmt.Printf("page:%+v \n", page)
}

func Test_mysqlCli_Update(t *testing.T) {
	cli := NewMysqlClient(dsn)
	ctx := context.Background()
	data := &pb.TestData{
		Id:         0,
		Name:       "hahahaha",
		Status:     0,
		Enable:     false,
		Content:    []byte("hahahahha"),
		Like:       80,
		Score:      9.36,
		CreateTime: 0,
		UpdateTime: 0,
	}
	row, err := cli.Update(ctx, data,
		client.WithTable("t_test_data"),
		client.WithWhereArgs("id = ?", 1),
		client.WithCanOpField(TestOpField),
	)
	if err != nil {
		fmt.Printf("Update err:%v \n", err)
		return
	}
	fmt.Printf("row:%d \n", row)
}

func Test_mysqlCli_Insert(t *testing.T) {
	cli := NewMysqlClient(dsn)
	ctx := context.Background()
	data := &pb.TestData{
		Id:         0,
		Name:       "insert",
		Status:     0,
		Enable:     true,
		Content:    []byte("insert"),
		Like:       0,
		Score:      9.2,
		CreateTime: 0,
		UpdateTime: 0,
	}
	id, err := cli.Insert(ctx, data,
		client.WithTable("t_test_data"),
		client.WithCanOpField(TestOpField),
	)
	if err != nil {
		fmt.Printf("Insert err:%v \n", err)
		return
	}
	fmt.Printf("id:%d \n", id)
}
