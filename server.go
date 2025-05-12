package main

import (
	_ "embed"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/hnatekmarorg/HCI/config"
	"log"
	"net/http"
)

/*
http server that serves custom configuration per node based on json config
*/
func serveIPXE(w http.ResponseWriter, r *http.Request) {
	// Load config each time so we can reload it through tf without restarting the server
	loadedConfig := config.GetConfig()
	if loadedConfig == nil {
		http.Error(w, "No config found", http.StatusInternalServerError)
		return
	}
	macAddress := r.URL.Query().Get("mac")
	if macAddress == "" {
		http.Error(w, "Missing MAC address", http.StatusBadRequest)
		return
	}
	node := loadedConfig.GetByMac(macAddress)
	if node == nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}
	_, err := w.Write(node.RenderResponse())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	err := env.Parse(&config.ServerConf)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ipxe", serveIPXE)
	fileServer := http.FileServer(http.Dir(config.ServerConf.ImageCacheDir))
	http.Handle("/image/", http.StripPrefix("/image/", fileServer))
	addressWithPort := fmt.Sprintf("%s:%d", config.ServerConf.ListenAddress, config.ServerConf.Port)
	log.Println("Listening on", addressWithPort)
	err = http.ListenAndServe(addressWithPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
