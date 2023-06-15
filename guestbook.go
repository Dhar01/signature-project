package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

// handling errors...
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// read a slice of string from user
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

// display a form to enter signature
func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

// add a new signature to the signature file - append!
func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")

	// _, err := writer.Write([]byte(signature))
	// check(err)

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("signatures.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, signature)
	check(err)
	err = file.Close()
	check(err)

	// redirect to the home page: guestbook
	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

// Read guestbook signatures and display them with a count
func viewHandler(writer http.ResponseWriter, request *http.Request) {
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
	http.HandleFunc("/guestbook/create", createHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
