package main

import (
	"log"
)

// ref: https://www.reddit.com/r/golang/comments/6irpt1/is_there_a_golang_logging_library_around_that/

var logger = make(chan []any, 10000)
var isBufferEnabled = false

func InitBuffer() {
	isBufferEnabled = true
	go func() {
		for msg := range logger {
			log.Println(msg...)
		}
	}()
}

func Print(v ...any) {

	if isBufferEnabled {
		logger <- []any{"hello"}
		return
	}

	go log.Println(v...)
}
