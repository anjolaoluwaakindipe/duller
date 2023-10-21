package discovery

type Message struct {
	Type string
	Data interface{}
}

const (
	registerServiceMsg = "registerServiceMsg"
)

type RegisterServiceMessage struct {
	ServerName string `json:"serverName"`
	Path       string `json:"path"`
	Address    string `json:"address"`
}

type GetAddressMessage struct {
	Path string `json:"path"`
}

type RegistryResponse struct {
	Code    int
	Message string
	Data    interface{}
}
