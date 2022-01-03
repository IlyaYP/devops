package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

//SendBuf("http://localhost:8080/update/", &buf)
func SendBuf(endpoint string, buf io.Reader) error {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, buf)
	if err != nil {
		log.Fatal(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	//request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	//request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Статус-код ", response.Status)
	return nil
}

func send(endpoint string) {
	data := `Hi, how are you? русский текст`
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//fmt.Println(string(body))
	_ = string(body)
}

func Receive(endpoint string) string {
	data := `Hi, how are you? русский текст`
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBufferString(data))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//fmt.Println(string(body))
	return string(body)
}
