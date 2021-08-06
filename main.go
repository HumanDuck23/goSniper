package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

	//info(strconv.Itoa(int(getPing())))
	// Configuration
	fmt.Println(yellow("What mode of authentication do you want to use?", true))
	fmt.Println(yellow("Mojang login (m), Microsoft login (ms) or a token (t)", true))

	mode := input("")
	mode = strings.ToLower(mode)

	var authConfig AuthConfig

	if mode == "m" {
		authConfig.EMAIL = input("Email: ")
		authConfig.PASSWORD = input("Password: ")
		authConfig.TOKEN = getToken(authConfig.EMAIL, authConfig.PASSWORD)
	} else if mode == "ms" {
		authConfig.EMAIL = input("Email: ")
		authConfig.PASSWORD = input("Password: ")
		authConfig.TOKEN = getTokenMS(authConfig.EMAIL, authConfig.PASSWORD)
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
		var config SnipeConfig
		config.USERNAME = input("Name you want to snipe: ")
		config.TOKEN = authConfig.TOKEN
		changeSkin(config.TOKEN)
		config.DROPTIME = getDropTime(config.USERNAME)
		if config.DROPTIME <= 0 {
			error("Unable to get droptime. Try start the sniper sooner before the name becomes available.")
			os.Exit(-1)
		}
		config.OFFSET, _ = strconv.Atoi(input("Offset you want to use (PING will be added): "))
		var wg sync.WaitGroup
		wg.Add(1)
		go snipe(config, &wg)
		wg.Wait()
		fmt.Println(time.Now())
	} else {
		error("Token invalid!")
		os.Exit(-1)
	}
}
