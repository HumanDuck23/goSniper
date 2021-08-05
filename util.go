package main

import (
	"fmt"
	"os"
	"time"
)

func input(query string) string {
	fmt.Printf(white("[", true) + yellow("input", true) + white("]", true) + " " + white(query, true))
	var input string
	fmt.Scanln(&input)
	return input
}

func error(message string) {
	fmt.Println(white("[", true) + red("error", true) + white("]", true) + " " + white(message, true))
}

func timeString() string {
	t := time.Now()
	str := fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return str
}

func timeError(message string) {
	fmt.Println(white("[", true) + magenta(timeString(), true) + white("]", true) + " " + white("[", true) + red("error", true) + white("]", true) + " " + white(message, true))
}

func success(message string) {
	fmt.Println(white("[", true) + green("success", true) + white("]", true) + " " + white(message, true))
}

func info(message string) {
	fmt.Println(white("[", true) + cyan("info", true) + white("]", true) + " " + white(message, true))
}

func timeInfo(message string) {
	fmt.Println(white("[", true) + magenta(timeString(), true) + white("]", true) + " " + white("[", true) + cyan("info", true) + white("]", true) + " " + white(message, true))
}

func log(str string) {
	f, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(str + "\n"); err != nil {
		fmt.Println(err)
	}
}
