package main

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"math/rand"
	"time"
	"flag"
	"encoding/csv"
	"log"
)

var t int
var wordsPath = "./words.csv"

func init() {
	flag.IntVar(&t, "t", 10, "制限時間")
}

func main() {
	flag.Parse()

	score := 0
	fmt.Println("TIME LIMIT is", t)
	timeout := time.After(time.Second * time.Duration(t))
	for sign := true; sign == true; {
		word := RandomWord()
		fmt.Println(word)
		c := imp(os.Stdin)
		select {
		case right := <-c:
			if right == word {
				fmt.Println("OK!")
				score += len(word)
			} else {
				fmt.Println("NG!")
			}
		case <-timeout:
			fmt.Println("It's time!")
			sign = false
		}
	}	
	fmt.Println("GAME OVER! Your score is", score * 10)
}

func RandomWord() (word string) {
	words, err := readCSV(wordsPath)

	if err != nil {
		panic(nil)
	}
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(words))
	return words[num]
}

func imp(r io.Reader) <-chan string {
	wordCh := make(chan string, 1)
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	wordCh <- scanner.Text()
	return wordCh
}

func readCSV(path string) ([]string, error) {
	var ret []string
	var row []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	csvFile := csv.NewReader(file)
	csvFile.TrimLeadingSpace = true

	for {
		row, err = csvFile.Read()
		if err != nil {
			break
		}
		ret = append(ret, row...)
	}
	return ret, nil
}