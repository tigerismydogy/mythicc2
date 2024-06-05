package utils

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// server configuration
	AdminUser               string
	AdminPassword           string
	DefaultOperationName    string
	DefaultOperationWebhook string
	DefaultOperationChannel string
	AllowedIPBlocks         []*net.IPNet
	DebugAgentMessage       bool
	DebugLevel              string
	ServerPort              uint
	ServerBindLocalhostOnly bool
	ServerDynamicPorts      []uint32
	ServerGRPCPort          uint
	GlobalServerName        string

	// rabbitmq configuration
	RabbitmqHost     string
	RabbitmqPort     uint
	RabbitmqUser     string
	RabbitmqPassword string
	RabbitmqVHost    string

	// postgres configuration
	PostgresHost     string
	PostgresPort     uint
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string

	// jwt configuration
	JWTSecret []byte
}

var (
	tigerConfig = Config{}
)

func Initialize() {
	tigerEnv := viper.New()
	// tiger config
	tigerEnv.SetDefault("tiger_debug_agent_message", false)
	tigerEnv.SetDefault("tiger_server_port", 17443)
	tigerEnv.SetDefault("tiger_server_bind_localhost_only", true)
	tigerEnv.SetDefault("tiger_server_dynamic_ports", "7000-7010")
	tigerEnv.SetDefault("tiger_server_grpc_port", 17444)
	tigerEnv.SetDefault("tiger_admin_user", "tiger_admin")
	tigerEnv.SetDefault("tiger_admin_password", "tiger_password")
	tigerEnv.SetDefault("allowed_ip_blocks", "0.0.0.0/0")
	tigerEnv.SetDefault("debug_level", "warning")
	// postgres configuration
	tigerEnv.SetDefault("postgres_host", "tiger_postgres")
	tigerEnv.SetDefault("postgres_port", 5432)
	tigerEnv.SetDefault("postgres_db", "tiger_db")
	tigerEnv.SetDefault("postgres_user", "tiger_user")
	tigerEnv.SetDefault("postgres_password", "")
	// rabbitmq configuration
	tigerEnv.SetDefault("rabbitmq_host", "tiger_rabbitmq")
	tigerEnv.SetDefault("rabbitmq_port", 5672)
	tigerEnv.SetDefault("rabbitmq_user", "tiger_user")
	tigerEnv.SetDefault("rabbitmq_password", "")
	tigerEnv.SetDefault("rabbitmq_vhost", "tiger_vhost")
	// jwt configuration
	tigerEnv.SetDefault("jwt_secret", "")
	// default operation configuration
	tigerEnv.SetDefault("default_operation_name", "Operation Chimera")
	tigerEnv.SetDefault("default_operation_webhook_url", "")
	tigerEnv.SetDefault("default_operation_webhook_channel", "")
	tigerEnv.SetDefault("global_server_name", "tiger")
	// pull in environment variables and configuration from .env if needed
	tigerEnv.SetConfigName(".env")
	tigerEnv.SetConfigType("env")
	tigerEnv.AddConfigPath(getCwdFromExe())
	tigerEnv.AutomaticEnv()
	if !fileExists(filepath.Join(getCwdFromExe(), ".env")) {
		_, err := os.Create(filepath.Join(getCwdFromExe(), ".env"))
		if err != nil {
			log.Fatalf("[-] .env doesn't exist and couldn't be created: %v", err)
		}
	}
	if err := tigerEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("[-] Error while reading in .env file: %v", err)
		} else {
			log.Fatalf("[-]Error while parsing .env file: %v", err)
		}
	}
	setConfigFromEnv(tigerEnv)
}

