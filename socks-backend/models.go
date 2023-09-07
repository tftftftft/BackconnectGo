package main

type ProxyInfo struct {
	ServerIP            string `json:"ServerIP"`
	ServerListeningPort string `json:"ServerListeningPort"`
	ProxyIP             string `json:"ProxyIP"`
	CountryCode         string `json:"CountryCode"`
	Region              string `json:"Region"`
	City                string `json:"City"`
	Zip                 string `json:"Zip"`
	Mobile              bool   `json:"Mobile"`
	Proxy               bool   `json:"Proxy"`
	Hosting             bool   `json:"Hosting"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
