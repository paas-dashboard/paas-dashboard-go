package module

type PulsarInstance struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	WebPort int    `json:"web_port"`
	TcpPort int    `json:"tcp_port"`
}
