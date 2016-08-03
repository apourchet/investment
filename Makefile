
PROTOC_LOCATION=/Users/antoine/Downloads/protoc3/bin/protoc

default: protos

protos: 
	 $(PROTOC_LOCATION) --go_out=plugins=grpc:. protos/*.proto

.PHONY: protos
