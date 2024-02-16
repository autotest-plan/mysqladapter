package mysqladapter

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/autotest-plan/log"
	pb "github.com/autotest-plan/rpcdefine/go/dbadapter"
	"gorm.io/gorm"
)

var _dsn = "user:password@tcp(ip:port)/testdatabase?charset=utf8mb4&parseTime=True&loc=Local"

func get_testdatabase() *gorm.DB {
	logger, err := log.NewDevelopmentLogger([]string{"/Users/mashazheng/Projects/log/test.log"})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	db, err := NewDBAdapter(context.Background(), logger, map[string]string{"_dsn": _dsn})
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	return db.forTest()
}

func insert_data_to_table_test(db *gorm.DB, assert func(bool, string)) {
	task := pb.Task{
		Name:   "insert_data_to_table_test1",
		Result: true,
	}
	result := db.Create(&task)
	fmt.Println(task.Id)
	assert(result.Error == nil, "insert_data_to_table_test failed")
	assert(result.RowsAffected == 1, "insert_data_to_table_test failed")
}

func query_no_order_test(db *gorm.DB, assert func(bool, string)) {
	task := pb.Task{}
	db.Where("Name = ?", "db_create_test1").First(&task)
	assert(task.Id == 2, "query_no_order_test incorrect")
}

func TestMaySqlAdapter(t *testing.T) {
	assert := func(condition bool, message string) {
		if !condition {
			t.Error(message)
		}
	}
	db := get_testdatabase()
	assert(db != nil, "get_testdatabase failed")
	insert_data_to_table_test(db, assert)
	query_no_order_test(db, assert)
}
