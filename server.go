package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v11"
	"log"
	"net/http"
	"os"
)

/*
	http server that serves custom configuration per node based on json config
*/

// Node represents single node that can be matched
type Node struct {
	// MacAddress of the node
	MacAddress string `json:"mac"`
	Response   string `json:"response"`
}

// NodeConfig array of nodes
type NodeConfig struct {
	// MacAddress of the node
	Nodes []Node `json:"nodes"`
}

type ServerConfig struct {
	Port          int    `env:"PORT" envDefault:"80"`
	ListenAddress string `env:"ADDRESS" envDefault:"0.0.0.0"`
	ConfigPath    string `env:"CONFIG_PATH"`
}

var config ServerConfig

func (c *NodeConfig) getByMac(macAddress string) *Node {
	for _, node := range c.Nodes {
		if node.MacAddress == macAddress {
			return &node
		}
	}
	return nil
}

func loadNodeConfig(path string) (*NodeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config NodeConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func getConfig() *NodeConfig {
	possibleConfigLocations := []string{
		"/etc/pxe_server.json",
		"~/.config/pxe_server.json",
	}
	if config.ConfigPath != "" {
		possibleConfigLocations = append(possibleConfigLocations, config.ConfigPath)
	}
	var config *NodeConfig
	var err error
	for _, location := range possibleConfigLocations {
		log.Printf("Loading config %s", possibleConfigLocations)
		config, err = loadNodeConfig(location)
		if err != nil {
			log.Println("Error loading config:", err)
		}
	}
	return config
}

func serveIPXE(w http.ResponseWriter, r *http.Request) {
	// Load config each time so we can reload it through tf without restarting the server
	config := getConfig()
	if config == nil {
		http.Error(w, "No config found", http.StatusInternalServerError)
		return
	}
	macAddress := r.URL.Query().Get("mac")
	if macAddress == "" {
		http.Error(w, "Missing MAC address", http.StatusBadRequest)
		return
	}
	node := config.getByMac(macAddress)
	if node == nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}
	_, err := w.Write([]byte(node.Response))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", serveIPXE)
	addresWithPort := fmt.Sprintf("%s:%d", config.ListenAddress, config.Port)
	log.Println("Listening on", addresWithPort)
	err = http.ListenAndServe(addresWithPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
