package omikuji

import (
	"time"
	"net/http"
	"math/rand"
	"encoding/json"
	"log"
	"fmt"
)

type Omikuji struct {
	time time.Time
}

type Result struct {
	Type string `json:"結果"`
}

var types = []string{
	"大吉",
	"吉",
	"中吉",
	"小吉",
	"末吉",
	"凶",
	"大凶",
}

func New(t time.Time) *Omikuji {
	rand.Seed(t.UnixNano())
	return &Omikuji{time: t}
}

func (o *Omikuji) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := draw(o.time)
	fmt.Println("結果", result.Type)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println("err: ", err)
	}
}

func draw(t time.Time) *Result {
	yd := t.YearDay()

	if yd == 1 || yd == 2 || yd == 3 {
		return &Result{Type: types[0]}
	}
	return &Result{Type: types[rand.Intn(len(types))]}
}
