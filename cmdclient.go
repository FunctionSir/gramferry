/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 14:42:24
 * @LastEditTime: 2025-09-03 15:45:06
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
	"sync"

	"github.com/spf13/cobra"
)

func cmdClient(cmd *cobra.Command, args []string) {
	var RemotesMapping sync.Map

	PrintBanner()
	fmt.Println("Listening UDP:", UDP)
	fmt.Println("Dialing TCP:", TCP)

	log.Println("starting client side service...")

	udpAddr, err := net.ResolveUDPAddr("udp", UDP)
	if LogOnErr(err) {
		panic(err)
	}
	udpListen, err := net.ListenUDP("udp", udpAddr)
	if LogOnErr(err) {
		panic(err)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", TCP)
	if LogOnErr(err) {
		panic(err)
	}

	for {
		buf := make([]byte, 65535)
		n, curRemote, err := udpListen.ReadFromUDP(buf)
		if LogOnErr(err) {
			panic(err)
		}
		go func(n int, remote *net.UDPAddr, buf []byte) {
			_, found := RemotesMapping.Load(remote.String())
			if !found {
				TCPConn, err := net.DialTCP("tcp", nil, tcpAddr)
				if LogOnErr(err) {
					panic(err)
				}
				RemotesMapping.Store(remote.String(), TCPConn)
				go func(TCPConn *net.TCPConn, remote *net.UDPAddr) {
					for {
						sz := make([]byte, 2)
						_, err := io.ReadFull(TCPConn, sz)
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

						_, err = udpListen.WriteToUDP(buf, remote)
						if LogOnErr(err) {
							panic(err)
						}
					}
				}(TCPConn, remote)
			}

			tmp, _ := RemotesMapping.Load(remote.String())

			TCPConn := tmp.(*net.TCPConn)

			toTCP := new(bytes.Buffer)
			if err := binary.Write(toTCP, binary.BigEndian, uint16(n)); LogOnErr(err) {
				panic(err)
			}

			toTCP.Write(buf[:n])
			if _, err = TCPConn.Write(toTCP.Bytes()); LogOnErr(err) {
				panic(err)
			}
		}(n, curRemote, buf)
	}
}
