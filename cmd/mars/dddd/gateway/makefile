.PHONY: lint
lint: $(LINTER)
	@${TOOLS_SHELL} lint
	@echo "lint check finished"

#cmd/mars/dddd/gateway/proto/protocol/third_party --proto_path=
.PHONY: proto
proto:
	protoc   --proto_path=./proto/protocol --proto_path=.  --proto_path=./proto/protocol/third_party --go_out=paths=source_relative:./proto/client/ ./proto/protocol/*.proto
