module cron

go 1.15

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	go.etcd.io/etcd v3.3.25+incompatible
	go.mongodb.org/mongo-driver v1.5.3
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.38.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
