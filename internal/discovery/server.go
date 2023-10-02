package discovery

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/mitchellh/mapstructure"
)

// Setting for the Registry Server
type RegistrySettings struct {
	REGISTRY_HOST string
	REGISTRY_PORT string
	REGISTRY_TYPE string
}

// Initiates a TCP server and accepts connections for the registry
func InitRegistryServer(rs RegistrySettings, registry *Registry) error {
	tcpServer, err := net.Listen(rs.REGISTRY_TYPE, rs.REGISTRY_HOST+":"+rs.REGISTRY_PORT)

	if err != nil {
		return err
	}

	defer tcpServer.Close()
	fmt.Println("Listening on " + rs.REGISTRY_HOST + ":" + rs.REGISTRY_PORT)

	go printRegistry(registry)
	for {
		conn, err := tcpServer.Accept()

		if err != nil {
			fmt.Printf("TCP connection error: %v \n", err)
			continue
		}

		go handleRequest(conn, registry)
	}
}

// handles every incoming request in a separate thread
func handleRequest(conn net.Conn, registry *Registry) {
	decoder := json.NewDecoder(conn)
	response := registryResponse{Code: 0}

	var msg Message
	decoder.Decode(&msg)

	handleMessage(msg, &response, registry)

	responseJson, _ := json.Marshal(response)
	conn.Write(responseJson)

	conn.Close()
}

// Useed to print out the registry every 5 seconds.
// Note: For development purpose. Remove this in production
func printRegistry(registry *Registry) {
	for {
		mutex.Lock()
		fmt.Printf("Registry Info: \n %v \n", registry)
		mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}

// handles messages comming from the clients and responds according to the
// message type e.g Register
func handleMessage(msg Message, response *registryResponse, registry *Registry) {
	// check for message type
	typeVal, typeExist := msg["type"]

	if !typeExist {
		response.Code = 1
		response.Message = "Message doesn't have 'type' property"
		return
	}

	//  get message data
	dataVal, dataExist := msg["data"]

	if !dataExist {
		response.Code = 1
		response.Message = "Message doesn't have 'data' property"
	}

	// based on the message type try and parse the message to its corresponding struct
	switch typeVal {
	case registerServiceMsg:
		var registerServerMessage registerServiceMessage
		if err := mapstructure.Decode(dataVal, &registerServerMessage); err != nil {
			response.Code = 1
			response.Message = "Data does not match type"
		}

		_, exist := registry.Services[registerServerMessage.Path]
		if exist {
			response.Code = 1
			response.Message = fmt.Sprintf(`service with path "%v" already exists in registry`, registerServerMessage.Path)
			break
		}
		mutex.Lock()
		registry.Services[registerServerMessage.Path] = ServiceInfo{timeCreated: time.Now(), serverName: registerServerMessage.ServerName, path: registerServerMessage.Path, address: registerServerMessage.Address}
		mutex.Unlock()
	default:
		fmt.Println(msg)
		response.Code = 1
		response.Message = "Invalid Message"
	}
}
