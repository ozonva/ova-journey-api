package main

import (
	"bufio"
	"fmt"
	"github.com/ozonva/ova-journey-api/internal/config"
	"os"
	"time"
)

const ConfigFile = "config/config.json"

func main() {
	fmt.Println("Hello, I'm ova-journey-api")

	appConfig := config.Configuration{}
	appConfig.LoadConfiguration(ConfigFile)

	go func() {
		for {
			appConfig.LoadConfiguration(ConfigFile)
			time.Sleep(time.Second)
		}
	}()

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
