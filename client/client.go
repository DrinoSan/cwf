package client
// Package for client. I'm too tired to think of a better explanation.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	// "go.uber.org/zap"
	"cwf/entities"
)

var baseURL string
var flagLookup = map[string]string{
	"-l":  "list",
	"-lt": "list-tree",
}

// Start client and handle action types.
func StartClient() {
	baseURL = "http://" + entities.MotherShip.MotherShipIP + ":" + entities.MotherShip.MotherShipPort + "/cwf"

	if fromPipe() {
		sendContent()
	} else if getFlagValue("l") {
		listFiles()
	} else if getFlagValue("lt") {
		// listTree()
	} else if getFlagValue("d") {
		deleteFile()
	} else {
		getContent()
	}
}

// Send content to server to save it encoded in specified file.
func sendContent() {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading content from StdIn!")
		return
	}

	encStr := base64.StdEncoding.EncodeToString(content)
	body, err := json.Marshal(entities.CWFBody_t{File: os.Args[1], Content: encStr})
	if err != nil {
		fmt.Println("Error encoding data!")
		return
	}

	res, err := http.Post(baseURL + "/copy",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error sending request!")
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(responseData))
		return
	}

	fmt.Println(string(responseData))
}

// Get content of clipboard file.
func getContent() {
	res, err := http.Get(baseURL + "/get?file=" + os.Args[1])
	if err != nil {
		fmt.Println("Error getting content!")
		return
	}

	bodyEncoded, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(bodyEncoded))
		return
	}

	bodyDecoded, err := base64.StdEncoding.DecodeString(string(bodyEncoded))
	if err != nil {
		fmt.Println("Failed to decode body!")
		return
	}

	fmt.Println(string(bodyDecoded))
}

// Get a list from server.
func listFiles() {
	requestUrl := baseURL + "/list"
	if len(os.Args) > 2 {
		requestUrl += "?dir=" + os.Args[2]
	}

	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("Error sending request!")
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(responseData))
		return
	}

	fmt.Println(string(responseData))
}

// Delete a filename from server.
func deleteFile() {
	if len(os.Args) < 3 {
		fmt.Println("No filename given to delete!")
		return
	}

	client := &http.Client{}
	requestUrl := baseURL + "/delete?file=" + os.Args[2]
	req, err := http.NewRequest("DELETE", requestUrl, nil)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request!")
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(responseData))
		return
	}

	fmt.Println(string(responseData))
}

// Check if we are getting content from pipe.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode()&os.ModeCharDevice == 0
}

// Get value of a flag.
func getFlagValue(flagName string) bool {
	return flag.Lookup(flagName).Value.(flag.Getter).Get().(bool)
}
