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
	const (
		URL int = iota
		USERNAME
		PASSWORD
		TOTP
		EXTRA
		NAME
		GROUPING
		FAV
	)

	filePath := os.Args[1]
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

		if record[FAV] == "fav" {
			// skip header
			continue
		}

		path := fmt.Sprintf("%v/%v", record[GROUPING], record[NAME])

		cmd := exec.Command("pass", "insert", "-m", path)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal("Error when opening pipe", err)
		}

		if err = cmd.Start(); err != nil {
			log.Fatal("An error occured: ", err)
		}

		if record[URL] == "http://sn" {
			log.Println("Inserting secure note:", path)

			io.WriteString(stdin, record[EXTRA])
		} else {
			log.Println("Inserting password:", path)

			url := getTldPlusOne(record[URL])
			content := ""
			content += fmt.Sprintf("%v\n", record[PASSWORD])
			content += fmt.Sprintf("URL: %v\n", url)
			content += fmt.Sprintf("Username: %v\n", record[USERNAME])
			content += fmt.Sprintf("Extra:\n%v\n", record[EXTRA])

			io.WriteString(stdin, content)
		}

		stdin.Close()
		cmd.Wait()
	}

	log.Println("Done!")
}
