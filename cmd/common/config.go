package common

type MySQL struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"passwd"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

type Node struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

type Client struct {
	Timeout uint32          `json:"timeout"`
	SetName string          `json:"setname"`
	Nodes   map[string]Node `json:"nodes"`
}

type Config struct {
	MySQL    MySQL    `json:"mysql"`
	BlackIPs []string `json:"blackip"`
	Client   Client   `json:"client"`
}

var Conf Config
