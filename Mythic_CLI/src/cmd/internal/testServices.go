package internal

import (
	"crypto/tls"
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/tigerMeta/tiger_CLI/cmd/manager"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func TesttigerConnection() {
	webAddress := "127.0.0.1"
	tigerEnv := config.GettigerEnv()
	if tigerEnv.GetString("NGINX_HOST") == "tiger_nginx" {
		if tigerEnv.GetBool("NGINX_USE_SSL") {
			webAddress = "https://127.0.0.1"
		} else {
			webAddress = "http://127.0.0.1"
		}
	} else {
		if tigerEnv.GetBool("NGINX_USE_SSL") {
			webAddress = "https://" + tigerEnv.GetString("NGINX_HOST")
		} else {
			webAddress = "http://" + tigerEnv.GetString("NGINX_HOST")
		}
	}
	maxCount := 10
	sleepTime := int64(10)
	count := make([]int, maxCount)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	log.Printf("[*] Waiting for tiger Server and Nginx to come online (Retry Count = %d)\n", maxCount)
	for i := range count {
		log.Printf("[*] Attempting to connect to tiger UI at %s:%d, attempt %d/%d\n", webAddress, tigerEnv.GetInt("NGINX_PORT"), i+1, maxCount)
		resp, err := http.Get(webAddress + ":" + strconv.Itoa(tigerEnv.GetInt("NGINX_PORT")))
		if err != nil {
			log.Printf("[-] Failed to make connection to host, retrying in %ds\n", sleepTime)
			log.Printf("%v\n", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == 200 || resp.StatusCode == 404 {
				log.Printf("[+] Successfully connected to tiger at " + webAddress + ":" + strconv.Itoa(tigerEnv.GetInt("NGINX_PORT")) + "\n\n")
				return
			} else if resp.StatusCode == 502 || resp.StatusCode == 504 {
				log.Printf("[-] Nginx is up, but waiting for tiger Server, retrying connection in %ds\n", sleepTime)
			} else {
				log.Printf("[-] Connection failed with HTTP Status Code %d, retrying in %ds\n", resp.StatusCode, sleepTime)
			}
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	log.Printf("[-] Failed to make connection to tiger Server\n")
	log.Printf("    This could be due to limited resources on the host (recommended at least 2CPU and 4GB RAM)\n")
	log.Printf("    If there is an issue with tiger server, use 'tiger-cli logs tiger_server' to view potential errors\n")
	Status(false)
	log.Printf("[*] Fetching logs from tiger_server now:\n")
	GetLogs("tiger_server", "500", false)
	os.Exit(1)
}
func TesttigerRabbitmqConnection() {
	rabbitmqAddress := "127.0.0.1"
	tigerEnv := config.GettigerEnv()
	rabbitmqPort := tigerEnv.GetString("RABBITMQ_PORT")
	if tigerEnv.GetString("RABBITMQ_HOST") != "tiger_rabbitmq" && tigerEnv.GetString("RABBITMQ_HOST") != "127.0.0.1" {
		rabbitmqAddress = tigerEnv.GetString("RABBITMQ_HOST")
	}
	if rabbitmqAddress == "127.0.0.1" && !manager.GetManager().IsServiceRunning("tiger_rabbitmq") {
		log.Printf("[-] Service tiger_rabbitmq should be running on the host, but isn't. Containers will be unable to connect.\nStart it by starting tiger ('sudo ./tiger-cli tiger start') or manually with 'sudo ./tiger-cli tiger start tiger_rabbitmq'\n")
		return
	}
	maxCount := 10
	var err error
	count := make([]int, maxCount)
	sleepTime := int64(10)
	log.Printf("[*] Waiting for RabbitMQ to come online (Retry Count = %d)\n", maxCount)
	for i := range count {
		log.Printf("[*] Attempting to connect to RabbitMQ at %s:%s, attempt %d/%d\n", rabbitmqAddress, rabbitmqPort, i+1, maxCount)
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/tiger_vhost", tigerEnv.GetString("RABBITMQ_USER"), tigerEnv.GetString("RABBITMQ_PASSWORD"), rabbitmqAddress, rabbitmqPort))
		if err != nil {
			log.Printf("[-] Failed to connect to RabbitMQ, retrying in %ds\n", sleepTime)
			time.Sleep(10 * time.Second)
		} else {
			conn.Close()
			log.Printf("[+] Successfully connected to RabbitMQ at amqp://%s:***@%s:%s/tiger_vhost\n\n", tigerEnv.GetString("RABBITMQ_USER"), rabbitmqAddress, rabbitmqPort)
			return
		}
	}
	log.Printf("[-] Failed to make a connection to the RabbitMQ server: %v\n", err)
	if manager.GetManager().IsServiceRunning("tiger_rabbitmq") {
		log.Printf("    The tiger_rabbitmq service is running, but tiger-cli is unable to connect\n")
	} else {
		if rabbitmqAddress == "127.0.0.1" {
			log.Printf("    The tiger_rabbitmq service isn't running, but should be running locally. Did you start it?\n")
		} else {
			log.Printf("    The tiger_rabbitmq service isn't running locally, check to make sure it's running with the proper credentials\n")
		}

	}
}
func TestPorts() error {
	intendedServices, _ := config.GetIntendedtigerServiceNames()
	manager.GetManager().TestPorts(intendedServices)
	return nil
}

func Status(verbose bool) {
	manager.GetManager().PrintConnectionInfo()
	manager.GetManager().Status(verbose)
	installedServices, err := manager.GetManager().GetInstalled3rdPartyServicesOnDisk()
	if err != nil {
		log.Fatalf("[-] failed to get installed services: %v\n", err)
	}
	tigerEnv := config.GettigerEnv()
	if len(installedServices) == 0 {
		log.Printf("[*] There are no services installed\n")
		log.Printf("    To install one, use \"sudo ./tiger-cli install github <url>\"\n")
		log.Printf("    Agents can be found at: https://github.com/tigerAgents\n")
		log.Printf("    C2 Profiles can be found at: https://github.com/tigerC2Profiles\n")
	}
	if tigerEnv.GetString("RABBITMQ_HOST") == "tiger_rabbitmq" && tigerEnv.GetBool("rabbitmq_bind_localhost_only") {
		log.Printf("\n[*] RabbitMQ is currently listening on localhost. If you have a remote Service, they will be unable to connect (i.e. one running on another server)")
		log.Printf("\n    Use 'sudo ./tiger-cli config set rabbitmq_bind_localhost_only false' and restart tiger ('sudo ./tiger-cli restart') to change this\n")
	}
	if tigerEnv.GetString("tiger_SERVER_HOST") == "tiger_server" && tigerEnv.GetBool("tiger_server_bind_localhost_only") {
		log.Printf("\n[*] tigerServer is currently listening on localhost. If you have a remote Service, they will be unable to connect (i.e. one running on another server)")
		log.Printf("\n    Use 'sudo ./tiger-cli config set tiger_server_bind_localhost_only false' and restart tiger ('sudo ./tiger-cli restart') to change this\n")
	}
	log.Printf("[*] If you are using a remote PayloadType or C2Profile, they will need certain environment variables to properly connect to tiger.\n")
	log.Printf("    Use 'sudo ./tiger-cli config service' for configs for these services.\n")
}
func GetLogs(containerName string, numLogs string, follow bool) {
	logCount, err := strconv.Atoi(numLogs)
	if err != nil {
		log.Fatalf("[-] Bad log count: %v\n", err)
	}
	manager.GetManager().GetLogs(containerName, logCount, follow)
}
func ListServices() {
	manager.GetManager().PrintAllServices()
}
