package config

import (
	"bufio"
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/utils"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var tigerPossibleServices = []string{
	"tiger_postgres",
	"tiger_react",
	"tiger_server",
	"tiger_nginx",
	"tiger_rabbitmq",
	"tiger_graphql",
	"tiger_documentation",
	"tiger_jupyter",
	"tiger_sync",
	"tiger_grafana",
	"tiger_prometheus",
	"tiger_postgres_exporter",
}
var tigerEnv = viper.New()
var tigerEnvInfo = make(map[string]string)

// GetIntendedtigerServiceNames uses tigerEnv host values for various services to see if they should be local or remote
func GetIntendedtigerServiceNames() ([]string, error) {
	// need to see about adding services back in if they were for remote hosts before
	containerList := []string{}
	for _, service := range tigerPossibleServices {
		// service is a tiger service, but it's not in our current container list (i.e. not in docker-compose)
		switch service {
		case "tiger_react":
			if tigerEnv.GetString("tiger_REACT_HOST") == "127.0.0.1" || tigerEnv.GetString("tiger_REACT_HOST") == "tiger_react" {
				containerList = append(containerList, service)
			}
		case "tiger_nginx":
			if tigerEnv.GetString("NGINX_HOST") == "127.0.0.1" || tigerEnv.GetString("NGINX_HOST") == "tiger_nginx" {
				containerList = append(containerList, service)
			}
		case "tiger_rabbitmq":
			if tigerEnv.GetString("RABBITMQ_HOST") == "127.0.0.1" || tigerEnv.GetString("RABBITMQ_HOST") == "tiger_rabbitmq" {
				containerList = append(containerList, service)
			}
		case "tiger_server":
			if tigerEnv.GetString("tiger_SERVER_HOST") == "127.0.0.1" || tigerEnv.GetString("tiger_SERVER_HOST") == "tiger_server" {
				containerList = append(containerList, service)
			}
		case "tiger_postgres":
			if tigerEnv.GetString("POSTGRES_HOST") == "127.0.0.1" || tigerEnv.GetString("POSTGRES_HOST") == "tiger_postgres" {
				containerList = append(containerList, service)
			}
		case "tiger_graphql":
			if tigerEnv.GetString("HASURA_HOST") == "127.0.0.1" || tigerEnv.GetString("HASURA_HOST") == "tiger_graphql" {
				containerList = append(containerList, service)
			}
		case "tiger_documentation":
			if tigerEnv.GetString("DOCUMENTATION_HOST") == "127.0.0.1" || tigerEnv.GetString("DOCUMENTATION_HOST") == "tiger_documentation" {
				containerList = append(containerList, service)
			}
		case "tiger_jupyter":
			if tigerEnv.GetString("JUPYTER_HOST") == "127.0.0.1" || tigerEnv.GetString("JUPYTER_HOST") == "tiger_jupyter" {
				containerList = append(containerList, service)
			}
		case "tiger_grafana":
			if tigerEnv.GetBool("postgres_debug") {
				containerList = append(containerList, service)
			}
		case "tiger_prometheus":
			if tigerEnv.GetBool("postgres_debug") {
				containerList = append(containerList, service)
			}
		case "tiger_postgres_exporter":
			if tigerEnv.GetBool("postgres_debug") {
				containerList = append(containerList, service)
			}
			/*
				case "tiger_sync":
					if tigerSyncPath, err := filepath.Abs(filepath.Join(utils.GetCwdFromExe(), InstalledServicesFolder, "tiger_sync")); err != nil {
						fmt.Printf("[-] Failed to get the absolute path to tiger_sync: %v\n", err)
					} else if _, err = os.Stat(tigerSyncPath); !os.IsNotExist(err) {
						// this means that the tiger_sync folder _does_ exist
						containerList = append(containerList, service)
					}

			*/
		}
	}
	return containerList, nil
}
func GettigerEnv() *viper.Viper {
	return tigerEnv
}
func settigerConfigDefaultValues() {
	// global configuration ---------------------------------------------
	tigerEnv.SetDefault("debug_level", "warning")
	tigerEnvInfo["debug_level"] = `This sets the logging level for tiger_server and all installed services. Valid options are debug, info, and warning`

	tigerEnv.SetDefault("global_server_name", "tiger")
	tigerEnvInfo["global_server_name"] = `This sets the name of the tiger server that's sent down as part of webhook and logging data. This makes it easier to identify which tiger server is sending data to webhooks or logs.`

	tigerEnv.SetDefault("global_manager", "docker")
	tigerEnvInfo["global_manager"] = `This sets the management software used to control tiger. The default is "docker" which uses Docker and Docker Compose. Valid options are currently: docker. Additional PRs can be made to implement the CLIManager Interface and provide more options.`

	tigerEnv.SetDefault("global_restart_policy", "always")
	tigerEnvInfo["global_restart_policy"] = `This sets the restart policy for the containers within tiger. Valid options should only be 'always', 'unless-stopped', and 'on-failure'. The default of 'always' will ensure that tiger comes back up even when the server reboots. The 'unless-stopped' value means that tiger should come back online after reboot unless you specifically ran './tiger-cli stop' first.`

	// nginx configuration ---------------------------------------------
	tigerEnv.SetDefault("nginx_port", 7443)
	tigerEnvInfo["nginx_port"] = `This sets the port used for the Nginx reverse proxy - this port is used by the React UI and tiger's Scripting`

	tigerEnv.SetDefault("nginx_host", "tiger_nginx")
	tigerEnvInfo["nginx_host"] = `This specifies the ip/hostname for where the Nginx container executes. If this is "tiger_nginx" or "127.0.0.1", then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("nginx_bind_localhost_only", false)
	tigerEnvInfo["nginx_bind_localhost_only"] = `This specifies if the Nginx container will expose the nginx_port on 0.0.0.0 or 127.0.0.1`

	tigerEnv.SetDefault("nginx_use_ssl", true)
	tigerEnvInfo["nginx_use_ssl"] = `This specifies if the Nginx reverse proxy uses http or https`

	tigerEnv.SetDefault("nginx_use_ipv4", true)
	tigerEnvInfo["nginx_use_ipv4"] = `This specifies if the Nginx reverse proxy should bind to IPv4 or not`

	tigerEnv.SetDefault("nginx_use_ipv6", true)
	tigerEnvInfo["nginx_use_ipv6"] = `This specifies if the Nginx reverse proxy should bind to IPv6 or not`

	tigerEnv.SetDefault("nginx_use_volume", false)
	tigerEnvInfo["nginx_use_volume"] = `The Nginx container gets dynamic configuration from a variety of .env values as well as dynamically created SSL certificates. If this is True, then a docker volume is created and mounted into the container to host these pieces. If this is false, then the local filesystem is mounted inside the container instead. `

	tigerEnv.SetDefault("nginx_use_build_context", false)
	tigerEnvInfo["nginx_use_build_context"] = `The Nginx container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/nginx-docker/Dockerfile is used to generate the image used for the tiger_nginx container instead of the hosted image.`

	// tiger react UI configuration ---------------------------------------------
	tigerEnv.SetDefault("tiger_react_host", "tiger_react")
	tigerEnvInfo["tiger_react_host"] = `This specifies the ip/hostname for where the React UI container executes. If this is 'tiger_react' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("tiger_react_port", 3000)
	tigerEnvInfo["tiger_react_port"] = `This specifies the port that the React UI server listens on. This is normally accessed through the nginx reverse proxy though via /new`

	tigerEnv.SetDefault("tiger_react_bind_localhost_only", true)
	tigerEnvInfo["tiger_react_bind_localhost_only"] = `This specifies if the tiger_react container will expose the tiger_react_port on 0.0.0.0 or 127.0.0.1. Binding the localhost will still allow internal reverse proxying to work, but won't allow the service to be hit remotely. It's unlikely this will ever need to change since you should be connecting through the nginx_proxy, but would be necessary to change if the React UI were hosted on a different server.`

	tigerEnv.SetDefault("tiger_react_use_volume", false)
	tigerEnvInfo["tiger_react_use_volume"] = `This specifies if the tiger_react container will mount mount the local filesystem to serve content or use the pre-build data within the image itself. If you want to change the website that's shown, you need to mount locally and change the tiger_react_use_build_context to true'`

	tigerEnv.SetDefault("tiger_react_use_build_context", false)
	tigerEnvInfo["tiger_react_use_build_context"] = `This specifies if the tiger_react container should use the pre-built docker image hosted on GitHub's container registry (ghcr.io) or if the local tiger-react-docker/Dockerfile should be used to generate the base image for the tiger_react container`

	// documentation configuration ---------------------------------------------
	tigerEnv.SetDefault("documentation_host", "tiger_documentation")
	tigerEnvInfo["documentation_host"] = `This specifies the ip/hostname for where the documentation container executes. If this is 'documentation_host' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("documentation_port", 8090)
	tigerEnvInfo["documentation_port"] = `This specifies the port that the Documentation UI server listens on. This is normally accessed through the nginx reverse proxy though via /docs`

	tigerEnv.SetDefault("documentation_bind_localhost_only", true)
	tigerEnvInfo["documentation_bind_localhost_only"] = `This specifies if the documentation container will expose the documentation_port on 0.0.0.0 or 127.0.0.1`

	tigerEnv.SetDefault("documentation_use_volume", false)
	tigerEnvInfo["documentation_use_volume"] = `The documentation container gets dynamic from installed agents and c2 profiles. If this is True, then a docker volume is created and mounted into the container to host these pieces. If this is false, then the local filesystem is mounted inside the container instead. `

	tigerEnv.SetDefault("documentation_use_build_context", false)
	tigerEnvInfo["documentation_use_build_context"] = `The documentation container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/documentation-docker/Dockerfile is used to generate the image used for the tiger_documentation container instead of the hosted image.`

	// tiger server configuration ---------------------------------------------
	tigerEnv.SetDefault("tiger_debug_agent_message", false)
	tigerEnvInfo["tiger_debug_agent_message"] = `When this is true, tiger will send a message to the operational event log for each step of processing every agent's message. This can be a lot of messages, so do it with care, but it can be extremely valuable in figuring out issues with agent messaging. This setting can also be toggled at will in the UI on the settings page by an admin.`

	tigerEnv.SetDefault("tiger_server_port", 17443)
	tigerEnvInfo["tiger_server_port"] = `This specifies the port that the tiger_server listens on. This is normally accessed through the nginx reverse proxy though via /new. Agent and C2 Profile containers will directly access this container and port when fetching/uploading files/payloads.`

	tigerEnv.SetDefault("tiger_server_grpc_port", 17444)
	tigerEnvInfo["tiger_server_grpc_port"] = `This specifies the port that the tiger_server's gRPC functionality listens on. Translation containers will directly access this container and port when establishing gRPC functionality. C2 Profile containers will directly access this container and port when using Push Style C2 connections.`

	tigerEnv.SetDefault("tiger_server_host", "tiger_server")
	tigerEnvInfo["tiger_server_host"] = `This specifies the ip/hostname for where the tiger server container executes. If this is 'tiger_server' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("tiger_server_bind_localhost_only", true)
	tigerEnvInfo["tiger_server_bind_localhost_only"] = `This specifies if the tiger_server container will expose the tiger_server_port and tiger_server_grpc_port on 0.0.0.0 or 127.0.0.1. If you have a remote agent container connecting to tiger, you MUST set this to false so that the remote agent container can do file transfers with tiger.`

	tigerEnv.SetDefault("tiger_server_cpus", "2")
	tigerEnvInfo["tiger_server_cpus"] = `Set this to limit the maximum number of CPUs this service is able to consume`

	tigerEnv.SetDefault("tiger_server_mem_limit", "")
	tigerEnvInfo["tiger_server_mem_limit"] = `Set this to limit the maximum amount of RAM this service is able to consume`

	tigerEnv.SetDefault("tiger_server_dynamic_ports", "7000-7010")
	tigerEnvInfo["tiger_server_dynamic_ports"] = `These ports are exposed through the Docker container and provide access to SOCKS, Reverse Port Forward, and Interactive Tasking ports opened up by the tiger Server. This is a comma-separated list of ranges, so you could do 7000-7010,7012,713-720`

	tigerEnv.SetDefault("tiger_server_dynamic_ports_bind_localhost_only", false)
	tigerEnvInfo["tiger_server_dynamic_ports_bind_localhost_only"] = `This specifies if the tiger_server container will expose the dynamic_ports on 0.0.0.0 or 127.0.0.1. If you have a remote agent container connecting to tiger, you MUST set this to false so that the remote agent container can connect to gRPC.`

	tigerEnv.SetDefault("tiger_server_use_volume", false)
	tigerEnvInfo["tiger_server_use_volume"] = `The tiger_server container saves uploaded and downloaded files. If this is True, then a docker volume is created and mounted into the container to host these pieces. If this is false, then the local filesystem is mounted inside the container instead. `

	tigerEnv.SetDefault("tiger_server_use_build_context", false)
	tigerEnvInfo["tiger_server_use_build_context"] = `The tiger_server container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/tiger-docker/Dockerfile is used to generate the image used for the tiger_server container instead of the hosted image. If you want to modify the local tiger_server code then you need to set this to true and uncomment the sections of the tiger-docker/Dockerfile that copy over the existing code and build it. If you don't do this then you won't see any of your changes take effect`

	tigerEnv.SetDefault("tiger_sync_cpus", "2")
	tigerEnvInfo["tiger_sync_cpus"] = `Set this to limit the maximum number of CPUs this service is able to consume`

	tigerEnv.SetDefault("tiger_sync_mem_limit", "")
	tigerEnvInfo["tiger_sync_mem_limit"] = `Set this to limit the maximum amount of RAM this service is able to consume`

	// postgres configuration ---------------------------------------------
	tigerEnv.SetDefault("postgres_host", "tiger_postgres")
	tigerEnvInfo["postgres_host"] = `This specifies the ip/hostname for where the postgres database container executes. If this is 'tiger_postgres' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("postgres_port", 5432)
	tigerEnvInfo["postgres_port"] = `This specifies the port that the Postgres database server listens on.`

	tigerEnv.SetDefault("postgres_bind_localhost_only", true)
	tigerEnvInfo["postgres_bind_localhost_only"] = `This specifies if the tiger_postgres container will expose the postgres_port on 0.0.0.0 or 127.0.0.1`

	tigerEnv.SetDefault("postgres_db", "tiger_db")
	tigerEnvInfo["postgres_db"] = `This configures the name of the database tiger uses to store its data`

	tigerEnv.SetDefault("postgres_user", "tiger_user")
	tigerEnvInfo["postgres_user"] = `This configures the name of the database user tiger uses
`
	tigerEnv.SetDefault("postgres_password", utils.GenerateRandomPassword(30))
	tigerEnvInfo["postgres_password"] = `This is the randomly generated password that tiger_server and tiger_graphql use to connect to the tiger_postgres container`

	tigerEnv.SetDefault("postgres_cpus", "2")
	tigerEnvInfo["postgres_cpus"] = `Set this to limit the maximum number of CPUs this service is able to consume`

	tigerEnv.SetDefault("postgres_mem_limit", "")
	tigerEnvInfo["postgres_mem_limit"] = `Set this to limit the maximum amount of RAM this service is able to consume`

	tigerEnv.SetDefault("postgres_use_volume", false)
	tigerEnvInfo["postgres_use_volume"] = `The tiger_postgres container saves a database of everything that happens within tiger. If this is True, then a docker volume is created and mounted into the container to host these pieces. If this is false, then the local filesystem is mounted inside the container instead. `

	tigerEnv.SetDefault("postgres_use_build_context", false)
	tigerEnvInfo["postgres_use_build_context"] = `The tiger_postgres container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/postgres-docker/Dockerfile is used to generate the image used for the tiger_postgres container instead of the hosted image. `

	// rabbitmq configuration ---------------------------------------------
	tigerEnv.SetDefault("rabbitmq_host", "tiger_rabbitmq")
	tigerEnvInfo["rabbitmq_host"] = `This specifies the ip/hostname for where the RabbitMQ container executes. If this is 'rabbitmq_host' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("rabbitmq_port", 5672)
	tigerEnvInfo["postgres_port"] = `This specifies the port that the RabbitMQ server listens on.`

	tigerEnv.SetDefault("rabbitmq_bind_localhost_only", true)
	tigerEnvInfo["rabbitmq_bind_localhost_only"] = `This specifies if the tiger_rabbitmq container will expose the rabbitmq_port on 0.0.0.0 or 127.0.0.1. If you have a remote agent container connecting to tiger, you MUST set this to false so that the remote agent container can connect to tiger.`

	tigerEnv.SetDefault("rabbitmq_user", "tiger_user")
	tigerEnvInfo["rabbitmq_user"] = `This is the user that all containers use to connect to RabbitMQ queues`

	tigerEnv.SetDefault("rabbitmq_password", utils.GenerateRandomPassword(30))
	tigerEnvInfo["rabbitmq_password"] = `This is the randomly generated password that all containers use to connect to RabbitMQ queues`
	tigerEnv.SetDefault("rabbitmq_vhost", "tiger_vhost")

	tigerEnv.SetDefault("rabbitmq_cpus", "2")
	tigerEnvInfo["rabbitmq_cpus"] = `Set this to limit the maximum number of CPUs this service is able to consume`

	tigerEnv.SetDefault("rabbitmq_mem_limit", "")
	tigerEnvInfo["rabbitmq_mem_limit"] = `Set this to limit the maximum amount of RAM this service is able to consume`

	tigerEnv.SetDefault("rabbitmq_use_volume", false)
	tigerEnvInfo["rabbitmq_use_volume"] = `The tiger_rabbitmq container saves data about the messages queues used and their stats. If this is True, then a docker volume is created and mounted into the container to host these pieces. If this is false, then the local filesystem is mounted inside the container instead. `

	tigerEnv.SetDefault("rabbitmq_use_build_context", false)
	tigerEnvInfo["rabbitmq_use_build_context"] = `The tiger_rabbitmq container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/rabbitmq-docker/Dockerfile is used to generate the image used for the tiger_rabbitmq container instead of the hosted image. `

	// jwt configuration ---------------------------------------------
	tigerEnv.SetDefault("jwt_secret", utils.GenerateRandomPassword(30))
	tigerEnvInfo["jwt_secret"] = `This is the randomly generated password used to sign JWTs to ensure they're valid for this tiger instance`

	// hasura configuration ---------------------------------------------
	tigerEnv.SetDefault("hasura_host", "tiger_graphql")
	tigerEnvInfo["hasura_host"] = `This specifies the ip/hostname for where the Hasura GraphQL container executes. If this is 'tiger_graphql' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("hasura_port", 8080)
	tigerEnvInfo["postgres_port"] = `This specifies the port that the Hasura GraphQL server listens on. This is normally accessed through the Nginx reverse proxy though via /console`

	tigerEnv.SetDefault("hasura_bind_localhost_only", true)
	tigerEnvInfo["hasura_bind_localhost_only"] = `This specifies if the tiger_graphql container will expose the hasura_port on 0.0.0.0 or 127.0.0.1. `

	tigerEnv.SetDefault("hasura_secret", utils.GenerateRandomPassword(30))
	tigerEnvInfo["hasura_secret"] = `This is the randomly generated password you can use to connect to Hasura through the /console route through the nginx proxy`

	tigerEnv.SetDefault("hasura_cpus", "2")
	tigerEnvInfo["hasura_cpus"] = `Set this to limit the maximum number of CPUs this service is able to consume`

	tigerEnv.SetDefault("hasura_mem_limit", "2gb")
	tigerEnvInfo["hasura_mem_limit"] = `Set this to limit the maximum amount of RAM this service is able to consume`

	tigerEnv.SetDefault("hasura_use_volume", false)
	tigerEnvInfo["hasura_use_volume"] = `The tiger_graphql container has data about the roles within tiger and their permissions for various graphQL endpoints. If this is True, then the internal settings are used from the built image. If this is false, then the local filesystem is mounted inside the container instead. If you want to make any changes to the Hasura permissions, columns, or actions, then you need to make sure you first set this to false and restart tiger_graphql so that your changes are saved to disk and loaded up each time properly.`

	tigerEnv.SetDefault("hasura_use_build_context", false)
	tigerEnvInfo["hasura_use_build_context"] = `The tiger_graphql container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/hasura-docker/Dockerfile is used to generate the image used for the tiger_graphql container instead of the hosted image.`

	// docker-compose configuration ---------------------------------------------
	tigerEnv.SetDefault("COMPOSE_PROJECT_NAME", "tiger")
	tigerEnvInfo["compose_project_name"] = `This is the project name for Docker Compose - it sets the prefix of the container names and shouldn't be changed`

	tigerEnv.SetDefault("REBUILD_ON_START", false)
	tigerEnvInfo["rebuild_on_start"] = `This identifies if a container's backing image should be re-built (or re-fetched) each time you start the container. This can cause agent and c2 profile containers to have their volumes wiped on each start (and thus deleting any changes). This also drastically increases the start time for tiger overall. This should only be needed if you're doing a bunch of development on tiger itself. If you need to rebuild a specific container, you should use './tiger-cli build [container name]' instead to just rebuild that one container`

	// tiger instance configuration ---------------------------------------------
	tigerEnv.SetDefault("tiger_admin_user", "tiger_admin")
	tigerEnvInfo["tiger_admin_user"] = `This configures the name of the first user in tiger when tiger starts for the first time. After the first time tiger starts, this value is unused.`

	tigerEnv.SetDefault("tiger_admin_password", utils.GenerateRandomPassword(30))
	tigerEnvInfo["tiger_admin_password"] = `This randomly generated password is used when tiger first starts to set the password for the tiger_admin_user account. After the first time tiger starts, this value is unused`

	tigerEnv.SetDefault("default_operation_name", "Operation Chimera")
	tigerEnvInfo["default_operation_name"] = `This is used to name the initial operation created for the tiger_admin account. After the first time tiger starts, this value is unused`

	tigerEnv.SetDefault("allowed_ip_blocks", "0.0.0.0/0,::/0")
	tigerEnvInfo["allowed_ip_blocks"] = `This comma-separated set of HOST-ONLY CIDR ranges specifies where valid logins can come from. These values are used by tiger_server to block potential downloads as well as by tiger_nginx to block connections from invalid addresses as well.`

	tigerEnv.SetDefault("default_operation_webhook_url", "")
	tigerEnvInfo["default_operation_webhook_url"] = `If an operation doesn't specify their own webhook URL, then this value is used. You must instal a webhook container to have access to webhooks.`

	tigerEnv.SetDefault("default_operation_webhook_channel", "")
	tigerEnvInfo["default_operation_webhook_channel"] = `If an operation doesn't specify their own webhook channel, then this value is used. You must install a webhook container to have access to webhooks.`

	// jupyter configuration ---------------------------------------------
	tigerEnv.SetDefault("jupyter_port", 8888)
	tigerEnvInfo["jupyter_port"] = `This specifies the port for the tiger_jupyter container to expose outside of its container. This is typically accessed through the nginx proxy via /jupyter`

	tigerEnv.SetDefault("jupyter_host", "tiger_jupyter")
	tigerEnvInfo["jupyter_host"] = `This specifies the ip/hostname for where the Jupyter container executes. If this is 'jupyter_host' or '127.0.0.1', then tiger-cli assumes this container is running locally. If it's anything else, tiger-cli will not spin up this container as it assumes it lives elsewhere`

	tigerEnv.SetDefault("jupyter_token", "tiger")
	tigerEnvInfo["jupyter_token"] = `This value is used to authenticate to the Jupyter instance via the /jupyter route in the React UI`

	tigerEnv.SetDefault("jupyter_cpus", "2")
	tigerEnvInfo["jupyter_cpus"] = `Set this to limit the maximum number of CPUs this service is able to consume`

	tigerEnv.SetDefault("jupyter_mem_limit", "")
	tigerEnvInfo["jupyter_mem_limit"] = `Set this to limit the maximum amount of RAM this service is able to consume`

	tigerEnv.SetDefault("jupyter_bind_localhost_only", true)
	tigerEnvInfo["jupyter_bind_localhost_only"] = `This specifies if the tiger_jupyter container will expose the jupyter_port on 0.0.0.0 or 127.0.0.1. `

	tigerEnv.SetDefault("jupyter_use_volume", false)
	tigerEnvInfo["jupyter_use_volume"] = `The tiger_jupyter container saves data about script examples. If this is True, then a docker volume is created and mounted into the container to host these pieces. If this is false, then the local filesystem is mounted inside the container instead. `

	tigerEnv.SetDefault("jupyter_use_build_context", false)
	tigerEnvInfo["jupyter_use_build_context"] = `The tiger_jupyter container by default pulls configuration from a pre-compiled Docker image hosted on GitHub's Container Registry (ghcr.io). Setting this to "true" means that the local tiger/jupyter-docker/Dockerfile is used to generate the image used for the tiger_jupyter container instead of the hosted image.`

	// debugging help ---------------------------------------------
	tigerEnv.SetDefault("postgres_debug", false)
	tigerEnv.SetDefault("tiger_react_debug", false)
	tigerEnvInfo["tiger_react_debug"] = `Setting this to true switches the React UI from using a pre-built React UI to a live hot-reloading development server. You should only need to do this if you're planning on working on the tiger UI. Once you're doing making changes to the UI, you can run 'sudo ./tiger-cli build_ui' to compile your changes and save them to the tiger-react-docker folder. Assuming you have tiger_react_use_volume set to false, then when you disable debugging, you'll be using the newly compiled version of the UI`

	// installed service configuration ---------------------------------------------
	tigerEnv.SetDefault("installed_service_cpus", "1")
	tigerEnvInfo["installed_service_cpus"] = `Set this to limit the maximum number of CPUs that installed Agents/C2 Profile containers are allowed to consume`

	tigerEnv.SetDefault("installed_service_mem_limit", "")
	tigerEnvInfo["installed_service_mem_limit"] = `Set this to limit the maximum amount of RAM that installed Agents/C2 Profile containers are allowed to consume`

	tigerEnv.SetDefault("webhook_default_url", "")
	tigerEnvInfo["webhook_default_url"] = `This is the default webhook URL to use if one isn't configured for an operation`

	tigerEnv.SetDefault("webhook_default_callback_channel", "")
	tigerEnvInfo["webhook_default_callback_channel"] = `This is the default channel to use for new callbacks with the specified webhook url`

	tigerEnv.SetDefault("webhook_default_feedback_channel", "")
	tigerEnvInfo["webhook_default_feedback_channel"] = `This is the default channel to use for new feedback with the specified webhook url`

	tigerEnv.SetDefault("webhook_default_startup_channel", "")
	tigerEnvInfo["webhook_default_startup_channel"] = `This is the default channel to use for new startup notifications with the specified webhook url`

	tigerEnv.SetDefault("webhook_default_alert_channel", "")
	tigerEnvInfo["webhook_default_alert_channel"] = `This is the default channel to use for new alerts with the specified webhook url`

	tigerEnv.SetDefault("webhook_default_custom_channel", "")
	tigerEnvInfo["webhook_default_custom_channel"] = `This is the default channel to use for new custom messages with the specified webhook url`

}
func parsetigerEnvironmentVariables() {
	settigerConfigDefaultValues()
	tigerEnv.SetConfigName(".env")
	tigerEnv.SetConfigType("env")
	tigerEnv.AddConfigPath(utils.GetCwdFromExe())
	tigerEnv.AutomaticEnv()
	if !utils.FileExists(filepath.Join(utils.GetCwdFromExe(), ".env")) {
		_, err := os.Create(filepath.Join(utils.GetCwdFromExe(), ".env"))
		if err != nil {
			log.Fatalf("[-] .env doesn't exist and couldn't be created\n")
		}
	}
	if err := tigerEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("[-] Error while reading in .env file: %s\n", err)
		} else {
			log.Fatalf("[-]Error while parsing .env file: %s\n", err)
		}
	}
	portChecks := map[string][]string{
		"tiger_SERVER_HOST": {
			"tiger_SERVER_PORT",
			"tiger_server",
		},
		"POSTGRES_HOST": {
			"POSTGRES_PORT",
			"tiger_postgres",
		},
		"HASURA_HOST": {
			"HASURA_PORT",
			"tiger_graphql",
		},
		"RABBITMQ_HOST": {
			"RABBITMQ_PORT",
			"tiger_rabbitmq",
		},
		"DOCUMENTATION_HOST": {
			"DOCUMENTATION_PORT",
			"tiger_documentation",
		},
		"NGINX_HOST": {
			"NGINX_PORT",
			"tiger_nginx",
		},
		"tiger_REACT_HOST": {
			"tiger_REACT_PORT",
			"tiger_react",
		},
		"tiger_JUPYTER_HOST": {
			"tiger_JUPYTER_PORT",
			"tiger_jupyter",
		},
	}
	for key, val := range portChecks {
		if tigerEnv.GetString(key) == "127.0.0.1" {
			tigerEnv.Set(key, val[1])
		}
	}
	tigerEnv.Set("global_docker_latest", tigerDockerLatest)
	tigerEnvInfo["global_docker_latest"] = `This is the latest Docker Image version available for all tiger services (tiger_server, tiger_postgres, tiger-cli, etc). This is determined by the tag on the tiger branch and stamped into tiger-cli. Even if you change or remove this locally, tiger-cli will always put it back to what it was. For each of the main tiger services, if you set their *_use_build_context to false, then it's this specified Docker image version that will be fetched and used.`
	writetigerEnvironmentVariables()
}
func writetigerEnvironmentVariables() {
	c := tigerEnv.AllSettings()
	// to make it easier to read and look at, get all the keys, sort them, and display variables in order
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	f, err := os.Create(filepath.Join(utils.GetCwdFromExe(), ".env"))
	if err != nil {
		log.Fatalf("[-] Error writing out environment!\n%v", err)
	}
	defer f.Close()
	for _, key := range keys {
		if len(tigerEnv.GetString(key)) == 0 {
			_, err = f.WriteString(fmt.Sprintf("%s=\n", strings.ToUpper(key)))
		} else {
			_, err = f.WriteString(fmt.Sprintf("%s=\"%s\"\n", strings.ToUpper(key), tigerEnv.GetString(key)))
		}

		if err != nil {
			log.Fatalf("[-] Failed to write out environment!\n%v", err)
		}
	}
	return
}
func GetConfigAllStrings() map[string]string {
	c := tigerEnv.AllSettings()
	// to make it easier to read and look at, get all the keys, sort them, and display variables in order
	keys := make([]string, 0)
	for k := range c {
		keys = append(keys, k)
	}
	resultMap := make(map[string]string)
	for _, key := range keys {
		resultMap[key] = tigerEnv.GetString(key)
	}
	return resultMap
}
func GetConfigStrings(args []string) map[string]string {
	resultMap := make(map[string]string)
	allSettings := tigerEnv.AllKeys()
	for i := 0; i < len(args[0:]); i++ {
		searchRegex, err := regexp.Compile(args[i])
		if err != nil {
			log.Fatalf("[!] bad regex: %v", err)
		}
		for _, setting := range allSettings {
			if searchRegex.MatchString(strings.ToUpper(setting)) || searchRegex.MatchString(strings.ToLower(setting)) {
				resultMap[setting] = tigerEnv.GetString(setting)
			}
		}
	}
	return resultMap
}
func SetConfigStrings(key string, value string) {
	allSettings := tigerEnv.AllKeys()
	searchRegex, err := regexp.Compile(key)
	if err != nil {
		log.Fatalf("[!] bad regex: %v", err)
	}
	found := false
	for _, setting := range allSettings {
		if searchRegex.MatchString(strings.ToUpper(setting)) || searchRegex.MatchString(strings.ToLower(setting)) {
			tigerEnv.Set(setting, value)
			found = true
		}
	}
	if !found {
		log.Printf("[-] Failed to find any matching keys for %s\n", key)
		return
	}
	log.Println("[+] Configuration successfully updated. Bring containers down and up for changes to take effect.")
	writetigerEnvironmentVariables()
}
func SetNewConfigStrings(key string, value string) {
	tigerEnv.Set(key, value)
	writetigerEnvironmentVariables()
}
func GetBuildArguments() []string {
	var buildEnv = viper.New()
	buildEnv.SetConfigName("build.env")
	buildEnv.SetConfigType("env")
	buildEnv.AddConfigPath(utils.GetCwdFromExe())
	buildEnv.AutomaticEnv()
	if !utils.FileExists(filepath.Join(utils.GetCwdFromExe(), "build.env")) {
		return []string{}
	}
	if err := buildEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("[-] Error while reading in build.env file: %s\n", err)
		} else {
			log.Fatalf("[-]Error while parsing build.env file: %s\n", err)
		}
	}
	c := buildEnv.AllSettings()
	// to make it easier to read and look at, get all the keys, sort them, and display variables in order
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var args []string
	for _, key := range keys {
		args = append(args, fmt.Sprintf("%s=%s", strings.ToUpper(key), buildEnv.GetString(key)))
	}
	return args
}
func GetConfigHelp(entries []string) map[string]string {
	allSettings := tigerEnv.AllKeys()
	output := make(map[string]string)
	for i := 0; i < len(entries[0:]); i++ {
		searchRegex, err := regexp.Compile(entries[i])
		if err != nil {
			log.Fatalf("[!] bad regex: %v", err)
		}
		for _, setting := range allSettings {
			if searchRegex.MatchString(strings.ToUpper(setting)) || searchRegex.MatchString(strings.ToLower(setting)) {
				if _, ok := tigerEnvInfo[setting]; ok {
					output[setting] = tigerEnvInfo[setting]
				} else if strings.HasSuffix(setting, "use_volume") {
					output[setting] = `This creates a new volume in the format [agent]_volume that's mapped into an agent's '/tiger' directory. The first time this is created with an empty volume, the contents of your pre-built agent/c2 gets copied to the volume. If you install a new version of the agent/c2 profile or rebuild the container (either via rebuild_on_start or directly calling ./tiger-cli build [agent]) then the volume is deleted and recreated. This will DELETE any changes you made in the container. For agents this could be temporary files created as part of builds. For C2 Profiles this could be profile updates. If this is set to 'false' the no new volume is created and instead the local InstalledServices/[agent] directory is mapped into the container, preserving all changes on rebuild and restart.`
				} else if strings.HasSuffix(setting, "use_build_context") {
					output[setting] = `This setting determines if you use the pre-built image hosted on GitHub/DockerHub or if you use your local InstalledService/[agent]/Dockerfile to build a new local image for your container. If you're wanting to make changes to an agent or c2 profile (adding commands, updating code, etc), then you need to set this to 'true' and update your Dockerfile to copy in your modified code. If your container code is Golang, then you'll also need to make sure your modified Dockerfile rebuilds that Go code based on your changes (probably with a 'go build' or 'make build' command depending on the agent). If your container code is Python, then simply copying in the changes should be sufficient. Also make sure you change [agent]_use_volume to 'false' so that your changes don't get overwritten by an old volume. Alternatively, you could remove the old volume with './tiger-cli volume rm [agent]_volume' and then build your new container with './tiger-cli build [agent]'.`
				} else if strings.HasSuffix(setting, "remote_image") {
					output[setting] = `This setting configures the remote image to use if you have *_use_build_context set to false. This value should get automatically updated by the agent/c2 profile's repo as new releases are created. This value will also get updated each time you install an agent. So if you want to pull an agent's latest image, just re-install the agent (or manually update this value and restart the local container).`
				} else {

				}
			}
		}
	}
	return output
}

// https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func AskConfirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[-] Failed to read user input\n")
			return false
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		}
	}
}

// https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func AskVariable(prompt string, environmentVariable string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[-] Failed to read user input\n")
		}
		input = strings.TrimSpace(input)
		tigerEnv.Set(environmentVariable, input)
		writetigerEnvironmentVariables()
	}
}
func Initialize() {
	parsetigerEnvironmentVariables()
	writetigerEnvironmentVariables()
}
