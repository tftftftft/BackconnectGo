package main

type InitialProxyData struct {
	ProxyIP      string `json:"ProxyIP"`
	ProxyPort    string `json:"ProxyPort"`
	UserID       string `json:"UserID"`
	BuildVersion string `json:"BuildVersion"`
}
type IPInfo struct {
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`
	Continent     string `json:"continent"`
	ContinentCode string `json:"continentCode"`
	Country       string `json:"country"`
	CountryCode   string `json:"countryCode"`
	Region        string `json:"region"`
	RegionName    string `json:"regionName"`
	City          string `json:"city"`
	Zip           string `json:"zip"`
	Timezone      string `json:"timezone"`
	ISP           string `json:"isp"`
	Org           string `json:"org"`
	ASName        string `json:"asname"`
	Mobile        bool   `json:"mobile"`
	Proxy         bool   `json:"proxy"`
	Hosting       bool   `json:"hosting"`
}

type BotInfo struct {
	InitialProxyData
	IPInfo
}
