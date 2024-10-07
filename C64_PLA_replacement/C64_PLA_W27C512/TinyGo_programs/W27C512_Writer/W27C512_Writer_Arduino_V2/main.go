// W27C512 Writer 64K - Arduino MEGA - Version V2
// Gustavo Murta   2024_10_04
// tinygo version 0.32.0 windows/amd64 (using go version go1.22.4 and LLVM version 18.1.2)
// C:\Users\jgust\tinygo\programas\programador_eprom\EEPROM WC27C512 Programmer\W27C512_Writer_Arduino
// tinygo flash -target=arduino-mega2560 main.go
// Program CHIP => Connect -OE/VPP pin to 12V
// Pulse 100 uS to -CE pin

package main

import (
	"device/avr"
	"machine"
	"time"
)

var address uint32

var eepromData byte
var delayPeriod = 1 * time.Microsecond //  us delay

func main() {

	machine.UART0.Configure(machine.UARTConfig{BaudRate: 74880}) // setup UART0 = 74880 Bps

	// 115200 74880 57600 38400 19200
	// here is timeout of 915 miliseconds
	// Set PORT registers before DDR registers

	avr.PORTA.Set(0xFF) // set EEPROM A0 to A7 = 0xFF
	avr.PORTC.Set(0xFF) // set EEPROM A8 to A15 = 0xFF
	avr.PORTL.Set(0xFF) // set EEPROM D0 to D7 = 0xFF
	avr.PORTG.Set(0x0F) // set OE#(PG0) = 1 and CE#(PG1) = 1

	avr.DDRA.Set(0xFF) // configure PORTA pins as output - EEPROM A0 to A7
	avr.DDRC.Set(0xFF) // configure PORTC pins as output - EEPROM A8 to A15
	avr.DDRG.Set(0x0F) // configure PORT G0 to G3 pins as output - pin PG0 = OE# / pin PG1 = CE#
	avr.DDRL.Set(0xFF) // configure PORTL pins as output - EEPROM D0 to D7

	// EEPROM = 64 X 8 bits  64K=65536

	for address = 0; address < 65536; {

		if machine.UART0.Buffered() > 0 {

			eepromData, err := machine.UART0.ReadByte() // Read byte from Binary File
			avr.PORTG.Set(0x0E)                         // set OE#(PG0) = 0 and CE#(PG1) = 1  0xE: 1110 activate 12V

			if err == nil {

				avr.PORTC.Set(uint8(address >> 8)) // shift right 8 bits - Memory address A8 to A15
				avr.PORTA.Set(uint8(address))      // set Memory Address A0 to A7

				avr.PORTL.Set(eepromData)         // write Data Memory - D0 to D7
				time.Sleep(delayPeriod)           // delay
				avr.PORTG.Set(0x0C)               // set OE#(PG0) = 0 and CE#(PG1) = 0  0xB: 1100
				time.Sleep(83 * time.Microsecond) // CE Pulse 100 us delay - VALUE MUST BE 83 to serial syn
				avr.PORTG.Set(0x0E)               // set OE#(PG0) = 0 and CE#(PG1) = 1  0xE: 1110
				address++
			}
		}
	}
	time.Sleep(delayPeriod) // delay
	avr.PORTG.Set(0x0F)     // set OE#(PG0) = 1 and CE#(PG1) = 1
	avr.PORTA.Set(0x00)     // set EEPROM A0 to A7 = 0x00
	avr.PORTC.Set(0x00)     // set EEPROM A8 to A15 = 0x00
	avr.PORTL.Set(0xFF)     // set EEPROM D0 to D7 = 0xFF
}
