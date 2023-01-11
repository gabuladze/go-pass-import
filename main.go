package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"

	"golang.org/x/net/publicsuffix"
)

func getTldPlusOne(fullUrl string) string {
	result := "*"

	u, err := url.Parse(fullUrl)
	if err != nil {
		return result
	}

	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(u.Host)
	if err != nil {
		return result
	}

	result = fmt.Sprintf("%v://*.%v/*", u.Scheme, tldPlusOne)
	return result
}

func main() {
	cols := map[string]int{
		"url":      0,
		"username": 1,
		"password": 2,
		"totp":     3,
		"extra":    4,
		"name":     5,
		"grouping": 6,
		"fav":      7,
	}

	filePath := "./in.csv"
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[cols["fav"]] == "fav" {
			// skip header
			continue
		}

		path := fmt.Sprintf("%v/%v", record[cols["grouping"]], record[cols["name"]])

		cmd := exec.Command("pass", "insert", "-m", path)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal("Error when opening pipe", err)
		}

		if err = cmd.Start(); err != nil {
			log.Fatal("An error occured: ", err)
		}

		if record[cols["url"]] == "http://sn" {
			log.Println("Inserting secure note:", path)

			io.WriteString(stdin, record[cols["extra"]])
		} else {
			log.Println("Inserting password:", path)

			url := getTldPlusOne(record[cols["url"]])
			content := ""
			content += fmt.Sprintf("%v\n", record[cols["password"]])
			content += fmt.Sprintf("URL: %v\n", url)
			content += fmt.Sprintf("Username: %v\n", record[cols["username"]])
			content += fmt.Sprintf("Extra:\n%v\n", record[cols["extra"]])

			io.WriteString(stdin, content)
		}

		stdin.Close()
		cmd.Wait()
	}

	log.Println("Done!")
}