func setConfigFromEnv(tigerEnv *viper.Viper) {
	// tiger server configuration
	tigerConfig.DebugAgentMessage = tigerEnv.GetBool("tiger_debug_agent_message")
	tigerConfig.ServerPort = tigerEnv.GetUint("tiger_server_port")
	tigerConfig.ServerBindLocalhostOnly = tigerEnv.GetBool("tiger_server_bind_localhost_only")
	tigerConfig.ServerGRPCPort = tigerEnv.GetUint("tiger_server_grpc_port")
	dynamicPorts := tigerEnv.GetString("tiger_server_dynamic_ports")
	for _, port := range strings.Split(dynamicPorts, ",") {
		if strings.Contains(port, "-") {
			rangePieces := strings.Split(port, "-")
			if len(rangePieces) != 2 {
				log.Printf("[-] tiger_server_dynamic_ports value range was malformed: %s:%s\n", "port", port)
			} else {
				lowerRange, err := strconv.Atoi(rangePieces[0])
				if err != nil {
					log.Printf("[-] Failed to parse port for tiger_server_dynamic_ports: %v\n", err)
					continue
				}
				upperRange, err := strconv.Atoi(rangePieces[1])
				if err != nil {
					log.Printf("[-] Failed to parse port for tiger_server_dynamic_ports: %v\n", err)
					continue
				}
				if lowerRange > upperRange {
					log.Printf("[-] lower range port for tiger_server_dynamic_ports is higher than upper range: %s\n", port)
					continue
				}
				for i := lowerRange; i <= upperRange; i++ {
					tigerConfig.ServerDynamicPorts = append(tigerConfig.ServerDynamicPorts, uint32(i))
				}
			}
		} else {
			intPort, err := strconv.Atoi(port)
			if err == nil {
				tigerConfig.ServerDynamicPorts = append(tigerConfig.ServerDynamicPorts, uint32(intPort))
			} else {
				log.Printf("[-] Failed to parse port for tiger_server_dynamic_ports: %v - %v\n", port, err)
			}

		}
	}
	tigerConfig.AdminUser = tigerEnv.GetString("tiger_admin_user")
	tigerConfig.AdminPassword = tigerEnv.GetString("tiger_admin_password")
	tigerConfig.DefaultOperationName = tigerEnv.GetString("default_operation_name")
	tigerConfig.DefaultOperationWebhook = tigerEnv.GetString("default_operation_webhook_url")
	tigerConfig.DefaultOperationChannel = tigerEnv.GetString("default_operation_webhook_channel")
	tigerConfig.GlobalServerName = tigerEnv.GetString("global_server_name")
	allowedIPBlocks := []*net.IPNet{}
	for _, ipBlock := range strings.Split(tigerEnv.GetString("allowed_ip_blocks"), ",") {
		if _, subnet, err := net.ParseCIDR(ipBlock); err != nil {
			log.Printf("[-] Failed to parse CIDR block: %s\n", ipBlock)
		} else {
			allowedIPBlocks = append(allowedIPBlocks, subnet)
		}
	}
	tigerConfig.AllowedIPBlocks = allowedIPBlocks
	tigerConfig.DebugLevel = tigerEnv.GetString("debug_level")
	// postgres configuration
	tigerConfig.PostgresHost = tigerEnv.GetString("postgres_host")
	tigerConfig.PostgresPort = tigerEnv.GetUint("postgres_port")
	tigerConfig.PostgresDB = tigerEnv.GetString("postgres_db")
	tigerConfig.PostgresUser = tigerEnv.GetString("postgres_user")
	tigerConfig.PostgresPassword = tigerEnv.GetString("postgres_password")
	// rabbitmq configuration
	tigerConfig.RabbitmqHost = tigerEnv.GetString("rabbitmq_host")
	tigerConfig.RabbitmqPort = tigerEnv.GetUint("rabbitmq_port")
	tigerConfig.RabbitmqUser = tigerEnv.GetString("rabbitmq_user")
	tigerConfig.RabbitmqPassword = tigerEnv.GetString("rabbitmq_password")
	tigerConfig.RabbitmqVHost = tigerEnv.GetString("rabbitmq_vhost")
	// jwt configuration
	tigerConfig.JWTSecret = []byte(tigerEnv.GetString("jwt_secret"))
}

func getCwdFromExe() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("[-] Failed to get path to current executable: %v", err)
	}
	return filepath.Dir(exe)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}

func SetConfigValue(configKey string, configValue interface{}) error {
	switch configKey {
	case "tiger_DEBUG_AGENT_MESSAGE":
		tigerConfig.DebugAgentMessage = configValue.(bool)
	default:
		return errors.New("unknown configKey to update")
	}
	return nil
}

func GetGlobalConfig() map[string]interface{} {
	return map[string]interface{}{
		"tiger_DEBUG_AGENT_MESSAGE": tigerConfig.DebugAgentMessage,
	}
}
