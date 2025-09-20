package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// define json structure
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	IMG    string  `json:"image_url"`
}

// define commands structure
type commands struct {
	get    string
	add    bool
	update bool
	delete bool
}

func main() {
	var input Book
	var cmd commands

	// flags
	// cmd
	flag.StringVar(&cmd.get, "get", "false,all", "get all books or get book by id")
	flag.BoolVar(&cmd.add, "add", false, "add book")
	flag.BoolVar(&cmd.update, "update", false, "update book by id")
	flag.BoolVar(&cmd.delete, "delete", false, "delete book by id")

	// book info
	flag.StringVar(&input.ID, "id", "1", "ID of book")
	flag.StringVar(&input.Title, "title", "default title", "title of book")
	flag.StringVar(&input.Author, "author", "default author", "author of book")
	flag.Float64Var(&input.Price, "price", 0.0, "price of book")
	flag.StringVar(&input.IMG, "img", "", "image url of book")

	flag.Parse()

	// open json file define jsonFile as var define err as error
	jsonFile, err := os.Open("books.json")
	// check if error has occured
	if err != nil {
		//print error
		fmt.Println(err)
		err = nil
	}

	// read json file
	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		//print error
		fmt.Println(err)
		err = nil
	}

	// define books as slice of Book
	// slices are a dynamic length array, in go arrays have a fixed length
	// Book is a struct defined above
	var books []Book

	// unmarshal byteValue into books
	// unmarshal is a function that turns json data into a go struct
	json.Unmarshal(byteValue, &books)

	// closes jsonFile
	defer jsonFile.Close()

	if cmd.get[0:4] != "false" {
		if cmd.get[len(cmd.get)-3:len(cmd.get)] == "all" {
			// get all books
			for i := 0; i < len(books); i++ {
				fmt.Println(books[i])
			}
		}

		if cmd.get[len(cmd.get)-2:] == "id" {
			// get book by id
			idx, err := strconv.Atoi(input.ID)
			if err != nil {
				fmt.Println(err)
				err = nil
			} else {
				found := false
				for i, b := range books {
					if b.ID == input.ID {
						fmt.Println(b)
						found = true
						_ = i
						break
					}
				}

				if !found {
					fmt.Println("Book with ID", idx, "not found")
				}
			}
		}
	}
}
