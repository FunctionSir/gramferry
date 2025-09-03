/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 14:14:54
 * @LastEditTime: 2025-09-03 15:43:54
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/cmdserver.go
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
	"sync/atomic"

	"github.com/spf13/cobra"
)

func serverShutdown(tcpConn *net.Conn, udpConn *net.UDPConn, stopFlag *atomic.Bool) {
	LogOnErr((*tcpConn).Close())
	LogOnErr(udpConn.Close())
	stopFlag.Store(true)
}

func cmdServer(cmd *cobra.Command, args []string) {
	PrintBanner()
	fmt.Println("Listening TCP:", TCP)
	fmt.Println("Dialing UDP:", UDP)

	log.Println("starting server side service...")

	tcpListen, err := net.Listen("tcp", TCP)
	if LogOnErr(err) {
		panic(err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", UDP)
	if LogOnErr(err) {
		panic(err)
	}

	for {

		TCPConn, err := tcpListen.Accept()
		if LogOnErr(err) {
			continue
		}

		UDPConn, err := net.DialUDP("udp", nil, udpAddr)
		if LogOnErr(err) {
			continue
		}

		var stop atomic.Bool
		stop.Store(false)

		var once sync.Once

		go func(TCPConn net.Conn, UDPConn *net.UDPConn, stop *atomic.Bool) {
			for !stop.Load() {
				szBuf := make([]byte, 2)
				if _, err := io.ReadFull(TCPConn, szBuf); LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}

				var szUint uint16
				if err := binary.Read(bytes.NewReader(szBuf), binary.BigEndian, &szUint); LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}

				buf := make([]byte, szUint)
				if _, err := io.ReadFull(TCPConn, buf); LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}

				if _, err := UDPConn.Write(buf); LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}
			}
		}(TCPConn, UDPConn, &stop)

		go func(TCPConn net.Conn, UDPConn *net.UDPConn, stop *atomic.Bool) {
			for !stop.Load() {
				buf := make([]byte, 65535)
				n, err := UDPConn.Read(buf)
				if LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}

				toTCP := new(bytes.Buffer)

				err = binary.Write(toTCP, binary.BigEndian, uint16(n))
				if LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}

				toTCP.Write(buf[:n])

				if _, err = TCPConn.Write(toTCP.Bytes()); LogOnErr(err) {
					once.Do(func() { serverShutdown(&TCPConn, UDPConn, stop) })
					return
				}
			}
		}(TCPConn, UDPConn, &stop)

	}
}
