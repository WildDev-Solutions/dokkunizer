package main

import (
	"bufio"
	"bytes"
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
		name = strings.TrimSpace(name)

		ports, err := getDokkuPortsUsed()
		if err != nil {
			fmt.Printf("Error getting ports: %v\n", err)
		}

		fmt.Print("Will this project use postgres? (y/n): ")
		use_pg, _ := reader.ReadString('\n')
		use_pg = strings.TrimSpace(use_pg)

		fmt.Print("Which port will this project use? (already used ports: " + ports + "): ")
		app_port, _ := reader.ReadString('\n')
		app_port = strings.TrimSpace(app_port)

		fmt.Print("What is the domain of the project: ")
		domain, _ := reader.ReadString('\n')
		domain = strings.TrimSpace(domain)

		fmt.Println("Creating app ", name, " with ports ", app_port, " and postgres ", use_pg, " on domain ", domain)

		fmt.Print("Is this correct? (y/n): ")
		correct, _ := reader.ReadString('\n')
		if correct == "y\n" {
			fmt.Println("Creating app ")
			createApp := exec.Command("dokku", "apps:create", name)
			err := createApp.Run()
			if err != nil {
				fmt.Printf("Error creating app: %v\n", err)
				break
			}

			if use_pg == "y" {
				fmt.Println("Creating postgres")
				createPostgres := exec.Command("dokku", "postgres:create", name)
				err := createPostgres.Run()
				if err != nil {
					fmt.Printf("Error creating postgres: %v\n", err)
					break
				}

				fmt.Println("Linking postgres")
				linkPostgres := exec.Command("dokku", "postgres:link", name, name)
				err = linkPostgres.Run()
				if err != nil {
					fmt.Printf("Error linking postgres: %v\n", err)
					break
				}
			}

			fmt.Println("Setting domain")
			setDomain := exec.Command("dokku", "domains:add", name, domain)
			err = setDomain.Run()
			if err != nil {
				fmt.Printf("Error setting domain: %v\n", err)
				break
			}

			fmt.Println("Running letsencrypt")
			var buffer bytes.Buffer
			letsencrypt := exec.Command("dokku", "letsencrypt:set", name, "email", "pedro.leoti.dev@gmail.com")
			letsencrypt.Stdout = &buffer
			err = letsencrypt.Run()
			fmt.Println(buffer.String())
			if err != nil {
				fmt.Printf("Error running letsencrypt set: %v\n", err)
				break
			}

			letsencryptEnable := exec.Command("dokku", "letsencrypt:enable", name)
			letsencryptEnable.Stdout = &buffer
			err = letsencryptEnable.Run()
			fmt.Println(buffer.String())
			if err != nil {
				fmt.Printf("Error running letsencrypt enable: %v\n", err)
				break
			}

			fmt.Println("Setting ports")
			setPortsHttp := exec.Command("dokku", "ports:set", name, "http:80:", app_port)
			err = setPortsHttp.Run()
			if err != nil {
				fmt.Printf("Error setting ports: %v\n", err)
				break
			}
			setPortsHttps := exec.Command("dokku", "ports:set", name, "http:433:", app_port)
			err = setPortsHttps.Run()
			if err != nil {
				fmt.Printf("Error setting ports: %v\n", err)
				break
			}

			fmt.Println("All done!")
			break

		} else {
			fmt.Println("Exiting...")
			break
		}

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
