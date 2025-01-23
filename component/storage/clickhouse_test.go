package storage

import (
	"context"
	"fmt"
	"testing"
)

func TestClickhouseDatabase(t *testing.T) {
	ctx := context.Background()
	dsn := "clickhouse://localhost:9000/default"
	err := InitClickhouseDatabase(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	conn := ClickhouseConn()
	// 从 t_order_rmt 表查询id = testId 的数据
	testId := 101
	rows, err := conn.QueryContext(ctx, "SELECT id, sku_id FROM `t_order_rmt` WHERE id = ?", testId)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	fmt.Println(rows.Columns())
	for rows.Next() {
		var id int
		var skuId string
		rows.Scan(&id, &skuId)
		fmt.Println(id, skuId)
	}
}
