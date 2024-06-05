BINARY_NAME?=tiger_server
JWT_SECRET?=
POSTGRES_HOST?="127.0.0.1"
POSTGRES_PASSWORD?="0zoHFSTbiT7AzJvnbvQ5Lt65YkIsuT"
DEBUG_LEVEL?="debug"
RABBITMQ_HOST?="127.0.0.1"
RABBITMQ_PASSWORD?="PqR9XJ957sfHqcxj6FsBMj4p"
ALLOWED_IP_BLOCKS?="0.0.0.0/0"
tiger_SERVER_GRPC_PORT?="17444"
tiger_SERVER_BIND_LOCALHOST_ONLY?="false"
tiger_SERVER_DYNAMIC_PORTS?="7000-7010"
tiger_ADMIN_USERNAME?="tiger_admin"
tiger_ADMIN_PASSWORD?="tiger_password"
tiger_DEFAULT_OPERATION?="Operation Chimera"
tiger_DEBUG_AGENT_MESSAGE?=false

export

build:
	go mod tidy
	CGO_ENABLED=0 go build -o ${BINARY_NAME} -tags="nomsgpack" .
copy_to_final:
	cp ${BINARY_NAME} /
copy_from_final:
	cp /${BINARY_NAME} .

run_custom_values:
	JWT_SECRET=${JWT_SECRET} \
POSTGRES_HOST=${POSTGRES_HOST} \
POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
DEBUG_LEVEL=${DEBUG_LEVEL} \
tiger_DEBUG_AGENT_MESSAGE=${tiger_DEBUG_AGENT_MESSAGE} \
RABBITMQ_HOST=${RABBITMQ_HOST} \
RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD} \
ALLOWED_IP_BLOCKS=${ALLOWED_IP_BLOCKS} \
tiger_ADMIN_USERNAME=${tiger_ADMIN_USERNAME} \
tiger_ADMIN_PASSWORD=${tiger_ADMIN_PASSWORD} \
tiger_DEFAULT_OPERATION=${tiger_DEFAULT_OPERATION} \
tiger_SERVER_BIND_LOCALHOST_ONLY=${tiger_SERVER_BIND_LOCALHOST_ONLY} \
tiger_SERVER_DYNAMIC_PORTS=${tiger_SERVER_DYNAMIC_PORTS} \
tiger_SERVER_GRPC_PORT=${tiger_SERVER_GRPC_PORT} \
./${BINARY_NAME}

run_local_build:
	./${BINARY_NAME}

build_protobuf_go:
	protoc --go_out=`pwd`/grpc/services --go_opt=paths=source_relative \
--go-grpc_out=`pwd`/grpc/services --go-grpc_opt=paths=source_relative \
--proto_path=`pwd`/grpc/services \
`pwd`/grpc/services/*.proto

build_protobuf_python:
	python3 -m grpc_tools.protoc --python_out=`pwd`/../../../tigerMeta/tigerContainerPyPi/tiger_container/grpc \
--pyi_out=`pwd`/../../../tigerMeta/tigerContainerPyPi/tiger_container/grpc \
--grpc_python_out=`pwd`/../../../tigerMeta/tigerContainerPyPi/tiger_container/grpc \
--proto_path=`pwd`/grpc/services \
`pwd`/grpc/services/*.proto

build_protobuf: build_protobuf_go
build_and_build_protobuf: build_protobuf build
build_final: build copy_to_final
build_and_run: build_and_build_protobuf run_local_build
local_testing: build run_custom_values

clean:
	go clean
	rm ${BINARY_NAME}
