package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) error {
	log.Println("Downloading " + url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = out.Write(content)
	if err != nil {
		return err
	}
	return nil
}
