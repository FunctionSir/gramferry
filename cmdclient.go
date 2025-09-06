/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 14:42:24
 * @LastEditTime: 2025-09-04 13:38:23
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/cmdclient.go
 */

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func cmdClient(cmd *cobra.Command, args []string) {
	var RemotesMapping sync.Map
	var RemotesLastSeen sync.Map

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

	var bufPool = sync.Pool{
		New: func() any {
			return make([]byte, 65535)
		},
	}

	// Connection Cleaner
	go func() {
		for range time.Tick(60 * time.Second) {
			toRemove := make([]string, 0)
			RemotesLastSeen.Range(func(key, value any) bool {
				if value.(time.Time).Add(300 * time.Second).Before(time.Now()) {
					toRemove = append(toRemove, key.(string))
				}
				return true
			})
			for _, r := range toRemove {
				conn, found := RemotesMapping.Load(r)
				if !found {
					RemotesLastSeen.Delete(r)
					continue
				}
				LogOnErr(conn.(*net.TCPConn).Close())
				RemotesMapping.Delete(r)
				RemotesLastSeen.Delete(r)
				log.Println("the TCP connection for UDP remote", r, "closed since it was not used more than 300s")
			}
		}
	}()

	for {
		bufInPool := bufPool.Get()
		buf := bufInPool.([]byte)
		n, curRemote, err := udpListen.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		go func(remote *net.UDPAddr) {
			RemotesLastSeen.Store(remote.String(), time.Now())
		}(curRemote)

		go func(n int, remote *net.UDPAddr, buf []byte, bufInPool any) {
			defer bufPool.Put(bufInPool)
			_, found := RemotesMapping.Load(remote.String())
			if !found {
				TCPConn, err := net.DialTCP("tcp", nil, tcpAddr)
				if err != nil {
					panic(err)
				}
				log.Println("established TCP connection for UDP remote", remote.String())
				RemotesMapping.Store(remote.String(), TCPConn)
				go func(TCPConn *net.TCPConn, remote *net.UDPAddr) {
					for {
						go func(remote *net.UDPAddr) {
							RemotesLastSeen.Store(remote.String(), time.Now())
						}(curRemote)

						sz := make([]byte, 2)
						_, err := io.ReadFull(TCPConn, sz)
						if err != nil {
							if !errors.Is(err, net.ErrClosed) {
								panic(err)
							}
							return
						}

						szReader := bytes.NewReader(sz)
						var szUint uint16
						err = binary.Read(szReader, binary.BigEndian, &szUint)
						if err != nil {
							if !errors.Is(err, net.ErrClosed) {
								panic(err)
							}
							return
						}

						buf := make([]byte, szUint)
						_, err = io.ReadFull(TCPConn, buf)
						if err != nil {
							if !errors.Is(err, net.ErrClosed) {
								panic(err)
							}
							return
						}

						_, err = udpListen.WriteToUDP(buf, remote)
						if err != nil {
							if !errors.Is(err, net.ErrClosed) {
								panic(err)
							}
							return
						}
					}
				}(TCPConn, remote)
			}

			tmp, _ := RemotesMapping.Load(remote.String())

			TCPConn := tmp.(*net.TCPConn)

			toTCP := new(bytes.Buffer)
			if err := binary.Write(toTCP, binary.BigEndian, uint16(n)); err != nil {
				if !errors.Is(err, net.ErrClosed) {
					panic(err)
				}
				return
			}

			toTCP.Write(buf[:n])
			if _, err = TCPConn.Write(toTCP.Bytes()); err != nil {
				if !errors.Is(err, net.ErrClosed) {
					panic(err)
				}
				return
			}
		}(n, curRemote, buf, bufInPool)
	}
}
