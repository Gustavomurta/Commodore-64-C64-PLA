/* W27C512 Reader Serial 64K - Arduino MEGA - Version V2
Date: 2024_10_04  Gustavo Murta
go version go1.22.4 windows/amd64
https://zetcode.com/golang/readfile/
https://zetcode.com/golang/writefile/
https://pkg.go.dev/encoding/binary#example-Write
C:\Users\jgust\tinygo\programas\programador_eprom\EEPROM WC27C512 Programmer\W27C512_Reader_Serial
go build main.go
*/

package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

func main() {

	c := &serial.Config{
		Name:        "COM5",
		Baud:        115200,
		ReadTimeout: time.Second * 1,
		Size:        8,
		Parity:      0,
		StopBits:    1,
	} // Win10 Serial port setup 115200 Bps- timeout 1 second

	s, err := serial.OpenPort(c) // opens serial port
	if err != nil {              // if any error
		log.Fatal(err) // cancel and print error message
	}

	s.Flush() // garbage clean at serial port

	fw, err := os.Create("W27C512_.bin") // create a bin format file up to 512K Bytes FLASH_MX28F1000.bin

	if err != nil { // if any error
		log.Fatal(err) // cancel and print error message
	}

	defer fw.Close() // close file at end

	buf := make([]byte, 1)    // define read buffer
	bufROM := make([]byte, 0) // define ROM buffer

	log.Printf(" W27C512 read start")

	// EEPROM = 64 X 8 bits  64K=65536

	for i := 0; i < (65536); i++ { // read 64K memory with 8 bits data
		n, err := s.Read(buf) // read serial port data
		if err != nil {       // if any error
			log.Fatal(err, n) // cancel and print error message
		}
		bufROM = append(bufROM, buf...) // filling the ROM buffer
	}

	log.Printf(" W27C512 save file ") // print

	for _, data := range bufROM {

		err := binary.Write(fw, binary.LittleEndian, data) // write ROM buffer into BIN format file

		if err != nil { // if any error
			log.Fatal(err) // cancel and print error message
		}
	}
	log.Printf(" Save W27C512 file OK! ")
	fmt.Printf("%s", hex.Dump(bufROM[:])) // print HEX dump with ASCII data
	log.Printf(" End of program ")

}
