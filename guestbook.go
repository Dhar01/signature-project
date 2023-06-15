package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	// if the file doesn't exist
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	// ensure the file is closed
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())

	return lines
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	// displaying the loaded signatures
	signatures := getStrings("signatures.txt")
	html, err := template.ParseFiles("view.html")
	check(err)

	guestbook := Guestbook{
		SignatureCount: len(signatures),
		Signatures:     signatures,
	}

	err = html.Execute(writer, guestbook)
	check(err)
}

func main() {
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/new", newHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
