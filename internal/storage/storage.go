// Package storage stands for storing and retrieving information about books.
package storage

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
)

const HomePermission = 0700 // storage home directory permissions
const FilePermission = 0600 // storage file permissions

const bookIDChars = 7  // number of chars in the beginning of a book ID
const bookIDDigits = 3 // number of digits in the end of book ID

type Book struct {
	Title      string `json:"title"`
	PageNumber int    `json:"page_number"`
}

type Books map[string]Book

// GenerateBookID generates unique random readable book identifier.
func GenerateBookID(title string) string {
	title = strings.Join(strings.Fields(title), "") // remove all spaces
	title = strings.ToLower(title)                  // to lower case

	id := ""
	if len(title) >= bookIDChars {
		id = title[0:bookIDChars]
	} else {
		id = title[0:]
	}

	for len(id) < bookIDChars+bookIDDigits {
		id += strconv.Itoa(rand.Intn(9))
	}

	return id
}

// Storage represents storage file where books are stored.
type Storage struct {
	// Directory where storage file is located.
	Home string

	// Name of the storage file.
	FileName string

	// Path to the storage file (it is simply Home + FileName).
	path string
}

// New creates a new instance of Storage.
func New(home string, fileName string) *Storage {
	return &Storage{
		Home:     home,
		FileName: fileName,
		path:     path.Join(home, fileName),
	}
}

// Read reads data from books storage.
func (s *Storage) Read() (Books, error) {
	homeExists, err := pathExists(s.Home) // check if storage home exists
	if err != nil {
		return nil, err
	}

	if !homeExists { // create storage home if it does not exist
		if err := os.MkdirAll(s.Home, HomePermission); err != nil {
			return nil, err
		}
	}

	fileExists, err := pathExists(s.path)
	if err != nil {
		return nil, err
	}

	if !fileExists { // create empty json storage file if it does not exist.
		if err := s.Write(make(Books)); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	data := make(Books)
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// Write writes data to books storage.
func (s *Storage) Write(books Books) error {
	content, err := json.MarshalIndent(books, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, content, FilePermission)
}
