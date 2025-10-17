.PHONY: proto
PROTO_DIR=proto
PROTO_FILES=$(PROTO_DIR)/*.proto

proto:
	protoc --go_out=paths=source_relative:. \
	       --go-grpc_out=paths=source_relative:. \
	       $(PROTO_FILES)