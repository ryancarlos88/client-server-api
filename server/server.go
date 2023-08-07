package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func main(){
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}
const (
	apiCallTimeout = time.Millisecond * 200
	dbWriteTimeout = time.Millisecond * 10
	apiUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)
func cotacaoHandler(w http.ResponseWriter, r *http.Request){
	
	apiCtx, cancel := context.WithTimeout(context.Background(), apiCallTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, apiUrl, nil)
	checkErr(err)

	res, err := http.DefaultClient.Do(req)
	checkErr(err)

	bodyBytes, err := io.ReadAll(res.Body)
	checkErr(err)

	var c Cotacao

	err = json.Unmarshal(bodyBytes, &c)
	checkErr(err)



	w.Write([]byte(c.Usdbrl.Bid))

	
}


type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func checkErr(e error){
	if e != nil {
		panic(e)
	}
}