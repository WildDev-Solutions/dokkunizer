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

		ports, err := getDokkuPortsUsed()
		if err != nil {
			fmt.Printf("Error getting ports: %v\n", err)
		}

		// fmt.Println(dokkuPorts)

		fmt.Print("Will this project use postgres? (y/n): ")
		use_pg, _ := reader.ReadString('\n')

		fmt.Print("Which port will this project use? (already used ports: " + ports + "): ")
		app_ports, _ := reader.ReadString('\n')
		fmt.Println("Creating app " + name + " with ports " + app_ports, " and postgres " + use_pg)
		return

		// if err != nil {
		// 	fmt.Printf("Error executing command: %s\n", err)
		// }

		// Print the output
		// fmt.Println(string(output))
	}
}

func getDokkuPortsUsed() (string, error) {
	// Run 'dokku apps:list' to get all apps
	appsCmd := exec.Command("dokku", "apps:list")
	appsOutput, err := appsCmd.Output()
	if err != nil {
		// fmt.Printf("Error listing apps: %v\n", err)
		return "", err
	}

	// Convert output to string and split into lines
	apps := strings.Split(string(appsOutput), "\n")

	// Initialize a slice to store port numbers
	var ports []string

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

		// Scan through the output line by line
		scanner := bufio.NewScanner(strings.NewReader(string(proxyOutput)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Look for lines containing "container port"
			if strings.Contains(line, "container port") {
				// Split the line and grab the last part (port number)
				parts := strings.Fields(line)
				if len(parts) > 0 {
					port := parts[len(parts)-1]
					ports = append(ports, port)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning output for app %s: %v\n", app, err)
		}
	}

	p := strings.Join(ports, ",")
	return p, nil
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
