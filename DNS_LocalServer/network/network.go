package network

type NetConfig struct {
	BindAdress string `json:"bind_address"`
	Port       int    `json:"port"`
}

const (
	PACKET_GET_DNS = 48
	PACKET_REG_DNS = 49
)
