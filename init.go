package main

import (
	"flag"
	"log"
	"os"
	"time"
)

const (
	resFormat = `{"error": "%v", "data": "%+v"}`
	shortHand = " (short hand)"

	defaultPort = "8080"
	usagePort   = "PORT at which the server will run"
)

var (
	// PORT at which the server will run (default: 8080),
	// can be modified using flags:
	// 	`-port 80` or `-p 80`
	PORT string
	// Asia/Kolkata - Indian TimeZone's time.Location object
	loc, _ = time.LoadLocation("Asia/Kolkata")
)

func init() {
	flag.StringVar(&PORT, "port", defaultPort, usagePort)
	flag.StringVar(&PORT, "p", defaultPort, usagePort+shortHand)
}

// Is Debug mode - set or not
func Debug() bool {
	return os.Getenv("debug") != ""
}

// Ignore just ignores all the unhandled outputs
func Ignore(args ...interface{}) {
	if Debug() {
		log.Println(args...)
	}
}
