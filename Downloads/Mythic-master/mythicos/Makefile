.PHONY: default
default: linux ;

linux:
	cd tiger_CLI && make build_linux && mv tiger-cli ../

macos:
	cd tiger_CLI && make build_macos && mv tiger-cli ../

local:
	cd tiger_CLI && make build_local

linux_docker:
	cd tiger_CLI && make build_linux_docker && mv tiger-cli ../

macos_docker:
	cd tiger_CLI && make build_macos_docker && mv tiger-cli ../