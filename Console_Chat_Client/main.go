// Console_Chat_Client project main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var clients_ip_port []string
var clients_name []string
var not_reading_console_write bool = true
var not_reading_server_answer bool = true
var msg_log []string
var ip_port_server_local string = ""
var ip_port_local string = ""
var nick_name string = ""

const SERVER_IP_PORT string = ""
const LOCAL_IP string = ""
const STR_UDP = "udp"
const EMPTY_STRING = ""

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {

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

	ServerAddr, err := net.ResolveUDPAddr(STR_UDP, ip_port_server_local)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr(STR_UDP, ip_port_local)
	CheckError(err)

	Conn, err := net.DialUDP(STR_UDP, LocalAddr, ServerAddr)
	CheckError(err)
	fmt.Println("Введи ник нажми Enter")
	defer Conn.Close()
	for {
		if not_reading_console_write {
			go check_msg(Conn)
		}
		if not_reading_server_answer {
			go check_answer(Conn)
		}
		time.Sleep(time.Second * 1)
	}
}

//func auth(srv_ip_port string, clt_ip_port string, nick_name string) *net.UDPConn {

//}

func check_answer(conn *net.UDPConn) {
	not_reading_server_answer = false
	answer := make([]byte, 1024)
	n, _, _ := conn.ReadFromUDP(answer)
	print_to_chat(string(answer[0:n]))
	not_reading_server_answer = true
}

func print_to_chat(msg string) {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	msg_log = append(msg_log, msg)
	for _, msg_from_log := range msg_log {
		fmt.Println(msg_from_log)
	}
}

func input_local_ip() string {
	const MSG string = "Введите локальный ip"
	fmt.Println(MSG)
	return read_console_write()
}

func input_server_ip_port() string {
	const MSG string = "введите ip и port сервера в формате {ip}:{port}"
	fmt.Println(MSG)
	return read_console_write()
}

func check_msg(conn *net.UDPConn) {
	not_reading_console_write = false

	text := read_console_write()
	if text != "" {
		buf := []byte(text)
		_, err := conn.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
	}

	not_reading_console_write = true
}

func read_console_write() string {
	const COMMAND_START_PREF = "-"
	var text string
	reader := bufio.NewReader(os.Stdin)
	text, _ = reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, COMMAND_START_PREF) {
		run_command(text)
		text = ""
	}
	return text
}

func run_command(cmd_str string) {
	const COMMAND_ARGS_SEPARATE string = " "
	const COMMAND_EXEC_OK_CODE int = 1
	const COMMAND_NAME int = 0

	var commands_res_code int = 0
	var cmd_arr []string
	var cmd_args []string
	var cmd_arr_len int
	commands_map := map[string]func([]string) int{
		"-help": cmd_help,
		"-h":    cmd_help,
	}

	cmd_arr = strings.Split(cmd_str, COMMAND_ARGS_SEPARATE)
	cmd_arr_len = len(cmd_arr)
	if cmd_arr_len > 1 {
		cmd_args = cmd_arr[1:len(cmd_arr)]
	}
	commands_res_code = commands_map[cmd_arr[COMMAND_NAME]](cmd_args)
	if commands_res_code > COMMAND_EXEC_OK_CODE {
		fmt.Println(commands_res_code)
	}

}
func cmd_help(cmd_args []string) int {

	var temp_arg string
	help_arg_map := map[string]string{
		"help": "Помощь, без аргументов (-help,-h) выводить список команд, с аргументами (-help {название команды},-h {название команды}) выводит описание команды",
		"h":    "Помощь, без аргументов (-help,-h) выводить список команд, с аргументами (-help {название команды},-h {название команды}) выводит описание команды",
	}
	for _, elm := range cmd_args {
		temp_arg = help_arg_map[elm]
		if temp_arg != EMPTY_STRING {
			fmt.Println(temp_arg)
		}
	}
	return 1
}
