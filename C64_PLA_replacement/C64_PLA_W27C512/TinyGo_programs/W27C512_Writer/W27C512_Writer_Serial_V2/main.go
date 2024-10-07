/* W27C512 Writer Serial 64K - Arduino MEGA Version V2
Date: 2024_10_04  Gustavo Murta
go version go1.22.4 windows/amd64
https://zetcode.com/golang/readfile/
https://zetcode.com/golang/writefile/
https://pkg.go.dev/encoding/binary#example-Write
C:\Users\jgust\tinygo\programas\programador_eprom\EEPROM WC27C512 Programmer\W27C512_Writer_Serial
go build main.go
*/

package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

func main() {

	log.Printf(" Start program")

	c := &serial.Config{
		Name:        "COM5", // Arduino MEGA 115200
		Baud:        74880,
		ReadTimeout: time.Millisecond * 1000,
		Size:        8,
		Parity:      0,
		StopBits:    1,
	} // Win10 Serial port setup 74880 Bps- timeout 1 second

	// 74880 57600 38400 19200

	s, err := serial.OpenPort(c) // opens serial port
	if err != nil {              // if any error
		log.Fatal(err) // cancel and print error message
	}

	s.Flush() // garbage clean at serial port

	fw, err := os.Open("C64-pla_original.bin") // open binary file - PLA 00 to FF.bin / C64-pla_original.bin

	if err != nil { // if any error
		log.Fatal(err) // cancel and print error message
	}

	defer fw.Close() // close file at end

	// EEPROM = 64 X 8 bits  64K=65536

	reader := bufio.NewReader(fw)
	buf := make([]byte, 65536) // define PLA buffer 64K bytes

	time.Sleep(915 * time.Millisecond) // very important timeout! 915

	log.Printf(" Start file reading")

	for {
		_, err := reader.Read(buf) // read each byte of the Binary file

		n, erro := s.Write(buf) // send byte to Arduino Serial Port
		if erro != nil {
			log.Fatal(err, n) // if any error cancel and print error message
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}

	log.Printf(" End of file ")
	fmt.Printf("%s", hex.Dump(buf)) // print hexadecimal dump
	log.Printf(" HEX Dump printed ")
}
