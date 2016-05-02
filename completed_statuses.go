package main

// Types of completed status codes.
const (
	queueWait = iota
	completed
	failedCmd
	failedInvalidReturnFormat
)
