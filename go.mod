module github.com/autotest-plan/mysqladapter

go 1.22.0

require github.com/autotest-plan/rpcdefine v0.0.0

require github.com/autotest-plan/log v0.0.0

require github.com/autotest-plan/error v0.0.0

replace github.com/autotest-plan/rpcdefine => ../rpc-define

replace github.com/autotest-plan/log => ../log

replace github.com/autotest-plan/error => ../errors

require go.uber.org/zap v1.26.0

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240205150955-31a09d347014 // indirect
	google.golang.org/grpc v1.61.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gorm.io/driver/mysql v1.5.4
)
