package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-socks5"
)

type InitialProxyData struct {
	ProxyIP      string `json:"ProxyIP"`
	ProxyPort    string `json:"ProxyPort"`
	UserID       string `json:"UserID"`
	BuildVersion string `json:"BuildVersion"`
}

// IPs to be filled in at compile-time
var ipCSV string // Example default IPs
var ipList []string

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getPublicIP() (string, error) {
	res, err := http.Get("https://api4.ipext.org")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	ip := string(bodyBytes)
	return ip, nil
}

// functions that sends MachineData every 15 minutes
func sendMachineDataRoutine(proxy *InitialProxyData, ipList []string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		success := false
		for _, ip := range ipList {
			jsonproxy, err := json.Marshal(proxy)
			if err != nil {
				log.Printf("Error marshalling proxy data: %s\n", err)
				continue
			}

			resp, err := http.Post(fmt.Sprintf("http://%s:25000/route", ip), "application/json", bytes.NewBuffer(jsonproxy))
			if err != nil {
				log.Printf("Error sending HTTP request to %s: %s\n", ip, err)
				continue
			}

			if resp.StatusCode == 200 {
				log.Printf("Successfully sent machine data to %s\n", ip)
				success = true
				break
			} else {
				log.Printf("Received non-200 response code %d from %s\n", resp.StatusCode, ip)
			}
		}

		if !success {
			log.Println("Failed to send data to all IPs")
		}
	}
}

// function returns random port number from 20000 to 30000 range
func getRandomPort() string {
	min := 20000
	max := 30000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// AddFirewallRule allows the app to receive inbound traffic on a specified port
func AddFirewallRule(port string) error {
	// Find the path of the running executable
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	// Prepare the command
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name=MyAppRule",
		"dir=in",
		"action=allow",
		"protocol=TCP",
		"localport="+port,
		"program="+exePath,
		"enable=yes",
	)

	// Run the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	var proxy InitialProxyData

	proxy.ProxyIP, _ = getPublicIP()
	proxy.ProxyPort = getRandomPort()
	proxy.BuildVersion = "1.0"
	proxy.UserID = "test-user"

	if err := AddFirewallRule(proxy.ProxyPort); err != nil {
		log.Fatal(err)
	}
	log.Println("Proxy Info: ", proxy)

	ipList = strings.Split(ipCSV, ",")

	go sendMachineDataRoutine(&proxy, ipList)

	// Create a SOCKS5 server
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "0.0.0.0:"+proxy.ProxyPort); err != nil {
		panic(err)
	}
}
