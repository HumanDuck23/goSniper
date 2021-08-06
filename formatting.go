package main

import "fmt"

var useColors = false

func color(text string, code int, bold bool) string {
	if useColors {
		if bold {
			return fmt.Sprintf("[1;%dm%s[0m", code, text)
		} else {
			return fmt.Sprintf("[%dm%s[0m", code, text)
		}
	} else {
		return text
	}
}

func black(text string, bold bool) string {
	return color(text, 30, bold)
}

func red(text string, bold bool) string {
	return color(text, 31, bold)
}

func green(text string, bold bool) string {
	return color(text, 32, bold)
}

func yellow(text string, bold bool) string {
	return color(text, 33, bold)
}

func blue(text string, bold bool) string {
	return color(text, 34, bold)
}

func magenta(text string, bold bool) string {
	return color(text, 35, bold)
}

func cyan(text string, bold bool) string {
	return color(text, 36, bold)
}

func white(text string, bold bool) string {
	return color(text, 15, bold)
}
