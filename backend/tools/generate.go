package tools

//go:generate sh -c "protoc --experimental_editions --proto_path=../../proto --go_out=../gen --go_opt=paths=source_relative grpc/v1/penguin.proto"
//go:generate sh -c "protoc --experimental_editions --proto_path=../../proto --connect-go_out=../gen --connect-go_opt=paths=source_relative grpc/v1/penguin.proto"
