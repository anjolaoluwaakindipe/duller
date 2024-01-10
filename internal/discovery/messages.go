package discovery

type Message struct {
	Type string
	Data interface{}
}

const (
	registerServiceMsg = "registerServiceMsg"
	getAddressMsg      = "getAddressMsg"
)

type HeartBeatMessage struct {
	ServiceId string `json:"serviceId"`
	Path      string `json:"path"`
	IP        string `json:"ip"`
	Port      stirng `json:"port"`
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
