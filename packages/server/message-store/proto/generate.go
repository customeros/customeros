package gen

//go:generate mkdir -p ../gen/proto
//go:generate protoc --go_out=../gen/proto --go_opt=paths=source_relative --go-grpc_out=../gen/proto --go-grpc_opt=paths=source_relative ./messagestore.proto
