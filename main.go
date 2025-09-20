package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

// define json structure
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	IMG    string  `json:"image_url"`
}

func main() {
	var input Book

	// flags
	commandType := flag.String("command", "read", "type of command to run")
	flag.StringVar(&input.ID, "ID", "1", "ID of book")

	flag.Parse()

	println("Command Type:", *commandType)

	// open json file define jsonFile as var define err as error
	jsonFile, err := os.Open("books.json")
	// check if error has occured
	if err != nil {
		//print error
		fmt.Println(err)
		err = nil
	}

	// if no error print message
	fmt.Println("Successfully Opened books.json")

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
	fmt.Print(books)

	// closes jsonFile
	defer jsonFile.Close()

}
