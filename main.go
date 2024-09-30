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
		fmt.Print("What is the name of the project: ")
		name, _ := reader.ReadString('\n')

		// Trim any extra spaces or newlines
		name = strings.TrimSpace(name)

		ports, err := getDokkuPortsUsed()
		if err != nil {
			fmt.Printf("Error getting ports: %v\n", err)
		}

		fmt.Print("Will this project use postgres? (y/n): ")
		use_pg, _ := reader.ReadString('\n')

		fmt.Print("Which port will this project use? (already used ports: " + ports + "): ")
		app_ports, _ := reader.ReadString('\n')

		fmt.Println("Creating app ", name, " with ports ", app_ports, " and postgres ", use_pg)
		fmt.Println("Creating app ", name, " with pg ", use_pg, " and ports ", ports)

		fmt.Print("Is this correct? (y/n): ")
		correct, _ := reader.ReadString('\n')
		if correct == "y\n" {
			fmt.Println("Creating app ")
		} else {
			fmt.Println("Exiting...")
			break
		}

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
		fmt.Printf("Error listing apps: %v\n", err)
		return "", err
	}

	// Convert output to string and split into lines
	apps := strings.Split(string(appsOutput), "\n")

	// Initialize a slice to store port numbers
	var ports []string
	portMap := make(map[string]bool) // To track unique ports

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
			// fmt.Printf("Error getting port mappings for app %s: %v\n", app, err)
			continue
		}

		// Scan through the output line by line
		scanner := bufio.NewScanner(strings.NewReader(string(proxyOutput)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Look for lines containing "https" or "http" to get the port number
			if strings.HasPrefix(line, "https") || strings.HasPrefix(line, "http") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					// Get the port number (last item in the line)
					port := parts[2] // Host port is the 3rd element in the line
					// Check if port is already in the map
					if !portMap[port] {
						ports = append(ports, port)
						portMap[port] = true // Mark port as seen
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning output for app %s: %v\n", app, err)
		}
	}

	p := strings.Join(ports, ", ")
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
