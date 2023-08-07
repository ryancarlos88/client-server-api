package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	requestTimeout = time.Millisecond * 300
	requestURL = "http://localhost:8080/cotacao"
	fileName = "cotacao.txt"
)
func main(){
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	checkErr(err)

	res, err := http.DefaultClient.Do(req)
	checkErr(err)
	defer res.Body.Close()

	f, err := os.Create(fileName)
	checkErr(err)

	r, err := io.ReadAll(res.Body)
	checkErr(err)

	f.Write([]byte(fmt.Sprintf("Dólar: {%v}", string(r))))
}

func checkErr(e error) {
	if e != nil {
		log.Printf("erro encontrado: %v", e.Error())
		panic(e)
	}
}