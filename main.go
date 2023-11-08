package main

import (
	"fmt"
	"github.com/jieggii/bm/internal/args"
	"github.com/jieggii/bm/internal/env"
	"github.com/jieggii/bm/internal/logger"
	"github.com/jieggii/bm/internal/storage"
	"os"
	"strconv"
)

func fmtBookStatus(bookTitle string, bookID string, page int) string {
	return fmt.Sprintf("\"%v\" (%v): page %v.", bookTitle, bookID, page)
}

func CommandNew(arguments []string) error {
	bookTitle := arguments[0]

	vars := env.Read()
	db := storage.New(vars.StorageHome, vars.StorageFileName)

	// read books from storage:
	books, err := db.Read()
	if err != nil {
		return err
	}

	// add the new book to the storage
	bookID := storage.GenerateBookID(bookTitle)
	books[bookID] = storage.Book{Title: bookTitle, PageNumber: 1}

	// write changes to the storage
	err = db.Write(books)
	if err != nil {
		return err
	}

	// print the book ID
	logger.Plus.Printf("Created a new book \"%v\", id: %v.", bookTitle, bookID)
	return nil
}

func CommandShow(arguments []string) error {
	bookID := arguments[0]

	vars := env.Read()
	db := storage.New(vars.StorageHome, vars.StorageFileName)

	// read books from storage:
	books, err := db.Read()
	if err != nil {
		return err
	}

	// get the book:
	book, found := books[bookID]
	if !found {
		return fmt.Errorf("book %v does not exist", bookID)
	}

	// print information about the book
	logger.Stdout.Printf(fmtBookStatus(book.Title, bookID, book.PageNumber))

	return nil
}

func CommandLs(_ []string) error {
	vars := env.Read()
	db := storage.New(vars.StorageHome, vars.StorageFileName)

	// read books from storage:
	books, err := db.Read()
	if err != nil {
		return err
	}

	// print all of them:
	for bookID, book := range books {
		logger.Stdout.Printf("- %v", fmtBookStatus(book.Title, bookID, book.PageNumber))
	}

	return nil
}

func CommandSet(arguments []string) error {
	bookID := arguments[0]
	pageString := arguments[1]
	pageNumber, err := strconv.Atoi(pageString)
	if err != nil {
		return err
	}

	vars := env.Read()
	db := storage.New(vars.StorageHome, vars.StorageFileName)

	// fetch books from the storage:
	books, err := db.Read()
	if err != nil {
		return err
	}

	// find the book:
	book, found := books[bookID]
	if !found {
		return fmt.Errorf("book %v does not exist", bookID)
	}

	// update book page number:
	book.PageNumber = pageNumber
	books[bookID] = book

	// write changes to the storage:
	err = db.Write(books)
	if err != nil {
		return err
	}

	// print information about the book:
	logger.Update.Printf(fmtBookStatus(book.Title, bookID, book.PageNumber))

	return nil

}

func CommandRm(arguments []string) error {
	bookID := arguments[0]

	vars := env.Read()
	db := storage.New(vars.StorageHome, vars.StorageFileName)

	// read books from storage:
	books, err := db.Read()
	if err != nil {
		return err
	}

	// check if book exists:
	book, found := books[bookID]
	if !found {
		return fmt.Errorf("book %v does not exist", bookID)
	}
	delete(books, bookID)

	// write changes to the storage
	err = db.Write(books)
	if err != nil {
		return err
	}

	logger.Minus.Printf("Deleted book \"%v\" (%v).", book.Title, bookID)
	return nil

}

func main() {
	argsHandler := args.Handler{
		ProgramMeta: args.ProgramMeta{
			Name:        "bm",
			Description: "bm - simple bookmarking tool (literally for books)",
			Version:     "0.1.0",
			Author:      "jieggii <jieggii@protonmail.com>",
		},
		Args: os.Args,
		Commands: []*args.Command{
			{
				Name:        "new",
				Usage:       "new <title>",
				Description: "create a new book",
				ArgsNumber:  1,
				Handler:     CommandNew,
			},
			{
				Name:        "show",
				Usage:       "show <id>",
				Description: "show page number of the book <id>",
				ArgsNumber:  1,
				Handler:     CommandShow,
			},
			{
				Name:        "ls",
				Usage:       "ls",
				Description: "list all books",
				ArgsNumber:  0,
				Handler:     CommandLs,
			},
			{
				Name:        "set",
				Usage:       "set <id> <page>",
				Description: "set page number of book <id> to <page>",
				ArgsNumber:  2,
				Handler:     CommandSet,
			},
			{
				Name:        "rm",
				Usage:       "rm <id>",
				Description: "remove book",
				ArgsNumber:  1,
				Handler:     CommandRm,
			},
		},
	}
	if err := argsHandler.Handle(); err != nil {
		logger.Stderr.Printf("bm: %v", err)
	}
}
