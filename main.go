package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v2"
	"github.com/jlaffaye/ftp"
	"github.com/zan8in/masscan"
)

func runMasscan(ip string, ports ...string) ([]int, []int) {
	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets(ip),     // Injetar ip
		masscan.SetParamPorts(ports...), // Injetar porta
		masscan.SetParamRate(10000),
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
	tcpPorts := []int{}
	udpPorts := []int{}
	for _, port := range scanResult.Ports {
		if port.Proto == "tcp" {
			tcpPorts = append(tcpPorts, port.Port)
		} else {
			udpPorts = append(udpPorts, port.Port)
		}

	}
	return tcpPorts, udpPorts
}

func portsStringfy(ports []int) string {
	output := []string{}
	for _, port := range ports {
		output = append(output, fmt.Sprint(port))
	}
	return strings.Join(output, ",")
}

type Service struct {
	host     string
	port     uint16
	protocol string
	name     string
	product  string
	version  string
}

func (s Service) testFTP() {
	fmt.Printf("%s:%d\n", s.host, s.port)
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", s.host, s.port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		log.Fatal(err)
	}
	output, err := c.List(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range output {
		fmt.Printf("%#v", entry)
	}
	//Testar se permite anonimo
	//Testar subir arquivo
	//Enumerar arquivos (path)
}

type Services []Service

func (s Services) process() {
	for _, service := range s {
		fmt.Printf("%+v\n", service)

		switch service.name {
		case "ftp":
			service.testFTP()
		case "ssh":
			fallthrough
		case "http":
			fallthrough
		case "msrpc":
			fallthrough
		case "domain":
			fallthrough
		case "ms-wbt-server":
			fallthrough
		default:
			fmt.Println("Não implementei ainda o negócio do " + service.name)
		}
	}
}

var services Services

func runNmap(ip string, ports []int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// Equivalent to `/usr/local/bin/nmap -p 80,443,843 google.com facebook.com youtube.com`,
	// with a 5 minute timeout.
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(ip),
		nmap.WithPorts(portsStringfy(ports)),
		nmap.WithContext(ctx),
		nmap.WithSkipHostDiscovery(),
		nmap.WithServiceInfo(),
		// nmap.WithAggressiveScan(),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	fmt.Println(result.Args)
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	if warnings != nil {
		log.Printf("Warnings: \n %v", warnings)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		fmt.Printf("Host %q:\n", host.Addresses[0])

		for _, port := range host.Ports {
			services = append(services, Service{
				port:     port.ID,
				host:     host.Addresses[0].Addr,
				protocol: port.Protocol,
				name:     port.Service.Name,
				product:  port.Service.Product,
				version:  port.Service.Version,
			})
		}
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}

func main() {
	//ip := "10.10.59.3"
	// tcpPorts, _ := runMasscan(ip, "1-65535")
	//tcpPorts := []int{21}
	//log.Println(tcpPorts)
	//runNmap(ip, tcpPorts)
	//services.process()

	host, port := "10.10.59.3", 21
	fmt.Printf("%s:%d\n", host, port)
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", host, port), ftp.DialWithTimeout(5*time.Second), ftp.DialWithDebugOutput(os.Stdout))
	if err != nil {
		fmt.Println("Primeiro erro")
		log.Fatal(err)
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		fmt.Println("Segundo erro")
		log.Fatal(err)
	}
	output, err := c.List("")
	if err != nil {
		fmt.Println("Terceiro erro")
		log.Fatal(err)
	}
	for _, entry := range output {
		fmt.Printf("%#v\n", entry)
	}
	//Testar se permite anonimo
	//Testar subir arquivo
	//Enumerar arquivos (path)
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
