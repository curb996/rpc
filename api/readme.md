protoc --go_out=. --go-grpc_out=. api/driver.proto



protoc -I proto/ proto/kv.proto --go_out=plugins=grpc:proto/ 
