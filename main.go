package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func CAHandler(w http.ResponseWriter, r *http.Request) {
	serverAddress := r.URL.Query().Get("url")
	serverPort := r.URL.Query().Get("port")

	// fmt.Println("[server url]", serverAddress)
	// fmt.Println("[server port]", serverPort)

	// Check if valid domain name
	// https://stackoverflow.com/questions/7930751/regexp-for-subdomain
	match, _ := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9.-]+[a-zA-Z0-9]$", serverAddress)

	if !match {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Domain name invalid. \""+serverAddress+"\" is not a valid domain name."+func() string {

			if strings.Contains(serverAddress, "poweroff") ||
				strings.Contains(serverAddress, "reboot") ||
				strings.Contains(serverAddress, "sudo") {
				return " Nice try though."
			} else {
				return ""
			}
		}())
		return
	}

	// Try to parse port if valid integer
	_, serverPortConvErr := strconv.Atoi(serverPort)

	if serverPortConvErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Port invalid. \""+serverPort+"\" is not a valid port.")

		return
	}

	cmd := exec.Command("bash", "-c", "/usr/bin/openssl s_client -showcerts -servername "+serverAddress+" -connect "+serverAddress+":"+serverPort+" </dev/null")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing cmd", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Exec error. Maybe domain name invalid?")

		return
	}

	cmdOutput := string(stdout)

	// Split by -----BEGIN CERTIFICATE-----
	cmdSplit := strings.Split(cmdOutput, "-----BEGIN CERTIFICATE-----")

	if len(cmdSplit) > 1 {
		cmdSplit = cmdSplit[1:]
	}

	// Slice -----END CERTIFICATE-----
	for i, splitStr := range cmdSplit {
		endCertSubstringIndex := strings.Index(splitStr, "-----END CERTIFICATE-----")

		cmdSplit[i] = strings.TrimSpace(splitStr[:endCertSubstringIndex])
	}

	// for i, splitStr := range cmdSplit {
	// 	fmt.Println("[", i, "]", splitStr)
	// }

	if len(cmdSplit) > 1 {
		fmt.Fprintf(w, "-----BEGIN CERTIFICATE-----\n"+cmdSplit[1]+"\n-----END CERTIFICATE-----\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("Certificate index out of range"))

		return
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", CAHandler)

	fmt.Println("Listening on port " + os.Getenv("SERVER_PORT"))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), r))

}
