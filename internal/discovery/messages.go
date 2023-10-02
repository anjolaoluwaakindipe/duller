package discovery

type Message map[string]interface{}

const (
	registerServiceMsg = "registerServericeMsg"
)

type registerServiceMessage struct {
	ServerName string `json:"serverName"`
	Path       string `json:"path"`
	Address    string `json:"address"`
}

type registryResponse struct {
	Code    int
	Message string
}
