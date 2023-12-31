package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/mitchellh/mapstructure"
)

// RegistrySettings holds configuration for the Registry Server
type RegistrySettings struct {
	REGISTRY_HOST      string
	REGISTRY_PORT      string
	REGISTRY_TYPE      string
	HEARTBEAT_INTERVAL time.Duration
}

// InitRegistryServer initiates a TCP server and accepts connections for the registry
func InitRegistryServer(balancer balancer.LoadBalancer, rs RegistrySettings, registry Registry, ctx context.Context, status chan<- bool) error {
	tcpServer, err := net.Listen(rs.REGISTRY_TYPE, rs.REGISTRY_HOST+":"+rs.REGISTRY_PORT)
	fmt.Println(err)
	if err != nil {
		if status != nil {
			status <- false
			close(status)
		}
		return err
	}

	defer tcpServer.Close()
	log.Println("Listening on " + rs.REGISTRY_HOST + ":" + rs.REGISTRY_PORT)

	// Registry gorountines
	go printRegistry(registry, ctx)
	go registry.RefreshRegistry(rs.HEARTBEAT_INTERVAL, ctx)

	status <- true

	for {
		select {
		case <-ctx.Done():
			if status != nil {
				status <- false
				close(status)
			}
			return nil
		default:
			conn, err := tcpServer.Accept()
			if err != nil {
				fmt.Printf("TCP connection error: %v \n", err)
				continue
			}
			go HandleRequest(conn, registry, balancer)
		}
	}
}

// HandleRequest handles every incoming request. Intended to be used in a seaparate goroutine.
func HandleRequest(conn net.Conn, registry Registry, balancer balancer.LoadBalancer) {
	decoder := json.NewDecoder(conn)
	response := RegistryResponse{Code: 0}

	var msg Message
	err := decoder.Decode(&msg)

	if err != nil {
		response.Code = 1
		response.Message = "Invalid Message structure"
	} else {
		handleMessage(msg, &response, registry, balancer)
	}

	responseJson, _ := json.Marshal(response)
	if _, err := conn.Write(responseJson); err != nil {
		slog.Warn("Could not Send response to connection")
	}

	if err := conn.Close(); err != nil {
		slog.Warn("Could not close connection for client")
	}
}

// Useed to print out the registry every 5 seconds.
// Note: For development purpose. Remove this in production
func printRegistry(registry service.Registry, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Registry Info: \n %v \n", registry)
			time.Sleep(5 * time.Second)
		}
	}
}

// handleMessage handles incoming messages from the clients and responds according to the
// message type e.g re
func handleMessage(msg Message, response *RegistryResponse, registry service.Registry, balancer balancer.LoadBalancer) {
	// check for message type
	msgType := msg.Type

	//  get message data
	data := msg.Data

	// based on the message type try and parse the message to its corresponding struct
	switch msgType {
	case registerServiceMsg:
		var message RegisterServiceMessage
		if err := mapstructure.Decode(data, &message); err != nil {
			response.Code = 1
			response.Message = "Data does not match specified type"
		}
		if err := registry.RegisterService(service.ServiceInfo{ServiceId: message.ServiceName, Path: message.Path}); err != nil {
			response.Code = 1
			response.Message = fmt.Sprintf(`service with path "%v" already exists in registry`, message.Path)
			break
		}

	case getAddressMsg:
		var message GetAddressMessage
		if err := mapstructure.Decode(data, &message); err != nil {
			response.Code = 1
			response.Message = "Data does not match specified type"
		}
		serviceInfo, err := balancer.GetNextService(message.Path)
		if err != nil {
			response.Code = 1
			response.Message = fmt.Sprintf("Error occured while getting service: %v", err.Error())
			break
		}

		dataMap := make(map[string]interface{})
		dataMap["address"] = "http://" + serviceInfo.IP + ":" + serviceInfo.Port
		response.Data = dataMap

	default:
		fmt.Println(msg)
		response.Code = 1
		response.Message = "Invalid Message Type"
	}
}
