package module

type KafkaInstance struct {
	Name             string `json:"name"`
	BootstrapServers string `json:"bootstrap_servers"`
}
