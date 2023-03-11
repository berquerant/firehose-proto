.PHONY: dev
dev: lint test

#
# code generation
#

.PHONY: generate
generate: gen-proto re-go-generate

#
# docker images
#

.PHONY: docker
docker: docker-format firehose-docker-protobuf

.PHONY: firehose-docker-protobuf
firehose-docker-protobuf:
	@(cd docker/firehose-docker-protobuf && ./bin/build.sh)

.PHONY: docker-format
docker-format:
	@docker build docker/format --tag firehose-proto-format

#
# lint, format
#

.PHONY: lint
lint:
	@bin/format.sh dry
	@bin/lint.sh

.PHONY: format
format:
	@bin/format.sh

#
# protobuf
#

%.pb.go: %.proto
	@bin/protoc.sh $<

.PHONY: gen-proto
gen-proto:
	@bin/gen-proto.sh

#
# test
#

.PHONY: test
test:
	go test ./...

#
# go generate
#

.PHONY: re-go-generate
re-go-generate: clean-go-generated go-generate

.PHONY: clean-go-generated
clean-go-generated:
	find . -name "*_generated.go" -type f -delete

.PHONY: go-generate
go-generate:
	go generate ./...
