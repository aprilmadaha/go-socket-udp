package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1071")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("------udp客户端已经启动------")
	defer conn.Close()

	go func() { //发送心跳
		for {
			_, err = conn.Write([]byte("keeplive"))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	reader := bufio.NewReader(conn) //读取server发送来的数据
	for {
		dat, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(dat))
	}
}
