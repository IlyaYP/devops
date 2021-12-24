package handlers

import (
	"fmt"
	"net/http"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request URL:", r.URL)
	////fmt.Println("request Headers:", r.Header)
	//body, _ := io.ReadAll(r.Body)
	//fmt.Println("request Body:", string(body))

	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

