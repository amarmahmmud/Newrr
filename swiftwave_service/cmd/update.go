package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var updateCmd = &cobra.Command{
	Use:    "update",
	Short:  "Update Swiftwave to the latest minor patch version",
	Long:   `Update Swiftwave to the latest minor patch version`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		isUpdated := false
		distro, err := detectDistro()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if strings.Contains(distro, "debian") {
			isUpdated, err = updateDebianPackage("swiftwave")
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		} else if strings.Contains(distro, "redhat") {
			isUpdated, err = updateRedHatPackage("swiftwave")
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		} else {
			fmt.Println("Error: unknown distribution")
			return
		}

		if isUpdated {
			fmt.Println("Swiftwave has been updated successfully")
			fmt.Println("Trying to restart the service...")
			_ = exec.Command("systemctl", "restart", "swiftwave.service").Run()
			os.Exit(0)
		} else {
			fmt.Println("Swiftwave is already up to date")
			os.Exit(0)
		}
	},
}

func detectDistro() (string, error) {
	out, err := exec.Command("bash", "-c", "cat /etc/*release").Output()
	if err != nil {
		return "", err
	}

	output := strings.ToLower(string(out))
	if strings.Contains(output, "debian") || strings.Contains(output, "ubuntu") {
		return "debian", nil
	} else if strings.Contains(output, "redhat") || strings.Contains(output, "centos") || strings.Contains(output, "fedora") {
		return "redhat", nil
	}

	return "", fmt.Errorf("unknown distribution")
}

func updateDebianPackage(packageName string) (bool, error) {
	// run apt update first
	output, err := exec.Command("apt", "update", "-y").Output()
	fmt.Println(string(output))
	if err != nil {
		return false, err
	}
	output, err = exec.Command("apt", "install", "--only-upgrade", packageName).Output()
	fmt.Println(string(output))
	if err != nil {
		return false, err
	}
	// check if the package is already up to date
	line := "swiftwave is already the newest version"
	if strings.Contains(string(output), line) {
		return false, nil
	}
	return true, nil
}

func updateRedHatPackage(packageName string) (bool, error) {
	output, err := exec.Command("dnf", "update", packageName).Output()
	if err != nil {
		return false, err
	}
	if strings.Contains(string(output), "Nothing to do.") {
		return false, nil
	}
	return true, nil
}
