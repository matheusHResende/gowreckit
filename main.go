package main

import (
	"fmt"
	"log"

	"github.com/zan8in/masscan"
)

func main() {
	fmt.Println("Hello world!")
	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("146.56.202.100/24"), // Injetar ip
		masscan.SetParamPorts("80"),                  // Injetar porta
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
	)
	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}
	scanResult, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("masscan encountered an error: %v", err)
	}

	if scanResult != nil {
		for i, v := range scanResult.Hosts {
			fmt.Printf("Host: %s Port: %v\n", v.IP, scanResult.Ports[i].Port)
		}
		fmt.Println("hosts len : ", len(scanResult.Hosts))
	}
}
