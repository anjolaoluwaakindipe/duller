package discovery

type Message struct {
	Type string
	Data interface{}
}

const (
	registerServiceMsg = "registerServiceMsg"
	getAddressMsg      = "getAddressMsg"
)

type RegisterServiceMessage struct {
	ServiceName string `json:"serverName"`
	Path        string `json:"path"`
	Address     string `json:"address"`
}

type GetAddressMessage struct {
	Path string `json:"path"`
}

type RegistryResponse struct {
	Code    int
	Message string
	Data    interface{}
}

type GetServiceResponse struct {
	IP        string
	Port      string
	ServiceId string
}
