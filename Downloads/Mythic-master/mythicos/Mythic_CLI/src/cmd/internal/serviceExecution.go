package internal

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/tigerMeta/tiger_CLI/cmd/manager"
	"github.com/tigerMeta/tiger_CLI/cmd/utils"
	"log"
)

// ServiceStart is entrypoint from commands to start containers
func ServiceStart(containers []string) error {
	// first stop all the containers or the ones specified
	_ = manager.GetManager().StopServices(containers, config.GettigerEnv().GetBool("REBUILD_ON_START"))

	// get all the services on disk and in docker-compose currently
	diskAgents, err := manager.GetManager().GetInstalled3rdPartyServicesOnDisk()
	if err != nil {
		return err
	}
	dockerComposeContainers, err := manager.GetManager().GetAllInstalled3rdPartyServiceNames()
	if err != nil {
		return err
	}
	intendedtigerServices, err := config.GetIntendedtigerServiceNames()
	if err != nil {
		return err
	}
	currenttigerServices, err := manager.GetManager().GetCurrenttigerServiceNames()
	if err != nil {
		return err
	}
	for _, val := range currenttigerServices {
		if utils.StringInSlice(val, intendedtigerServices) {
		} else {
			_ = manager.GetManager().RemoveServices([]string{val})
		}
	}
	for _, val := range intendedtigerServices {
		AddtigerService(val, false)
	}
	// if the user didn't explicitly call out starting certain containers, then do all of them
	if len(containers) == 0 {
		containers = append(dockerComposeContainers, intendedtigerServices...)
		// make sure the ports are open that we're going to need
		TestPorts()
	}
	finalContainers := []string{}
	for _, val := range containers { // these are specified containers or all in docker compose
		if !utils.StringInSlice(val, dockerComposeContainers) && !utils.StringInSlice(val, config.tigerPossibleServices) {
			if utils.StringInSlice(val, diskAgents) {
				// the agent mentioned isn't in docker-compose, but is on disk, ask to add
				add := config.AskConfirm(fmt.Sprintf("\n%s isn't in docker-compose, but is on disk. Would you like to add it? ", val))
				if add {
					finalContainers = append(finalContainers, val)
					Add3rdPartyService(val, map[string]interface{}{}, config.GettigerEnv().GetBool("REBUILD_ON_START"))
				}
			} else {
				add := config.AskConfirm(fmt.Sprintf("\n%s isn't in docker-compose and is not on disk. Would you like to install it from https://github.com/? ", val))
				if add {
					finalContainers = append(finalContainers, val)
					installServiceByName(val)
				}
			}
		} else {
			finalContainers = append(finalContainers, val)
		}
	}
	// make sure we always update the config when starting in case .env variables changed\
	for _, service := range finalContainers {
		if utils.StringInSlice(service, config.tigerPossibleServices) {
			AddtigerService(service, config.GettigerEnv().GetBool("REBUILD_ON_START"))
		} else {
			Add3rdPartyService(service, map[string]interface{}{}, config.GettigerEnv().GetBool("REBUILD_ON_START"))
		}
	}
	manager.GetManager().TestPorts(finalContainers)
	err = manager.GetManager().StartServices(finalContainers, config.GettigerEnv().GetBool("REBUILD_ON_START"))
	err = manager.GetManager().RemoveImages()
	if err != nil {
		fmt.Printf("[-] Failed to remove images\n%v\n", err)
		return err
	}
	updateNginxBlockLists()
	generateCerts()
	TesttigerRabbitmqConnection()
	TesttigerConnection()
	Status(false)
	return nil
}
func ServiceStop(containers []string) error {
	return manager.GetManager().StopServices(containers, config.GettigerEnv().GetBool("REBUILD_ON_START"))
}
func ServiceBuild(containers []string) error {
	composeServices, err := manager.GetManager().GetAllInstalled3rdPartyServiceNames()
	if err != nil {
		log.Fatalf("[-] Failed to get installed service list: %v", err)
	}
	for _, container := range containers {
		if utils.StringInSlice(container, config.tigerPossibleServices) {
			// update the necessary docker compose entries for tiger services
			AddtigerService(container, true)
		} else if utils.StringInSlice(container, composeServices) {
			Add3rdPartyService(container, map[string]interface{}{}, true)
		}
	}
	err = manager.GetManager().BuildServices(containers)
	if err != nil {
		return err
	}
	return nil
}
func ServiceRemoveContainers(containers []string) error {
	return manager.GetManager().RemoveContainers(containers)
}

// Docker Save / Load commands

func DockerSave(containers []string) error {
	return manager.GetManager().SaveImages(containers, "saved_images")
}
func DockerLoad() error {
	return manager.GetManager().LoadImages("saved_images")
}
func DockerHealth(containers []string) {
	manager.GetManager().GetHealthCheck(containers)
}

// Build new Docker UI

func DockerBuildReactUI() error {
	if config.GettigerEnv().GetBool("tiger_REACT_DEBUG") {
		return manager.GetManager().BuildUI()
	}
	log.Fatalf("[-] Not using tiger_REACT_DEBUG to generate new UI, aborting...\n")
	return nil
}

// Docker Volume commands

func VolumesList() {
	manager.GetManager().PrintVolumeInformation()
}
func DockerRemoveVolume(volumeName string) error {
	return manager.GetManager().RemoveVolume(volumeName)
}

func DockerCopyIntoVolume(sourceFile string, destinationFileName string, destinationVolume string) {
	manager.GetManager().CopyIntoVolume(sourceFile, destinationFileName, destinationVolume)
}
func DockerCopyFromVolume(sourceVolumeName string, sourceFileName string, destinationName string) {
	manager.GetManager().CopyFromVolume(sourceVolumeName, sourceFileName, destinationName)
}
