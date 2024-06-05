package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// configCmd represents the config command
var updateCmd = &cobra.Command{
	Use:   "update [branch name]",
	Short: "Check for tiger updates",
	Long:  `Check for a tiger update against a specific branch or against HEAD by default.`,
	Run:   updateCheck,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(
		&branch,
		"branch",
		"b",
		"",
		`Check update status from a specific branch instead of HEAD`,
	)
}

func updateCheck(cmd *cobra.Command, args []string) {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var client = &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	urlBase := "https://raw.githubusercontent.com/its-a-feature/tiger/"

	targetURL := urlBase + "master"
	if len(args) == 1 {
		branch = args[0]
	}
	if len(branch) > 0 {
		targetURL = urlBase + branch
	}

	if tigerNeedsUpdate, err := checktigerVersion(client, targetURL); err != nil {

	} else if tigerNeedsUpdate {

	} else if uiNeedsUpdate, err := checkUIVersion(client, targetURL); err != nil {

	} else if uiNeedsUpdate {

	} else if cliNeedsUpdate, err := checkCLIVersion(client, targetURL); err != nil {

	} else if cliNeedsUpdate {

	}

}
func checktigerVersion(client *http.Client, urlBase string) (needsUpdate bool, err error) {
	targetURL := urlBase + "/VERSION"
	if req, err := http.NewRequest("GET", targetURL, nil); err != nil {
		fmt.Printf("[!] Failed to make new request: %v\n", err)
		return false, err
	} else if resp, err := client.Do(req); err != nil {
		fmt.Printf("[!] Error client.Do: %v\n", err)
		return false, err
	} else if resp.StatusCode != 200 {
		fmt.Printf("[!] Error resp.StatusCode: %v\n", resp.StatusCode)
		return false, err
	} else {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err != nil {
			fmt.Printf("[!] Failed to read file contents: %v\n", err)
			return false, err
		} else if fileContents, err := os.ReadFile("VERSION"); err != nil {
			fmt.Printf("[!] Failed to get tiger version: %v\n", err)
			return false, err
		} else {
			remoteVersion := "v" + string(body)
			localVersion := "v" + string(fileContents)
			versionComparison := semver.Compare(localVersion, remoteVersion)
			fmt.Printf("[*] tiger Local Version:  %s\n", localVersion)
			fmt.Printf("[*] tiger Remote Version: %s\n", remoteVersion)
			if !semver.IsValid(localVersion) {
				fmt.Printf("[-] Local version isn't valid\n")
				return false, errors.New("Local version isn't valid\n")
			} else if !semver.IsValid(remoteVersion) {
				fmt.Printf("[-] Remote version isn't valid\n")
				return false, errors.New("Remote version isn't valid\n")
			} else if versionComparison == 0 {
				fmt.Printf("[+] tiger is up to date!\n")
				return false, nil
			} else if versionComparison < 0 {
				fmt.Printf("[+] tiger update available!\n")
				if semver.Major(localVersion) != semver.Major(remoteVersion) {
					fmt.Printf("[!] Major version update available. This means a major update in how tiger operates\n")
					fmt.Printf("This will require a completely new clone of tiger\n")
				} else if semver.MajorMinor(localVersion) != semver.MajorMinor(remoteVersion) {
					fmt.Printf("[!] Minor version update available. This means there has been some database updates, but not a major change to how tiger operates.\n")
					fmt.Printf("This will require doing a 'git pull' and making a new 'tiger-cli' via 'sudo make'. Then restart tiger\n")
				} else {
					fmt.Printf("[+] A patch is available. This means no database schema has changed and only bug fixes applied. This is safe to update now.\n")
					fmt.Printf("This will require doing a 'git pull' and making a new 'tiger-cli' via 'sudo make'. Then restart tiger\n")
				}
				return true, nil
			} else {
				fmt.Printf("[+] Local version is ahead of remote!\n")
				return false, nil
			}
		}
	}
}
func checkUIVersion(client *http.Client, urlBase string) (needsUpdate bool, err error) {
	targetURL := urlBase + "/tigerReactUI/src/index.js"
	if req, err := http.NewRequest("GET", targetURL, nil); err != nil {
		fmt.Printf("[!] Failed to make new request: %v\n", err)
		return false, err
	} else if resp, err := client.Do(req); err != nil {
		fmt.Printf("[!] Error client.Do: %v\n", err)
		return false, err
	} else if resp.StatusCode != 200 {
		fmt.Printf("[!] Error trying to fetch UI version resp.StatusCode: %v\n", resp.StatusCode)
		return false, err
	} else {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err != nil {
			fmt.Printf("[!] Failed to read file contents: %v\n", err)
			return false, err
		} else if fileContents, err := os.ReadFile(filepath.Join(".", "tigerReactUI", "src", "index.js")); err != nil {
			fmt.Printf("[!] Failed to get tigerReactUI version: %v\n", err)
			return false, err
		} else {
			remoteVersion := "v" + getUIVersionFromFileContents(string(body))
			localVersion := "v" + string(getUIVersionFromFileContents(string(fileContents)))
			versionComparison := semver.Compare(localVersion, remoteVersion)
			fmt.Printf("[*] UI Local Version:  %s\n", localVersion)
			fmt.Printf("[*] UI Remote Version: %s\n", remoteVersion)
			if !semver.IsValid(localVersion) {
				fmt.Printf("[-] Local version isn't valid\n")
				return false, errors.New("Local version isn't valid\n")
			} else if !semver.IsValid(remoteVersion) {
				fmt.Printf("[-] Remote version isn't valid\n")
				return false, errors.New("Remote version isn't valid\n")
			} else if versionComparison == 0 {
				fmt.Printf("[+] Your UI is up to date!\n")
				return false, nil
			} else if versionComparison < 0 {
				fmt.Printf("[+] UI update available! This is safe to update now.\n")
				return true, nil
			} else {
				fmt.Printf("[+] Local version is ahead of remote!\n")
				return false, nil
			}
		}
	}
}
func getUIVersionFromFileContents(contents string) string {
	fileLines := strings.Split(contents, "\n")
	for _, line := range fileLines {
		if strings.Contains(line, "tigerUIVersion") {
			uiVersionPieces := strings.Split(line, "=")
			uiVersion := uiVersionPieces[1]
			return uiVersion[2 : len(uiVersion)-2]
		}
	}
	return "Failed to find version"
}
func checkCLIVersion(client *http.Client, urlBase string) (needsUpdate bool, err error) {
	targetURL := urlBase + "/tiger_CLI/src/cmd/config/vars.go"
	if req, err := http.NewRequest("GET", targetURL, nil); err != nil {
		fmt.Printf("[!] Failed to make new request: %v\n", err)
		return false, err
	} else if resp, err := client.Do(req); err != nil {
		fmt.Printf("[!] Error client.Do: %v\n", err)
		return false, err
	} else if resp.StatusCode != 200 {
		fmt.Printf("[!] Error trying to fetch tiger-cli version resp.StatusCode: %v\n", resp.StatusCode)
		return false, err
	} else {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err != nil {
			fmt.Printf("[!] Failed to read file contents: %v\n", err)
			return false, err
		} else {
			remoteVersion := getCLIVersionFromFileContents(string(body))
			localVersion := config.Version
			versionComparison := semver.Compare(localVersion, remoteVersion)
			fmt.Printf("[*] tiger-cli Local Version:  %s\n", localVersion)
			fmt.Printf("[*] tiger-cli Remote Version: %s\n", remoteVersion)
			if !semver.IsValid(localVersion) {
				fmt.Printf("[-] Local version isn't valid\n")
				return false, errors.New("Local version isn't valid\n")
			} else if !semver.IsValid(remoteVersion) {
				fmt.Printf("[-] Remote version isn't valid\n")
				return false, errors.New("Remote version isn't valid\n")
			} else if versionComparison == 0 {
				fmt.Printf("[+] Your tiger-cli is up to date!\n")
				return false, nil
			} else if versionComparison < 0 {
				fmt.Printf("[+] tiger-cli update available! This is safe to update now.\n")
				fmt.Printf("Update with the following:\n")
				fmt.Printf("1. git pull\n")
				fmt.Printf("2. make")
				return true, nil
			} else {
				fmt.Printf("[+] Local version is ahead of remote!\n")
				return false, nil
			}
		}
	}
}
func getCLIVersionFromFileContents(contents string) string {
	fileLines := strings.Split(contents, "\n")
	for _, line := range fileLines {
		if strings.Contains(line, "Version =") {
			uiVersionPieces := strings.Split(line, "=")
			uiVersion := uiVersionPieces[1]
			return uiVersion[2 : len(uiVersion)-1]
		}
	}
	return "Failed to find version"
}
