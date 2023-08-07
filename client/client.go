package main

import (
	"context"
	"fmt"
	"io"
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

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	r, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	f.Write([]byte(fmt.Sprintf("Dólar: {%v}", string(r))))


}