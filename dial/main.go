package main

import (
	"fmt"
	"net"
	"sort"
)

func main() {
	ports := make(chan int, 100)
	results := make(chan int, 100)
	portsLen := 1024
	openPorts := make([]int, 0)
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	go func() {
		for i := 1; i <= portsLen; i++ {
			ports <- i
		}
	}()
	for i := 0; i < portsLen; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}
}

func worker(ports <-chan int, results chan<- int) {
	for {
		p, ok := <-ports
		if !ok {
			break
		}
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		fmt.Printf("try dial %s\n", address)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		defer conn.Close()
		results <- p
	}
}
