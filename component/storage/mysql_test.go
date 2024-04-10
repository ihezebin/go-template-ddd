package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	pwd, _ := os.Getwd()
	config.GetConfig().Pwd = filepath.Join(pwd, "../../")
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	err := InitMySQLStorageClient(ctx, "mysql://root:root@127.0.0.1:3306/go-template-ddd?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}

	db := MySQLStorageDatabase()
	tx := db.Begin()
	// 执行一系列数据库操作
	// 如果出现错误，可以回滚事务
	if err = tx.Create(&entity.Example{
		Id:       primitive.NewObjectID().Hex(),
		Username: "admin",
		Password: "123456",
		Email:    "6wqz8@example.com",
		Salt:     "123456",
	}).Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// 提交事务
	tx.Commit()
}
