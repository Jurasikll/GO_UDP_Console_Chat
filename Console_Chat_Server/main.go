// Console_Chat_Server project main.go
package main

import (
	"fmt"
	"net"
	"os"
)

var clients_ip_port []*net.UDPAddr
var clients_name []string

//var ServerConn *net.UDPConn

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		var temp_client_id int
		temp_client_id = cheack_new_client(addr)
		if temp_client_id == -1 {
			set_new_client(addr, string(buf[0:n]))
			ServerConn.WriteToUDP([]byte("Вы зарегистрированны"), addr)
		} else {
			broadcast(clients_name[temp_client_id], string(buf[0:n]), ServerConn)
		}

		fmt.Println("Received ", string(buf[0:n]), " from ", addr, "  ", n)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func broadcast(client_name string, msg string, srv *net.UDPConn) {
	var msgToAll string
	msgToAll = client_name + ":" + msg
	//msgToAll = msg
	for _, client := range clients_ip_port {
		temp_client, _ := net.ResolveUDPAddr("udp", client.String())
		fmt.Printf("This is addres - %s", client.String())
		srv.WriteToUDP([]byte(msgToAll), temp_client)
	}
}

func set_new_client(addr *net.UDPAddr, name string) int {
	var clints_count int
	clints_count = len(clients_ip_port)

	clients_ip_port = append(clients_ip_port, addr)
	clients_name = append(clients_name, name)

	return clints_count

}

func cheack_new_client(addr *net.UDPAddr) int {
	var res int = -1
	for inx, client := range clients_ip_port {
		if client.String() == addr.String() {
			res = inx
		}
	}
	return res
}
