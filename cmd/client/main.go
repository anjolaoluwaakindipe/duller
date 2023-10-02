package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:9876")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	message := struct {
		Type string `json:"type"`

		Data struct {
			ServerName string `json:"serverName"`
			Path       string `json:"path"`
			Address    string `json:"address"`
		} `json:"data"`
	}{
		Type: "registerServiceMsg",
		Data: struct {
			ServerName string `json:"serverName"`
			Path       string `json:"path"`
			Address    string `json:"address"`
		}{

			ServerName: "hello",
			Path:       "tan",
			Address:    "http://localhost:3000/",
		},
	}

	jsonMessage, _ := json.Marshal(message)

	for {
		_, err = conn.Write(jsonMessage)
		if err != nil {
			fmt.Println("Error sending JSON message to server:", err)
			return
		}

		decoder, _ := io.ReadAll(conn)

		fmt.Println("Sent JSON message to server:", string(decoder))
	}
}
