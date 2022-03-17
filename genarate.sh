protoc pb/simple.proto --go_out=plugins=grpc:.

protoc -I pb --go_out=plugins=grpc:. pb/simple.proto