PROTO_DIR = proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT = .


.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_FILES)

.PHONY: build
build:
	docker build -t pchhalotre/relief-ops.${SERVICE_NAME}:latest -f ${DOCKERFILE} . && \
	docker image prune -f && \
	docker push pchhalotre/relief-ops.${SERVICE_NAME}:latest

.PHONY: build-all
build-all:
	make build SERVICE_NAME=api-gateway DOCKERFILE=services/api-gateway/Dockerfile
	make build SERVICE_NAME=user-service DOCKERFILE=services/user-service/Dockerfile
	make build SERVICE_NAME=disaster-service DOCKERFILE=services/disaster-service/Dockerfile
	make build SERVICE_NAME=resource-service DOCKERFILE=services/resource-service/Dockerfile