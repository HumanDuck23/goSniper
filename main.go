package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	for _, arg := range os.Args {
		if arg == "-c" || arg == "-color" {
			useColors = true
		}
	}
	fmt.Println(
		blue(
			"\n░██████╗░░█████╗░  ░██████╗███╗░░██╗██╗██████╗░███████╗██████╗░\n██╔════╝░██╔══██╗  ██╔════╝████╗░██║██║██╔══██╗██╔════╝██╔══██╗\n██║░░██╗░██║░░██║  ╚█████╗░██╔██╗██║██║██████╔╝█████╗░░██████╔╝\n██║░░╚██╗██║░░██║  ░╚═══██╗██║╚████║██║██╔═══╝░██╔══╝░░██╔══██╗\n╚██████╔╝╚█████╔╝  ██████╔╝██║░╚███║██║██║░░░░░███████╗██║░░██║\n░╚═════╝░░╚════╝░  ╚═════╝░╚═╝░░╚══╝╚═╝╚═╝░░░░░╚══════╝╚═╝░░╚═╝", false),
		)
	fmt.Println("")
	fmt.Println(blue("GOLANG Sniper made by ", true) + green("_Spqghett1#6969", true))
	fmt.Println("")
	// Configuration
	fmt.Println(yellow("What mode of authentication do you want to use?", true))
	fmt.Println(yellow("Mojang login (m) or a token (t)", true))

	mode := input("")
	mode = strings.ToLower(mode)

	var authConfig AuthConfig

	if mode == "m" {
		authConfig.EMAIL = input("Email: ")
		authConfig.PASSWORD = input("Password: ")
		authConfig.TOKEN = getToken(authConfig.EMAIL, authConfig.PASSWORD)
	} else if mode == "t" {
		tokenFile := input("Token file: ")
		dat, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			error("Unable to read token file!")
			error(err.Error())
			os.Exit(-1)
		}
		authConfig.TOKEN = string(dat)

	} else {
		error("Invalid input")
		os.Exit(-1)
	}
	if validateToken(authConfig.TOKEN) {
		success("Token validated!")
	} else {
		error("Token invalid!")
		os.Exit(-1)
	}
}