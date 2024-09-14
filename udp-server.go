package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Session struct { //保存会话结构体
	Address    string    //名字
	ClientAddr net.Addr  //socket句柄
	LastActive time.Time //最后刷新时间
}

var sessions = make(map[string]*Session) // 存储活跃会话
var cleanupInterval = 30 * time.Second   // 超时时间间隔

func main() {
	conn, err := net.ListenPacket("udp", ":1071")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("------udp服务器已经启动------")
	defer conn.Close()

	go cleanupSessions() // 启动会话清理 goroutine
	go oneSession(conn)  //给一个会话单独发送
	go allSeesion(conn)  //广播

	buf := make([]byte, 1024)
	for {
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		fmt.Println(addr.String() + " 客户端上线了")
		fmt.Println(string(buf))

		updateSession(addr.String(), addr)
	}
}

func cleanupSessions() { //定时清理超时的会话
	for {
		now := time.Now()
		for addr, session := range sessions {
			if now.Sub(session.LastActive) > cleanupInterval {
				delete(sessions, addr)
				fmt.Println("Session expired:", addr)
			}
		}
		time.Sleep(cleanupInterval)
	}
}

func updateSession(addrString string, addr net.Addr) { //发送心跳的就刷新时间
	session := sessions[addrString]
	if session == nil {
		// 新会话
		sessions[addrString] = &Session{Address: addrString, ClientAddr: addr, LastActive: time.Now()}
		fmt.Println("New Session:", string(sessions[addrString].Address))
	} else {
		// 更新最后活动时间
		session.LastActive = time.Now()
		fmt.Println("update Session:", string(session.Address), session.LastActive)
	}
}

func allSeesion(pc net.PacketConn) { //广播所有会话
	for {
		if len(sessions) > 0 {
			for _, num := range sessions {
				_, err := pc.WriteTo([]byte("大家还好吗？\n"), num.ClientAddr)	//坑：这里要加上\n不然客户端收不到
				if err != nil {
					log.Println(err)
				}
				log.Println("广播发送成功", num.Address, num.ClientAddr)
			}
			time.Sleep(time.Second * 5)
		}

	}
}

func oneSession(pc net.PacketConn) {
	for {
		var s string
		fmt.Print("请输入客户端标识：")
		fmt.Scanf("%s", &s)
		if len(s) > 0 {
			if _, ok := sessions[s]; ok {
				_, err := pc.WriteTo([]byte("单独找你聊天\n"), sessions[s].ClientAddr)
				if err != nil {
					log.Println(err)
				}
				log.Println("发送成功")
			} else {
				log.Println(s + "服务器已经宕机")
			}
		}
	}
}
