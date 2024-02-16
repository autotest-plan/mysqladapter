package mysqladapter

import (
	"context"
	"fmt"
	"net"
	"sort"

	merr "github.com/autotest-plan/error"
	"github.com/autotest-plan/log"
	pb "github.com/autotest-plan/rpcdefine/go/dbadapter"
	"google.golang.org/grpc"
)

type Server struct {
	*MySqlAdapter
	pb.DBAdapterServer
}

var dsn = "user:password@tcp(ip:port)/testdatabase?charset=utf8mb4&parseTime=True&loc=Local"

func NewServer(ctx context.Context, path []string) *Server {
	logger, err := log.NewProductLogger(path)
	if err != nil {
		fmt.Println("log模块初始化失败")
		// TODO: 实装context之后在这里优雅退出
		return nil
	}
	dba, err := NewDBAdapter(ctx, logger, map[string]string{
		"dsn": dsn,
	})
	if err != nil {
		logger.Errorf("请检查数据库连接: \n%+v\n", err)
		// TODO: 实装context之后在这里优雅退出
		return nil
	}
	return &Server{MySqlAdapter: dba}
}

// TODO: 处理Mysql查询中的范围过滤？？？
func (s *Server) Load(ctx context.Context, in *pb.Filter) (*pb.Tasks, error) {
	tasks := []*pb.Task{}
	for k, v := range in.Kv {
		tsk := s.WhereFirst(k, v)
		tasks = append(tasks, tsk)
	}
	if len(tasks) == 0 {
		return nil, merr.Error(pb.DBCode, "查询到的任务数量为0")
	}
	return &pb.Tasks{Tasks: tasks}, nil
}

func (s *Server) LoadSorted(ctx context.Context, in *pb.Filter) (*pb.Tasks, error) {
	tasks, err := s.Load(ctx, in)
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks.Tasks, func(i, j int) bool {
		iTask := tasks.Tasks[i]
		jTask := tasks.Tasks[j]
		iCorrectRate := (float64(iTask.Correct) / float64(iTask.Correct+iTask.Fault))
		jCorrectRate := (float64(jTask.Correct) / float64(jTask.Correct+jTask.Fault))
		if iTask.Parent == jTask.Parent {
			// 同级目录下，最后执行结果相同的，失败率越高优先级越高
			if iTask.Result == jTask.Result {
				return iCorrectRate < jCorrectRate
			}
			// 不同的最后执行结果，失败的优先级较高
			return !iTask.Result
		}
		// 不同目录下相同的最后结果，失败率高的优先
		if iTask.Result == jTask.Result {
			return iCorrectRate < jCorrectRate
		}
		// 不同目录下的不同最后执行结果，失败的优先
		return !iTask.Result
	})
	return tasks, nil
}

func (s *Server) Store(ctx context.Context, in *pb.Tasks) (*pb.Task, error) {
	for _, task := range in.Tasks {
		// 已存在的接口改用更新
		tsk, err := s.Load(ctx, &pb.Filter{Kv: map[string]string{
			"Name = ?": task.Name,
		}})
		if err != nil || len(tsk.Tasks) != 0 {
			s.SaveExisted(task)
			err = nil
		} else {
			_, err = s.InsertOne(task)
		}
		if err != nil {
			return &pb.Task{Result: false}, err
		}
	}
	return &pb.Task{Result: true}, nil
}

func (s *Server) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return merr.Error(pb.DBCode, "监听TCP端口失败")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDBAdapterServer(grpcServer, s)
	return grpcServer.Serve(lis)
}
