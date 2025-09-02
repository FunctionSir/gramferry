/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 14:42:24
 * @LastEditTime: 2025-09-02 17:44:55
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/cmdclient.go
 */

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
)

var Remote *net.UDPAddr

func cmdClient(cmd *cobra.Command, args []string) {
	PrintBanner()
	fmt.Println("Listening UDP:", UDP)
	fmt.Println("Dialing TCP:", TCP)

	log.Println("starting client side service...")

	udpAddr, err := ParseUDPAddr(UDP)
	if LogOnErr(err) {
		panic(err)
	}
	udpListen, err := net.ListenUDP("udp", &udpAddr)
	if LogOnErr(err) {
		panic(err)
	}
	tcpAddr, err := ParseTCPAddr(TCP)
	if LogOnErr(err) {
		panic(err)
	}
	TCPConn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if LogOnErr(err) {
		panic(err)
	}

	go func() {
		for {
			buf := make([]byte, 65535)
			n := 0
			n, Remote, err = udpListen.ReadFromUDP(buf)
			if LogOnErr(err) {
				panic(err)
			}

			toTCP := new(bytes.Buffer)
			if err := binary.Write(toTCP, binary.BigEndian, uint16(n)); LogOnErr(err) {
				panic(err)
			}

			toTCP.Write(buf[:n])
			if _, err = TCPConn.Write(toTCP.Bytes()); LogOnErr(err) {
				panic(err)
			}
		}
	}()

	go func() {
		for {
			sz := make([]byte, 2)
			_, err = io.ReadFull(TCPConn, sz)
			if LogOnErr(err) {
				panic(err)
			}

			szReader := bytes.NewReader(sz)
			var szUint uint16
			err = binary.Read(szReader, binary.BigEndian, &szUint)
			if LogOnErr(err) {
				panic(err)
			}

			buf := make([]byte, szUint)
			_, err = io.ReadFull(TCPConn, buf)
			if LogOnErr(err) {
				panic(err)
			}

			_, err = udpListen.WriteToUDP(buf, Remote)
			if LogOnErr(err) {
				panic(err)
			}
		}
	}()
	for {
		time.Sleep(24 * 365 * time.Hour)
	}
}
