package ent

//go:generate mkdir -p ../ent/proto
//go:generate protoc --go_out=../ent/proto --go_opt=paths=source_relative --go-grpc_out=../ent/proto --go-grpc_opt=paths=source_relative ./messagestore.proto

