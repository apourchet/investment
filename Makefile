
PROTOC_LOCATION=`which protoc`

default: protos

protos: 
	 $(PROTOC_LOCATION) --go_out=plugins=grpc:. protos/*.proto

clean:
	rm logs/*; \
	echo "drop measurement moment" | influx --database testdb

test:
	go test ./...


.PHONY: protos
