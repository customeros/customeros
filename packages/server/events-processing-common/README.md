### Execute from folder with proto files
protoc --go_out=. --go-grpc_out=. *.proto

### Execute from root of the project
find . -type f -name "*.proto" -print0 | xargs -0 -n1 -I{} sh -c 'dir=$(dirname "{}") && protoc --go_out="$dir" "{}" --go-grpc_out="$dir" "{}"'