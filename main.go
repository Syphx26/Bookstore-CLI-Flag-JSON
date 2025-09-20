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
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  string `json:"price"`
	IMG    string `json:"image_url"`
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
	flag.StringVar(&cmd.get, "get", "", "get all books or get book by id")
	flag.BoolVar(&cmd.add, "add", false, "add book")
	flag.BoolVar(&cmd.update, "update", false, "update book by id")
	flag.BoolVar(&cmd.delete, "delete", false, "delete book by id")

	// book info
	flag.StringVar(&input.ID, "id", "", "ID of book")
	flag.StringVar(&input.Title, "title", "", "title of book")
	flag.StringVar(&input.Author, "author", "", "author of book")
	flag.StringVar(&input.Price, "price", "", "price of book")
	flag.StringVar(&input.IMG, "image_url", "", "image url of book")

	flag.Parse()

	// open json file for reading
	jsonFile, err := os.Open("books.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// read json file
	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		//print error
		fmt.Println(err)
		err = nil
	}

	// ensure file is closed as soon as we've finished reading it
	jsonFile.Close()

	// define books as slice of Book
	// slices are a dynamic length array, in go arrays have a fixed length
	// Book is a struct defined above
	var books []Book

	// unmarshal byteValue into books
	// unmarshal is a function that turns json data into a go struct
	if err := json.Unmarshal(byteValue, &books); err != nil {
		fmt.Println("failed to parse JSON:", err)
		return
	}

	// get all books
	if cmd.get == "all" {
		for i := 0; i < len(books); i++ {
			fmt.Println(books[i])
		}
	}

	// get book by id
	if cmd.get == "id" {
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

	// add book
	if cmd.add {
		for i, b := range books {
			if b.ID == input.ID {
				fmt.Println("Book with ID", input.ID, "already exists")
				_ = i
				return
			}
		}

		if input.ID == "" || input.Title == "" || input.Author == "" || input.Price == "" {
			fmt.Println("Please provide all book details: id, title, author, price")
			return
		}
		// append the new book to the slice
		books = append(books, input)

		// marshal the entire slice back to JSON (use indent for readability)
		out, err := json.MarshalIndent(books, "", "  ")
		if err != nil {
			fmt.Println("failed to marshal books:", err)
			return
		}

		// write the updated JSON back to the file atomically
		tmpFile, err := os.CreateTemp(".", "books.json.tmp.*")
		if err != nil {
			fmt.Println("failed to create temp file:", err)
			return
		}
		tmpName := tmpFile.Name()

		// ensure temp file is removed on error
		removeTmp := func() {
			tmpFile.Close()
			os.Remove(tmpName)
		}

		if _, err := tmpFile.Write(out); err != nil {
			fmt.Println("failed to write temp file:", err)
			removeTmp()
			return
		}

		if err := tmpFile.Sync(); err != nil {
			fmt.Println("failed to sync temp file:", err)
			removeTmp()
			return
		}

		if err := tmpFile.Close(); err != nil {
			fmt.Println("failed to close temp file:", err)
			removeTmp()
			return
		}

		// set permissions to 0644
		if err := os.Chmod(tmpName, 0644); err != nil {
			fmt.Println("failed to chmod temp file:", err)
			removeTmp()
			return
		}

		// rename temp file into place (atomic on most OSes)
		if err := os.Rename(tmpName, "books.json"); err != nil {
			fmt.Println("failed to rename temp file into place:", err)
			removeTmp()
			return
		}

		fmt.Println("Book added successfully")
	}

	// update book by id
	if cmd.update {
		// ensure ID is provided
		if input.ID == "" {
			fmt.Println("Please provide the ID of the book to update")
			return
		} else {
			// find book by id and update fields if they are provided
			// if a field is not provided, it will not be updated
			// e.g. if only title is provided, only title will be updated
			updated := false
			for i, c := range books {
				if c.ID == input.ID {
					if input.Title != "" {
						books[i].Title = input.Title
						updated = true
					}
					if input.Author != "" {
						books[i].Author = input.Author
						updated = true
					}
					if input.Price != "" {
						books[i].Price = input.Price
						updated = true
					}
					if input.IMG != "" {
						books[i].IMG = input.IMG
						updated = true
					}
				}
			}
			// checks if updates were made
			if !updated {
				fmt.Println("Book with ID", input.ID, "not found")
				return
			} else {
				// marshal the entire slice back to JSON (use indent for readability)
				out, err := json.MarshalIndent(books, "", "  ")
				if err != nil {
					fmt.Println("failed to marshal books:", err)
					return
				}

				// write the updated JSON back to the file atomically
				tmpFile, err := os.CreateTemp(".", "books.json.tmp.*")
				if err != nil {
					fmt.Println("failed to create temp file:", err)
					return
				}
				tmpName := tmpFile.Name()

				// ensure temp file is removed on error
				removeTmp := func() {
					tmpFile.Close()
					os.Remove(tmpName)
				}

				if _, err := tmpFile.Write(out); err != nil {
					fmt.Println("failed to write temp file:", err)
					removeTmp()
					return
				}

				if err := tmpFile.Sync(); err != nil {
					fmt.Println("failed to sync temp file:", err)
					removeTmp()
					return
				}

				if err := tmpFile.Close(); err != nil {
					fmt.Println("failed to close temp file:", err)
					removeTmp()
					return
				}

				// set permissions to 0644
				if err := os.Chmod(tmpName, 0644); err != nil {
					fmt.Println("failed to chmod temp file:", err)
					removeTmp()
					return
				}

				// rename temp file into place (atomic on most OSes)
				if err := os.Rename(tmpName, "books.json"); err != nil {
					fmt.Println("failed to rename temp file into place:", err)
					removeTmp()
					return
				}

				fmt.Println("Book updated successfully")
			}
		}
	}

	// delete book by id
	if cmd.delete {
		if input.ID == "" {
			fmt.Println("Please provide the ID of the book to delete")
			return
		} else {
			deleted := false
			for i, c := range books {
				if c.ID == input.ID {
					// remove book from slice
					books = append(books[:i], books[i+1:]...)
					deleted = true
					break
				}
			}
			if !deleted {
				fmt.Println("Book with ID", input.ID, "not found")
				return
			}

			// marshal the entire slice back to JSON (use indent for readability)
			out, err := json.MarshalIndent(books, "", "  ")
			if err != nil {
				fmt.Println("failed to marshal books:", err)
				return
			}

			// write the updated JSON back to the file atomically
			tmpFile, err := os.CreateTemp(".", "books.json.tmp.*")
			if err != nil {
				fmt.Println("failed to create temp file:", err)
				return
			}
			tmpName := tmpFile.Name()

			// ensure temp file is removed on error
			removeTmp := func() {
				tmpFile.Close()
				os.Remove(tmpName)
			}

			if _, err := tmpFile.Write(out); err != nil {
				fmt.Println("failed to write temp file:", err)
				removeTmp()
				return
			}

			if err := tmpFile.Sync(); err != nil {
				fmt.Println("failed to sync temp file:", err)
				removeTmp()
				return
			}

			if err := tmpFile.Close(); err != nil {
				fmt.Println("failed to close temp file:", err)
				removeTmp()
				return
			}

			// set permissions to 0644
			if err := os.Chmod(tmpName, 0644); err != nil {
				fmt.Println("failed to chmod temp file:", err)
				removeTmp()
				return
			}

			// rename temp file into place (atomic on most OSes)
			if err := os.Rename(tmpName, "books.json"); err != nil {
				fmt.Println("failed to rename temp file into place:", err)
				removeTmp()
				return
			}

			fmt.Println("Book deleted successfully")
		}
	}
}
