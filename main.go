package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Give the name of the project: ")
		name, _ := reader.ReadString('\n')

		// Trim any extra spaces or newlines
		name = strings.TrimSpace(name)

		getDokkuPortsUsed()
		// fmt.Println(dokkuPorts)


		// fmt.Print("Will this project use postgres? (y/n): ")
		// use_pg, _ := reader.ReadString('\n')

		// fmt.Print("Which port will this project use? ")
		// app_ports, _ := reader.ReadString('\n')

		// if err != nil {
		// 	fmt.Printf("Error executing command: %s\n", err)
		// }

		// Print the output
		// fmt.Println(string(output))
	}
}

func getDokkuPortsUsed() {
	// Run 'dokku apps:list' to get all apps
	appsCmd := exec.Command("dokku", "apps:list")
	appsOutput, err := appsCmd.Output()
	if err != nil {
		fmt.Printf("Error listing apps: %v\n", err)
		return
	}

	// Convert output to string and split into lines
	apps := strings.Split(string(appsOutput), "\n")

	// Initialize a variable to store port mappings
	var portsMapping string

	// Iterate over each app, skipping the header
	for _, app := range apps {
		app = strings.TrimSpace(app)
		if app == "" || app == "=====> My Apps" {
			continue
		}

		// Run 'dokku proxy:ports <app>' for each app
		proxyCmd := exec.Command("dokku", "ports:list", app)
		proxyOutput, err := proxyCmd.Output()
		if err != nil {
			fmt.Printf("Error getting port mappings for app %s: %v\n", app, err)
			continue
		}

		// Append the app name and port mappings to the portsMapping string
		portsMapping += fmt.Sprintf("App: %s\n%s\n", app, string(proxyOutput))
	}

	// Print the final result with all port mappings
	fmt.Println(portsMapping)
	//
	// ///
	// cmd := exec.Command("dokku", "proxy:ports")
	// output, err := cmd.CombinedOutput()
	//
	// if err != nil {
	// 	fmt.Printf("Error executing command: %s\n", err)
	// }
	//
	// fmt.Println(string(output))
}

// If the user types "exit", break the loop and quit
// if input == "exit" {
// 	fmt.Println("Exiting...")
// 	break
// }

// // Split input into command and arguments
// parts := strings.Split(name, " ")
// cmd := exec.Command(parts[0], parts[1:]...)

// Get the command output
// output, err := cmd.CombinedOutput()
