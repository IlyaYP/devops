package internal

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SendBufRetry(endpoint string, buf *bytes.Buffer) error {
	var err error
	for i := 0; i < 3; i++ {
		if err = SendBuf(endpoint, buf); err != nil {
			log.Println(err)
			log.Println("once again")
			time.Sleep(time.Duration(1) * time.Second)
		} else {
			return nil
		}
	}
	return err
}

//SendBuf("http://localhost:8080/update/", &buf)
func SendBuf(endpoint string, buf *bytes.Buffer) error {

	var b bytes.Buffer

	gz, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return fmt.Errorf("failed init compress writer: %v", err)
	}

	_, err = buf.WriteTo(gz) //gz.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = gz.Flush()
	if err != nil {
		return fmt.Errorf("failed Flush data: %v", err)
	}
	err = gz.Close()
	if err != nil {
		return fmt.Errorf("failed compress data: %v", err)
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, &b)
	if err != nil {
		log.Println(err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")
	//request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	//request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()
	//log.Println("Статус-код ", response.Status)
	return nil
}

func Send(endpoint string) {
	data := `Hi, how are you? русский текст`
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
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
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body))
	return string(body)
}
