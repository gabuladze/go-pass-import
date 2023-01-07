package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

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

		if record[cols["url"]] == "http://sn" {
			log.Println("Inserting secure note:", record[cols["name"]])

			path := fmt.Sprintf("%v/%v", record[cols["grouping"]], record[cols["name"]])
			cmd := exec.Command("pass", "insert", "-m", path)
			stdin, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal("Error when opening pipe", err)
			}

			if err = cmd.Start(); err != nil {
				log.Fatal("An error occured: ", err)
			}

			io.WriteString(stdin, record[cols["extra"]])
			stdin.Close()
			cmd.Wait()
		} else {
			log.Println("Inserting password:", record[cols["name"]])

			path := fmt.Sprintf("%v/%v", record[cols["grouping"]], record[cols["name"]])

			cmd := exec.Command("pass", "insert", "-m", path)
			stdin, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal("Error when opening pipe", err)
			}

			if err = cmd.Start(); err != nil {
				log.Fatal("An error occured: ", err)
			}

			content := ""
			content += fmt.Sprintf("%v\n", record[cols["password"]])
			content += fmt.Sprintf("URL: %v\n", record[cols["url"]])
			content += fmt.Sprintf("Username: %v\n", record[cols["username"]])
			content += fmt.Sprintf("Extra:\n%v\n", record[cols["extra"]])

			io.WriteString(stdin, content)
			stdin.Close()
			cmd.Wait()
		}
	}

}
