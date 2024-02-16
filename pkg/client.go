package mysqladapter

import (
	"context"
	"fmt"
	"strings"

	merr "github.com/autotest-plan/error"
	mlog "github.com/autotest-plan/log"
	pb "github.com/autotest-plan/rpcdefine/go/dbadapter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySqlAdapter struct {
	db *gorm.DB
}

// 查询一条符合条件的表记录
func (dba *MySqlAdapter) WhereFirst(query interface{}, args ...interface{}) *pb.Task {
	task := pb.Task{}
	dba.db.Where(query, args).First(&task)
	return &task
}

// 往表中插入一条数据
func (dba *MySqlAdapter) InsertOne(task *pb.Task) (int64, error) {
	result := dba.db.Create(&task)
	return int64(task.Id), result.Error
}

// 更新一条已存在的数据
func (dba *MySqlAdapter) SaveExisted(task *pb.Task) {
	dba.db.Save(task)
}

// 测试用
func (ma *MySqlAdapter) forTest() *gorm.DB {
	return ma.db
}

// "dsn": "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func NewDBAdapter(ctx context.Context, log *mlog.Logger, config interface{}) (*MySqlAdapter, error) {
	m, ok := config.(map[string]string)
	if !ok {
		log.Error("配置参数要求map类型")
		return nil, merr.Error(pb.DBCode, "缺少gorm链接数据库的dsn")
	}
	dsn, ok := m[strings.ToLower("dsn")]
	if !ok {
		log.Error("缺少gorm链接数据库的dsn")
		return nil, merr.Error(pb.DBCode, "缺少gorm链接数据库的dsn")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("连接数据库失败, %s", err)
		return nil, merr.Error(pb.DBCode, fmt.Sprintf("连接数据库失败: %s", err))
	}
	return &MySqlAdapter{
		db: db,
	}, nil
}
