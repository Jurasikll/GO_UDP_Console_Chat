// Console_Chat_Client project main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var clients_ip_port []string
var clients_name []string

const SERVER_IP_PORT string = ""
const LOCAL_IP string = ""

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {

	var ip_port_server_local string
	var ip_port_local string
	if len(SERVER_IP_PORT) == 0 {
		ip_port_server_local = input_server_ip_port()
	} else {
		ip_port_server_local = SERVER_IP_PORT
	}
	if len(LOCAL_IP) == 0 {
		ip_port_local = input_local_ip() + ":0"
	} else {
		ip_port_local = LOCAL_IP
	}

	ServerAddr, err := net.ResolveUDPAddr("udp", ip_port_server_local)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", ip_port_local)
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	for {
		go check_answer(Conn)
		go check_msg(Conn)
		time.Sleep(time.Second * 1)
	}
}

func check_answer(conn *net.UDPConn) {
	time.Sleep(time.Second * 1)
	answer := make([]byte, 1024)
	n, _, _ := conn.ReadFromUDP(answer)
	fmt.Println(string(answer[0:n]))

}

func input_local_ip() string {
	const msg string = "Введите локальный ip"
	fmt.Println(msg)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func input_server_ip_port() string {
	const MSG string = "введите ip и port сервера в формате {ip}:{port}"
	fmt.Println(MSG)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func check_msg(conn *net.UDPConn) {
	time.Sleep(time.Second * 1)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	//fmt.Println(text)

	buf := []byte(text)
	_, err := conn.Write(buf)
	if err != nil {
		fmt.Println(err)
	}
}
