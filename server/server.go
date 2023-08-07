package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer() *Server {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		checkErr(err)
	}
	return &Server{db}
}

func main() {
	log.Println("iniciando server")
	s := NewServer()

	s.db.Table("cotacoes").AutoMigrate(&Cotacao{})

	http.HandleFunc("/cotacao", s.cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

const (
	apiCallTimeout = time.Millisecond * 200
	dbWriteTimeout = time.Millisecond * 10
	apiUrl         = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

func (s *Server) cotacaoHandler(w http.ResponseWriter, r *http.Request) {

	apiCtx, cancel := context.WithTimeout(context.Background(), apiCallTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, apiUrl, nil)
	checkErr(err)

	res, err := http.DefaultClient.Do(req)
	checkErr(err)

	bodyBytes, err := io.ReadAll(res.Body)
	checkErr(err)

	var a ApiResponse

	err = json.Unmarshal(bodyBytes, &a)
	checkErr(err)

	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()

	s.db.WithContext(dbCtx).Table("cotacoes").Create(&a.Cotacao)
	
	w.Write([]byte(a.Cotacao.Bid))
}

type Cotacao struct {
	ID         int    `gorm:"primaryKey"`
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
}
type ApiResponse struct {
	Cotacao Cotacao `json:"USDBRL"`
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
