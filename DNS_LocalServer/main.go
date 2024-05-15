package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"server/network"
	"server/service"
)

func main() {
	netConfig, err := parseNetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	dbConfig, err := parseDbConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.InitMysql(dbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server Start")
	startServer(netConfig)
}

func startServer(netConfig network.NetConfig) {
	address := fmt.Sprintf(":%d", netConfig.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	fmt.Println("Client Detected- ")

	recvData := make([]byte, 256)

	_, err := conn.Read(recvData)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	if int(recvData[0]) == network.PACKET_GET_DNS {
		handler_GetDNS(conn, recvData)
	} else if int(recvData[0]) == network.PACKET_REG_DNS {
		handler_RegisterDNS(conn, recvData)
	}
}

func handler_GetDNS(conn net.Conn, recvData []byte) {
	domainName := string(recvData[1:])
	ip := service.LoadDNS(domainName)

	if ip == "NoSuchData" {
		ip = getDNSFromUpperDNS(recvData)
	}

	sendData := []byte(ip)

	n, err := conn.Write(sendData)

	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	fmt.Println("Send Data:", domainName, ":", ip, ",", n, "bytes")
}

func handler_RegisterDNS(conn net.Conn, recvData []byte) {
	dnLength := binary.BigEndian.Uint16(recvData[1:3])
	ipLength := binary.BigEndian.Uint16(recvData[3:5])

	dns := string(recvData[5 : 5+dnLength])
	ip := string(recvData[5+dnLength : 5+dnLength+ipLength])
	result := service.RegsterDNS(ip, dns)

	regDNSToUpperDNS(recvData)

	sendData := []byte(result)

	n, err := conn.Write(sendData)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	fmt.Println("Send Data:", result, n, "bytes")
}

func parseNetConfig() (network.NetConfig, error) {
	var netConfig network.NetConfig
	file, err := os.Open("./config/net_config.json")
	if err != nil {
		return netConfig, err
	}

	defer file.Close()

	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(&netConfig)

	return netConfig, err
}

func parseDbConfig() (service.DbConfig, error) {
	var dbConfig service.DbConfig
	file, err := os.Open("./config/db_config.json")

	if err != nil {
		return dbConfig, nil
	}

	defer file.Close()

	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(&dbConfig)

	return dbConfig, nil
}

func getDNSFromUpperDNS(sendData []byte) string {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		conn.Close()
		return err.Error()
	}

	n, err := conn.Write(sendData)

	if err != nil {
		conn.Close()
		return err.Error()
	}

	recvBuff := make([]byte, 256)
	recvBytes, err := conn.Read(recvBuff[0:])

	if err != nil {
		conn.Close()
		return err.Error()
	}

	fmt.Println("Write/Read Byte: ", n, recvBytes)

	return string(recvBuff)
}

func regDNSToUpperDNS(sendData []byte) {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		conn.Close()
		return
	}

	n, err := conn.Write(sendData)

	if err != nil {
		conn.Close()
		return
	}

	recvBuff := make([]byte, 256)
	recvBytes, err := conn.Read(recvBuff[0:])

	if err != nil {
		conn.Close()
		return
	}

	fmt.Println("Write/Read Byte: ", n, recvBytes)
}
