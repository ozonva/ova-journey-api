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

	cu := config.NewConfigurationUpdater(time.Second, ConfigFile)
	cu.WatchConfigurationFile()
	fmt.Printf("Used configuration: %v \n", *cu.GetConfiguration())

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
