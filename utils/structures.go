package utils

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type ServerConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type Node struct {
	ID int16
	X  float64
	Y  float64
}

type Element struct {
	ID      int16
	NodesID []int16
}
