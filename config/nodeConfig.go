package config

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/hnatekmarorg/HCI/utils"
	"html/template"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type NodeType string

const (
	Talos   NodeType = "Talos"
	Generic NodeType = "Generic"
)

// TalosConfig only applicable if node is Talos
type TalosConfig struct {
	FactoryHash string `json:"factoryHash"`
	Version     string `json:"version"`
}

// Node represents single node that can be matched
type Node struct {
	// MacAddress of the node
	MacAddress string      `json:"mac"`
	Response   string      `json:"response"`
	Type       NodeType    `json:"type"`
	Talos      TalosConfig `json:"talos,omitempty"`
}

// NodeConfig array of nodes
type NodeConfig struct {
	// MacAddress of the node
	Nodes []Node `json:"nodes"`
}

//go:embed talos.ipxe.tmpl
var talosTemplate string

type talosTemplateConfig struct {
	Kernel    string
	Initramfs string
}

func (n *Node) RenderResponse() []byte {
	if n.Type == Generic {
		return []byte(n.Response)
	} else if n.Type == Talos {
		talosPath := n.Talos.FactoryHash + "/" + n.Talos.Version
		// Download images if they aren't present
		if _, err := os.Stat(ServerConf.ImageCacheDir + "/" + n.Talos.FactoryHash); err != nil {
			if os.IsNotExist(err) {
				n.downloadTalosImages(talosPath)
			} else {
				panic(err)
			}
		}
		textTemplate, err := template.New("talos_ipxe").Parse(talosTemplate)
		if err != nil {
			panic(err)
		}
		var buffer bytes.Buffer
		bufferWriter := bufio.NewWriter(&buffer)
		kernelPath, err := url.JoinPath(ServerConf.ServerAddress, "image", n.Talos.FactoryHash, n.Talos.Version, "kernel-amd64")
		if err != nil {
			panic(err)
		}
		initRamfs, err := url.JoinPath(ServerConf.ServerAddress, "image", n.Talos.FactoryHash, n.Talos.Version, "initramfs-amd64.xz")
		err = textTemplate.Execute(bufferWriter, talosTemplateConfig{
			Kernel:    kernelPath,
			Initramfs: initRamfs,
		})
		if err != nil {
			return nil
		}
		err = bufferWriter.Flush()
		if err != nil {
			panic(err)
		}
		log.Println(buffer.String())
		return buffer.Bytes()
	}
	panic("unreachable")
}

func (n *Node) downloadTalosImages(talosPath string) {
	ch := make(chan bool)
	err := os.Mkdir(ServerConf.ImageCacheDir+"/"+n.Talos.FactoryHash, 0777)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(ServerConf.ImageCacheDir+"/"+talosPath, 0777)
	urlRoot := ServerConf.TalosFactoryServer + "/" + talosPath
	go func() {
		whatToDownload := []string{
			"/kernel-amd64",
			"/initramfs-amd64.xz",
		}
		for _, remoteFile := range whatToDownload {
			for i := 0; i < 5; i++ {
				err := utils.DownloadFile(ServerConf.ImageCacheDir+"/"+talosPath+"/"+remoteFile, urlRoot+"/"+remoteFile)
				if err == nil {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
		ch <- true
	}()
	<-ch
}

func (c *NodeConfig) GetByMac(macAddress string) *Node {
	for _, node := range c.Nodes {
		if strings.ToLower(node.MacAddress) == strings.ToLower(macAddress) {
			return &node
		}
	}
	return nil
}

func LoadNodeConfig(path string) (*NodeConfig, error) {
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

func GetConfig() *NodeConfig {
	possibleConfigLocations := []string{
		"/etc/ipxe_server.json",
		"~/.config/pxe_server.json",
	}
	if ServerConf.ConfigPath != "" {
		possibleConfigLocations = append(possibleConfigLocations, ServerConf.ConfigPath)
	}
	var config *NodeConfig
	var err error
	for _, location := range possibleConfigLocations {
		log.Printf("Loading config %s", location)
		config, err = LoadNodeConfig(location)
		if err != nil {
			log.Println("Error loading config:", err)
			continue
		}
		return config
	}
	return config
}
