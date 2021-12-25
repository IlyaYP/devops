package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func Send(endpoint string) {
	// Имя метрики: "Alloc", тип: gauge
	// http://localhost:8080/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	//endpoint := "http://localhost:8080/update/gauge/Alloc/4000000000.0001"
	//endpoint := "http://localhost:8080/asedffggf"
	data := `Hi, how are you? русский текст`
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
