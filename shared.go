/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 14:20:36
 * @LastEditTime: 2025-09-04 10:06:04
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/shared.go
 */

package main

import (
	"fmt"
	"log"
)

func LogOnErr(err error) bool {
	if err != nil {
		log.Println(err)
	}
	return err != nil
}

func PrintBanner() {
	fmt.Println("GramFerry [ Version 0.1.0 (Wannai Kinuho) ]")
	if Pprof != "" {
		fmt.Println("Serving pprof:", Pprof)
	}
}
