package internal

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/tigerMeta/tiger_CLI/cmd/manager"
	"github.com/tigerMeta/tiger_CLI/cmd/utils"
	"log"
	"path/filepath"
	"strings"
)

func AddtigerService(service string, removeVolume bool) {
	pStruct, err := manager.GetManager().GetServiceConfiguration(service)
	if err != nil {
		log.Fatalf("[-] Failed to get current configuration information: %v\n", err)
	}
	if _, ok := pStruct["environment"]; !ok {
		pStruct["environment"] = []interface{}{}
	}
	pStruct["labels"] = map[string]string{
		"name": service,
	}
	pStruct["hostname"] = strings.ToLower(service)
	pStruct["logging"] = map[string]interface{}{
		"driver": "json-file",
		"options": map[string]string{
			"max-file": "1",
			"max-size": "10m",
		},
	}
	pStruct["restart"] = config.GettigerEnv().GetString("global_restart_policy")
	pStruct["container_name"] = strings.ToLower(service)
	tigerEnv := config.GettigerEnv()
	volumes, _ := manager.GetManager().GetVolumes()

	switch service {
	case "tiger_postgres":
		if tigerEnv.GetBool("postgres_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./postgres-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		pStruct["cpus"] = tigerEnv.GetInt("POSTGRES_CPUS")
		if tigerEnv.GetString("postgres_mem_limit") != "" {
			pStruct["mem_limit"] = tigerEnv.GetString("postgres_mem_limit")
		}
		if tigerEnv.GetBool("postgres_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${POSTGRES_PORT}:${POSTGRES_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${POSTGRES_PORT}:${POSTGRES_PORT}",
			}
		}
		pStruct["command"] = "postgres -c \"max_connections=100\" -p ${POSTGRES_PORT} -c config_file=/etc/postgresql.conf"
		environment := []string{
			"POSTGRES_DB=${POSTGRES_DB}",
			"POSTGRES_USER=${POSTGRES_USER}",
			"POSTGRES_PASSWORD=${POSTGRES_PASSWORD}",
			"POSTGRES_PORT=${POSTGRES_PORT}",
		}
		if _, ok := pStruct["environment"]; ok {
			pStruct["environment"] = utils.UpdateEnvironmentVariables(pStruct["environment"].([]interface{}), environment)
		} else {
			pStruct["environment"] = environment
		}
		if !tigerEnv.GetBool("postgres_use_volume") {
			pStruct["volumes"] = []string{
				"./postgres-docker/database:/var/lib/postgresql/data",
				"./postgres-docker/postgres.conf:/etc/postgresql.conf",
			}
		} else {
			pStruct["volumes"] = []string{
				"tiger_postgres_volume:/var/lib/postgresql/data",
			}

		}
		if _, ok := volumes["tiger_postgres"]; !ok {
			volumes["tiger_postgres_volume"] = map[string]interface{}{
				"name": "tiger_postgres_volume",
			}
		}
	case "tiger_documentation":
		if tigerEnv.GetBool("documentation_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./documentation-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		if tigerEnv.GetBool("documentation_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${DOCUMENTATION_PORT}:${DOCUMENTATION_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${DOCUMENTATION_PORT}:${DOCUMENTATION_PORT}",
			}
		}
		pStruct["environment"] = []string{
			"DOCUMENTATION_PORT=${DOCUMENTATION_PORT}",
		}
		if !tigerEnv.GetBool("documentation_use_volume") {
			pStruct["volumes"] = []string{
				"./documentation-docker/:/src",
			}
		} else {
			pStruct["volumes"] = []string{
				"tiger_documentation_volume:/src",
			}
		}
		if _, ok := volumes["tiger_documentation"]; !ok {
			volumes["tiger_documentation_volume"] = map[string]interface{}{
				"name": "tiger_documentation_volume",
			}
		}
	case "tiger_graphql":
		if tigerEnv.GetBool("hasura_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./hasura-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		pStruct["cpus"] = tigerEnv.GetInt("HASURA_CPUS")
		if tigerEnv.GetString("hasura_mem_limit") != "" {
			pStruct["mem_limit"] = tigerEnv.GetString("hasura_mem_limit")
		}
		environment := []string{
			"HASURA_GRAPHQL_DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}",
			"HASURA_GRAPHQL_METADATA_DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}",
			"HASURA_GRAPHQL_ENABLE_CONSOLE=true",
			"HASURA_GRAPHQL_DEV_MODE=false",
			"HASURA_GRAPHQL_ADMIN_SECRET=${HASURA_SECRET}",
			"HASURA_GRAPHQL_INSECURE_SKIP_TLS_VERIFY=true",
			"HASURA_GRAPHQL_SERVER_PORT=${HASURA_PORT}",
			"HASURA_GRAPHQL_METADATA_DIR=/metadata",
			"HASURA_GRAPHQL_LIVE_QUERIES_MULTIPLEXED_REFETCH_INTERVAL=1000",
			"HASURA_GRAPHQL_AUTH_HOOK=http://${tiger_SERVER_HOST}:${tiger_SERVER_PORT}/graphql/webhook",
			"tiger_ACTIONS_URL_BASE=http://${tiger_SERVER_HOST}:${tiger_SERVER_PORT}/api/v1.4",
			"HASURA_GRAPHQL_CONSOLE_ASSETS_DIR=/srv/console-assets",
		}
		if _, ok := pStruct["environment"]; ok {
			pStruct["environment"] = utils.UpdateEnvironmentVariables(pStruct["environment"].([]interface{}), environment)
		} else {
			pStruct["environment"] = environment
		}

		if tigerEnv.GetBool("hasura_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${HASURA_PORT}:${HASURA_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${HASURA_PORT}:${HASURA_PORT}",
			}
		}
		if !tigerEnv.GetBool("hasura_use_volume") {
			pStruct["volumes"] = []string{
				"./hasura-docker/metadata:/metadata",
			}
		} else {
			delete(pStruct, "volumes")
			/*
				pStruct["volumes"] = []string{
					"tiger_graphql_volume:/metadata",
				}

			*/
		}
		/*
			if _, ok := volumes["tiger_graphql"]; !ok {
				volumes["tiger_graphql_volume"] = map[string]interface{}{
					"name": "tiger_graphql_volume",
				}
			}

		*/
	case "tiger_nginx":
		if tigerEnv.GetBool("nginx_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./nginx-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		nginxUseSSL := "ssl"
		if !tigerEnv.GetBool("NGINX_USE_SSL") {
			nginxUseSSL = ""
		}
		nginxUseIPV4 := ""
		if !tigerEnv.GetBool("NGINX_USE_IPV4") {
			nginxUseIPV4 = "#"
		}
		nginxUseIPV6 := ""
		if !tigerEnv.GetBool("NGINX_USE_IPV6") {
			nginxUseIPV6 = "#"
		}
		environment := []string{
			"DOCUMENTATION_HOST=${DOCUMENTATION_HOST}",
			"DOCUMENTATION_PORT=${DOCUMENTATION_PORT}",
			"NGINX_PORT=${NGINX_PORT}",
			"tiger_SERVER_HOST=${tiger_SERVER_HOST}",
			"tiger_SERVER_PORT=${tiger_SERVER_PORT}",
			"HASURA_HOST=${HASURA_HOST}",
			"HASURA_PORT=${HASURA_PORT}",
			"tiger_REACT_HOST=${tiger_REACT_HOST}",
			"tiger_REACT_PORT=${tiger_REACT_PORT}",
			"JUPYTER_HOST=${JUPYTER_HOST}",
			"JUPYTER_PORT=${JUPYTER_PORT}",
			fmt.Sprintf("NGINX_USE_SSL=%s", nginxUseSSL),
			fmt.Sprintf("NGINX_USE_IPV4=%s", nginxUseIPV4),
			fmt.Sprintf("NGINX_USE_IPV6=%s", nginxUseIPV6),
		}
		if _, ok := pStruct["environment"]; ok {
			environment = utils.UpdateEnvironmentVariables(pStruct["environment"].([]interface{}), environment)
		}
		var finalNginxEnv []string
		for _, val := range environment {
			if !strings.Contains(val, "NEW_UI") {
				finalNginxEnv = append(finalNginxEnv, val)
			}
		}
		pStruct["environment"] = finalNginxEnv

		if tigerEnv.GetBool("nginx_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${NGINX_PORT}:${NGINX_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${NGINX_PORT}:${NGINX_PORT}",
			}
		}
		if !tigerEnv.GetBool("nginx_use_volume") {
			pStruct["volumes"] = []string{
				"./nginx-docker/ssl:/etc/ssl/private",
				"./nginx-docker/config:/etc/nginx",
			}
		} else {
			pStruct["volumes"] = []string{
				"tiger_nginx_volume_config:/etc/nginx",
				"tiger_nginx_volume_ssl:/etc/ssl/private",
			}
		}
		if _, ok := volumes["tiger_nginx"]; !ok {
			volumes["tiger_nginx_volume_config"] = map[string]interface{}{
				"name": "tiger_nginx_volume_config",
			}
			volumes["tiger_nginx_volume_ssl"] = map[string]interface{}{
				"name": "tiger_nginx_volume_ssl",
			}
		}
	case "tiger_rabbitmq":
		if tigerEnv.GetBool("rabbitmq_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./rabbitmq-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		pStruct["cpus"] = tigerEnv.GetInt("RABBITMQ_CPUS")
		if tigerEnv.GetString("rabbitmq_mem_limit") != "" {
			pStruct["mem_limit"] = tigerEnv.GetString("rabbitmq_mem_limit")
		}
		if tigerEnv.GetBool("rabbitmq_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${RABBITMQ_PORT}:${RABBITMQ_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${RABBITMQ_PORT}:${RABBITMQ_PORT}",
			}
		}
		environment := []string{
			"RABBITMQ_USER=${RABBITMQ_USER}",
			"RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}",
			"RABBITMQ_VHOST=${RABBITMQ_VHOST}",
			"RABBITMQ_PORT=${RABBITMQ_PORT}",
		}
		if _, ok := pStruct["environment"]; ok {
			environment = utils.UpdateEnvironmentVariables(pStruct["environment"].([]interface{}), environment)
		}
		var finalRabbitEnv []string
		badRabbitMqEnvs := []string{
			"RABBITMQ_DEFAULT_USER=${RABBITMQ_USER}",
			"RABBITMQ_DEFAULT_PASS=${RABBITMQ_PASSWORD}",
			"RABBITMQ_DEFAULT_VHOST=${RABBITMQ_VHOST}",
		}
		for _, val := range environment {
			if !utils.StringInSlice(val, badRabbitMqEnvs) {
				finalRabbitEnv = append(finalRabbitEnv, val)
			}
		}
		pStruct["environment"] = finalRabbitEnv
		if !tigerEnv.GetBool("rabbitmq_use_volume") {
			pStruct["volumes"] = []string{
				"./rabbitmq-docker/storage:/var/lib/rabbitmq",
				"./rabbitmq-docker/generate_config.sh:/generate_config.sh",
				"./rabbitmq-docker/rabbitmq.conf:/tmp/base_rabbitmq.conf",
			}
		} else {
			pStruct["volumes"] = []string{
				"tiger_rabbitmq_volume:/var/lib/rabbitmq",
			}
		}
		if _, ok := volumes["tiger_rabbitmq"]; !ok {
			volumes["tiger_rabbitmq_volume"] = map[string]interface{}{
				"name": "tiger_rabbitmq_volume",
			}
		}
	case "tiger_react":
		if tigerEnv.GetBool("tiger_react_debug") {
			pStruct["build"] = map[string]interface{}{
				"context": "./tigerReactUI",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
			pStruct["volumes"] = []string{
				"./tigerReactUI/src:/app/src",
				"./tigerReactUI/public:/app/public",
				"./tigerReactUI/package.json:/app/package.json",
				"./tigerReactUI/package-lock.json:/app/package-lock.json",
				"./tiger-react-docker/tiger/public:/app/build",
			}
		} else {
			if tigerEnv.GetBool("tiger_react_use_build_context") {
				pStruct["build"] = map[string]interface{}{
					"context": "./tiger-react-docker",
					"args":    config.GetBuildArguments(),
				}
				pStruct["image"] = service
			} else {
				pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
			}

			if !tigerEnv.GetBool("tiger_react_use_volume") {
				pStruct["volumes"] = []string{
					"./tiger-react-docker/config:/etc/nginx",
					"./tiger-react-docker/tiger/public:/tiger/new",
				}
			} else {
				if removeVolume {
					log.Printf("[*] Removing old volume, %s, if it exists to make room for updated configs", "tiger_react_volume_config")
					manager.GetManager().RemoveVolume("tiger_react_volume_config")
					log.Printf("[*] Removing old volume, %s, if it exists to make room for updated UI", "tiger_react_volume_public")
					manager.GetManager().RemoveVolume("tiger_react_volume_public")
				}
				pStruct["volumes"] = []string{
					"tiger_react_volume_config:/etc/nginx",
					"tiger_react_volume_public:/tiger/new",
				}
			}
		}
		if _, ok := volumes["tiger_react"]; !ok {
			volumes["tiger_react_volume_config"] = map[string]interface{}{
				"name": "tiger_react_volume_config",
			}
			volumes["tiger_react_volume_public"] = map[string]interface{}{
				"name": "tiger_react_volume_public",
			}
		}
		if tigerEnv.GetBool("tiger_react_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${tiger_REACT_PORT}:${tiger_REACT_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${tiger_REACT_PORT}:${tiger_REACT_PORT}",
			}
		}
		pStruct["environment"] = []string{
			"tiger_REACT_PORT=${tiger_REACT_PORT}",
		}
	case "tiger_jupyter":
		if tigerEnv.GetBool("jupyter_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./jupyter-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		pStruct["cpus"] = tigerEnv.GetInt("JUPYTER_CPUS")
		if tigerEnv.GetString("jupyter_mem_limit") != "" {
			pStruct["mem_limit"] = tigerEnv.GetString("jupyter_mem_limit")
		}
		if tigerEnv.GetBool("jupyter_bind_localhost_only") {
			pStruct["ports"] = []string{
				"127.0.0.1:${JUPYTER_PORT}:${JUPYTER_PORT}",
			}
		} else {
			pStruct["ports"] = []string{
				"${JUPYTER_PORT}:${JUPYTER_PORT}",
			}
		}

		pStruct["environment"] = []string{
			"JUPYTER_TOKEN=${JUPYTER_TOKEN}",
		}
		/*
			if curConfig.InConfig("services.tiger_jupyter.deploy") {
				pStruct["deploy"] = curConfig.Get("services.tiger_jupyter.deploy")
			}

		*/
		if !tigerEnv.GetBool("jupyter_use_volume") {
			pStruct["volumes"] = []string{
				"./jupyter-docker/jupyter:/projects",
			}
		} else {
			pStruct["volumes"] = []string{
				"tiger_jupyter_volume:/projects",
			}
		}
		if _, ok := volumes["tiger_jupyter"]; !ok {
			volumes["tiger_jupyter_volume"] = map[string]interface{}{
				"name": "tiger_jupyter_volume",
			}
		}
	case "tiger_server":
		if tigerEnv.GetBool("tiger_server_use_build_context") {
			pStruct["build"] = map[string]interface{}{
				"context": "./tiger-docker",
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
		} else {
			pStruct["image"] = fmt.Sprintf("ghcr.io/its-a-feature/%s:%s", service, tigerEnv.GetString("global_docker_latest"))
		}

		pStruct["cpus"] = tigerEnv.GetInt("tiger_SERVER_CPUS")
		if tigerEnv.GetString("tiger_server_mem_limit") != "" {
			pStruct["mem_limit"] = tigerEnv.GetString("tiger_server_mem_limit")
		}
		environment := []string{
			"POSTGRES_HOST=${POSTGRES_HOST}",
			"POSTGRES_PORT=${POSTGRES_PORT}",
			"POSTGRES_PASSWORD=${POSTGRES_PASSWORD}",
			"RABBITMQ_HOST=${RABBITMQ_HOST}",
			"RABBITMQ_PORT=${RABBITMQ_PORT}",
			"RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}",
			"JWT_SECRET=${JWT_SECRET}",
			"DEBUG_LEVEL=${DEBUG_LEVEL}",
			"tiger_DEBUG_AGENT_MESSAGE=${tiger_DEBUG_AGENT_MESSAGE}",
			"tiger_ADMIN_PASSWORD=${tiger_ADMIN_PASSWORD}",
			"tiger_ADMIN_USER=${tiger_ADMIN_USER}",
			"tiger_SERVER_PORT=${tiger_SERVER_PORT}",
			"tiger_SERVER_BIND_LOCALHOST_ONLY=${tiger_SERVER_BIND_LOCALHOST_ONLY}",
			"tiger_SERVER_GRPC_PORT=${tiger_SERVER_GRPC_PORT}",
			"ALLOWED_IP_BLOCKS=${ALLOWED_IP_BLOCKS}",
			"DEFAULT_OPERATION_NAME=${DEFAULT_OPERATION_NAME}",
			"DEFAULT_OPERATION_WEBHOOK_URL=${DEFAULT_OPERATION_WEBHOOK_URL}",
			"DEFAULT_OPERATION_WEBHOOK_CHANNEL=${DEFAULT_OPERATION_WEBHOOK_CHANNEL}",
			"NGINX_PORT=${NGINX_PORT}",
			"NGINX_HOST=${NGINX_HOST}",
			"tiger_SERVER_DYNAMIC_PORTS=${tiger_SERVER_DYNAMIC_PORTS}",
			"GLOBAL_SERVER_NAME=${GLOBAL_SERVER_NAME}",
		}
		tigerServerPorts := []string{
			"${tiger_SERVER_PORT}:${tiger_SERVER_PORT}",
			"${tiger_SERVER_GRPC_PORT}:${tiger_SERVER_GRPC_PORT}",
		}
		if tigerEnv.GetBool("tiger_SERVER_BIND_LOCALHOST_ONLY") {
			tigerServerPorts = []string{
				"127.0.0.1:${tiger_SERVER_PORT}:${tiger_SERVER_PORT}",
				"127.0.0.1:${tiger_SERVER_GRPC_PORT}:${tiger_SERVER_GRPC_PORT}",
			}
		}
		dynamicPortPieces := strings.Split(tigerEnv.GetString("tiger_SERVER_DYNAMIC_PORTS"), ",")
		for _, val := range dynamicPortPieces {
			if tigerEnv.GetBool("tiger_server_dynamic_ports_bind_localhost_only") {
				tigerServerPorts = append(tigerServerPorts, fmt.Sprintf("127.0.0.1:%s:%s", val, val))
			} else {
				tigerServerPorts = append(tigerServerPorts, fmt.Sprintf("%s:%s", val, val))
			}

		}
		pStruct["ports"] = tigerServerPorts
		if _, ok := pStruct["environment"]; ok {
			pStruct["environment"] = utils.UpdateEnvironmentVariables(pStruct["environment"].([]interface{}), environment)
		} else {
			pStruct["environment"] = environment
		}
		if !tigerEnv.GetBool("tiger_server_use_volume") {
			// mount the entire directory in so that you can see changes to code too
			pStruct["volumes"] = []string{
				"./tiger-docker/src:/usr/src/app",
			}
		} else {
			// when using a volume for tiger server, just have it save off the files
			pStruct["volumes"] = []string{
				"tiger_server_volume:/usr/src/app/files",
			}
		}
		if _, ok := volumes["tiger_server"]; !ok {
			volumes["tiger_server_volume"] = map[string]interface{}{
				"name": "tiger_server_volume",
			}
		}
	case "tiger_sync":
		if absPath, err := filepath.Abs(filepath.Join(manager.GetManager().GetPathTo3rdPartyServicesOnDisk(), service)); err != nil {
			fmt.Printf("[-] Failed to get abs path for tiger_sync\n")
			return
		} else {

			pStruct["build"] = map[string]interface{}{
				"context": absPath,
				"args":    config.GetBuildArguments(),
			}
			pStruct["image"] = service
			pStruct["cpus"] = tigerEnv.GetInt("tiger_SYNC_CPUS")
			if tigerEnv.GetString("tiger_sync_mem_limit") != "" {
				pStruct["mem_limit"] = tigerEnv.GetString("tiger_sync_mem_limit")
			}
			pStruct["environment"] = []string{
				"tiger_IP=${NGINX_HOST}",
				"tiger_PORT=${NGINX_PORT}",
				"tiger_USERNAME=${tiger_ADMIN_USER}",
				"tiger_PASSWORD=${tiger_ADMIN_PASSWORD}",
				"tiger_API_KEY=${tiger_API_KEY}",
				"GHOSTWRITER_API_KEY=${GHOSTWRITER_API_KEY}",
				"GHOSTWRITER_URL=${GHOSTWRITER_URL}",
				"GHOSTWRITER_OPLOG_ID=${GHOSTWRITER_OPLOG_ID}",
				"GLOBAL_SERVER_NAME=${GLOBAL_SERVER_NAME}",
			}
			if !tigerEnv.InConfig("GHOSTWRITER_API_KEY") {
				config.AskVariable("Please enter your GhostWriter API Key", "GHOSTWRITER_API_KEY")
			}
			if !tigerEnv.InConfig("GHOSTWRITER_URL") {
				config.AskVariable("Please enter your GhostWriter URL", "GHOSTWRITER_URL")
			}
			if !tigerEnv.InConfig("GHOSTWRITER_OPLOG_ID") {
				config.AskVariable("Please enter your GhostWriter OpLog ID", "GHOSTWRITER_OPLOG_ID")
			}
			if !tigerEnv.InConfig("tiger_API_KEY") {
				config.AskVariable("Please enter your tiger API Key (optional)", "tiger_API_KEY")
			}
		}
	}
	manager.GetManager().SetVolumes(volumes)
	_ = manager.GetManager().SetServiceConfiguration(service, pStruct)
}
func Add3rdPartyService(service string, additionalConfigs map[string]interface{}, removeVolume bool) error {
	existingConfig, _ := manager.GetManager().GetServiceConfiguration(service)
	if _, ok := existingConfig["environment"]; !ok {
		existingConfig["environment"] = []interface{}{}
	}
	existingConfig["labels"] = map[string]string{
		"name": service,
	}
	existingConfig["image"] = strings.ToLower(service)
	existingConfig["hostname"] = strings.ToLower(service)
	existingConfig["logging"] = map[string]interface{}{
		"driver": "json-file",
		"options": map[string]string{
			"max-file": "1",
			"max-size": "10m",
		},
	}
	existingConfig["restart"] = config.GettigerEnv().GetString("global_restart_policy")
	existingConfig["container_name"] = strings.ToLower(service)
	existingConfig["cpus"] = config.GettigerEnv().GetInt("INSTALLED_SERVICE_CPUS")
	existingConfig["build"] = map[string]interface{}{
		"context": filepath.Join(manager.GetManager().GetPathTo3rdPartyServicesOnDisk(), service),
		"args":    config.GetBuildArguments(),
	}
	existingConfig["network_mode"] = "host"
	existingConfig["extra_hosts"] = []string{
		"tiger_server:127.0.0.1",
		"tiger_rabbitmq:127.0.0.1",
	}
	/*
		pStruct := map[string]interface{}{
			"labels": map[string]string{
				"name": service,
			},
			"image":    strings.ToLower(service),
			"hostname": service,
			"logging": map[string]interface{}{
				"driver": "json-file",
				"options": map[string]string{
					"max-file": "1",
					"max-size": "10m",
				},
			},
			"restart":        config.GettigerEnv().GetString("global_restart_policy"),
			"container_name": strings.ToLower(service),
			"cpus":           config.GettigerEnv().GetInt("INSTALLED_SERVICE_CPUS"),
		}
		pStruct["build"] = map[string]interface{}{
			"context": filepath.Join(manager.GetManager().GetPathTo3rdPartyServicesOnDisk(), service),
			"args":    config.GetBuildArguments(),
		}

	*/
	agentConfigs := config.GetConfigStrings([]string{fmt.Sprintf("%s_.*", service)})
	agentUseBuildContextKey := fmt.Sprintf("%s_use_build_context", service)
	agentRemoteImageKey := fmt.Sprintf("%s_remote_image", service)
	agentUseVolumeKey := fmt.Sprintf("%s_use_volume", service)
	if useBuildContext, ok := agentConfigs[agentUseBuildContextKey]; ok {
		if useBuildContext == "false" {
			delete(existingConfig, "build")
			existingConfig["image"] = agentConfigs[agentRemoteImageKey]
		}
	}
	if useVolume, ok := agentConfigs[agentUseVolumeKey]; ok {
		if useVolume == "true" {
			volumeName := fmt.Sprintf("%s_volume", service)
			existingConfig["volumes"] = []string{
				volumeName + ":/tiger/",
			}
			if removeVolume {
				// blow away the old volume just in case to make sure we don't carry over old data
				log.Printf("[*] Removing old volume, %s, if it exists", volumeName)
				manager.GetManager().RemoveVolume(volumeName)
			}

			// add our new volume to the list of volumes if needed
			volumes, _ := manager.GetManager().GetVolumes()
			volumes[volumeName] = map[string]string{
				"name": volumeName,
			}
			manager.GetManager().SetVolumes(volumes)
		} else {
			delete(existingConfig, "volumes")
		}
	}
	if config.GettigerEnv().GetString("installed_service_mem_limit") != "" {
		existingConfig["mem_limit"] = config.GettigerEnv().GetString("installed_service_mem_limit")
	}
	for key, element := range additionalConfigs {
		existingConfig[key] = element
	}

	environment := []string{
		"tiger_ADDRESS=http://${tiger_SERVER_HOST}:${tiger_SERVER_PORT}/agent_message",
		"tiger_WEBSOCKET=ws://${tiger_SERVER_HOST}:${tiger_SERVER_PORT}/ws/agent_message",
		"RABBITMQ_USER=${RABBITMQ_USER}",
		"RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}",
		"RABBITMQ_PORT=${RABBITMQ_PORT}",
		"RABBITMQ_HOST=${RABBITMQ_HOST}",
		"tiger_SERVER_HOST=${tiger_SERVER_HOST}",
		"tiger_SERVER_PORT=${tiger_SERVER_PORT}",
		"tiger_SERVER_GRPC_PORT=${tiger_SERVER_GRPC_PORT}",
		"WEBHOOK_DEFAULT_URL=${WEBHOOK_DEFAULT_URL}",
		"WEBHOOK_DEFAULT_CALLBACK_CHANNEL=${WEBHOOK_DEFAULT_CALLBACK_CHANNEL}",
		"WEBHOOK_DEFAULT_FEEDBACK_CHANNEL=${WEBHOOK_DEFAULT_FEEDBACK_CHANNEL}",
		"WEBHOOK_DEFAULT_STARTUP_CHANNEL=${WEBHOOK_DEFAULT_STARTUP_CHANNEL}",
		"WEBHOOK_DEFAULT_ALERT_CHANNEL=${WEBHOOK_DEFAULT_ALERT_CHANNEL}",
		"WEBHOOK_DEFAULT_CUSTOM_CHANNEL=${WEBHOOK_DEFAULT_CUSTOM_CHANNEL}",
		"DEBUG_LEVEL=${DEBUG_LEVEL}",
		"GLOBAL_SERVER_NAME=${GLOBAL_SERVER_NAME}",
	}
	existingConfig["environment"] = utils.UpdateEnvironmentVariables(existingConfig["environment"].([]interface{}), environment)
	// only add in volumes if some aren't already listed
	if _, ok := existingConfig["volumes"]; !ok {
		existingConfig["volumes"] = []string{
			filepath.Join(manager.GetManager().GetPathTo3rdPartyServicesOnDisk(), service) + ":/tiger/",
		}
	}
	return manager.GetManager().SetServiceConfiguration(service, existingConfig)
}
func RemoveService(service string) error {
	return manager.GetManager().RemoveServices([]string{service})
}

func Initialize() {
	if !manager.GetManager().CheckRequiredManagerVersion() {
		log.Fatalf("[-] Bad %s version\n", manager.GetManager().GetManagerName())
	}
	manager.GetManager().GenerateRequiredConfig()
	// based on .env, find out which tiger services are supposed to be running and add them to docker compose
	intendedtigerContainers, _ := config.GetIntendedtigerServiceNames()
	for _, container := range intendedtigerContainers {
		AddtigerService(container, false)
	}

}
