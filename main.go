package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/eiannone/keyboard"
	"github.com/joeyak/hoin-printer"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var addr string

	flag.StringVar(&addr, "adder", "192.168.1.23:9100", "Address to connect printer too")
	flag.Parse()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("unable to dial:", err)
		return
	}
	defer conn.Close()

	printer := hoin.NewPrinter(conn)

	checkErr(keyboard.Open())
	defer keyboard.Close()

	fmt.Println("Press ctrl+c to escape")
	fmt.Println("Press ctrl+x to cut")
	fmt.Println("Press ctrl+f to feed")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatalln(err)
		}

		s := string(char)

		if key == keyboard.KeyEnter {
			s = "\n"
		} else if key == keyboard.KeyCtrlX {
			fmt.Println("~~CUT~~")
			checkErr(printer.CutFeed(20))
			continue
		} else if key == keyboard.KeyCtrlF {
			fmt.Println()
			checkErr(printer.FeedLines(5))
			continue
		} else if key == keyboard.KeySpace {
			s = " "
		}

		fmt.Print(s)
		printer.Print(s)
	}
}
