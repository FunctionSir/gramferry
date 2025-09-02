/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 14:20:36
 * @LastEditTime: 2025-09-02 17:43:57
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/shared.go
 */

package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func LogOnErr(err error) bool {
	if err != nil {
		log.Println(err)
	}
	return err != nil
}

func ParseUDPAddr(s string) (net.UDPAddr, error) {
	splited := strings.Split(s, ":")
	port, err := strconv.Atoi(splited[1])
	if err != nil {
		return net.UDPAddr{}, err
	}
	return net.UDPAddr{IP: net.ParseIP(splited[0]), Port: port, Zone: ""}, nil
}

func ParseTCPAddr(s string) (net.TCPAddr, error) {
	splited := strings.Split(s, ":")
	port, err := strconv.Atoi(splited[1])
	if err != nil {
		return net.TCPAddr{}, err
	}
	return net.TCPAddr{IP: net.ParseIP(splited[0]), Port: port, Zone: ""}, nil
}

func PrintBanner() {
	fmt.Println("GramFerry [ Version 0.1.0 (Wannai Kinuho) ]")
}
