package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func CAHandler(w http.ResponseWriter, r *http.Request) {
	serverAddress := r.URL.Query().Get("url")
	serverPort := r.URL.Query().Get("port")

	// fmt.Println("[server url]", serverAddress)
	// fmt.Println("[server port]", serverPort)

	cmd := exec.Command("bash", "-c", "/usr/bin/openssl s_client -showcerts -servername "+serverAddress+" -connect "+serverAddress+":"+serverPort+" </dev/null")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing cmd")
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

	if len(cmdSplit) > 2 {
		fmt.Fprintf(w, "-----BEGIN CERTIFICATE-----\n"+cmdSplit[1]+"\n-----END CERTIFICATE-----\n")
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
