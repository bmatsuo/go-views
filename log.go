package main

import (
	"log"
	"os"
)

var ErrLogger = log.New(os.Stderr, "", log.LstdFlags)
var OutLogger = log.New(os.Stdout, "", log.LstdFlags)
var debugLevel int

func Debug(level int, v ...interface{}) {
	if debugLevel > level {
		OutLogger.Print(v)
	}
}
func Print(v ...interface{}) {
	OutLogger.Print(v)
}
func Error(v ...interface{}) {
	ErrLogger.Print(v)
}
func Fatal(v ...interface{}) {
	ErrLogger.Fatal(v)
}

func Debugf(level int, format string, v ...interface{}) {
	if debugLevel > level {
		OutLogger.Printf(format, v...)
	}
}
func Printf(format string, v ...interface{}) {
	OutLogger.Printf(format, v...)
}
func Errorf(format string, v ...interface{}) {
	ErrLogger.Printf(format, v...)
}
func Fatalf(format string, v ...interface{}) {
	ErrLogger.Fatalf(format, v...)
}
