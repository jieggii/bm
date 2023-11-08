package logger

import (
	"log"
	"os"
)

var Stdout = log.New(os.Stdout, "", 0)
var Stderr = log.New(os.Stderr, "", 0)

var Plus = log.New(os.Stdout, "(+) ", 0)
var Minus = log.New(os.Stdout, "(-) ", 0)
var Update = log.New(os.Stdout, "(*) ", 0)
