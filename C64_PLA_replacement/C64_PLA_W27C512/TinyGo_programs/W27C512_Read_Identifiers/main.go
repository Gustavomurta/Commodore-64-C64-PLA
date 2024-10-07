// W27C512 Read Identifier - Arduino MEGA
// Gustavo Murta   2024_07_18
// tinygo version 0.32.0 windows/amd64 (using go version go1.22.4 and LLVM version 18.1.2)
// C:\Users\jgust\tinygo\programas\programador_eprom\EEPROM WC27C512 Programmer\W27C512_Read_Identifiers
// tinygo flash -target=arduino-mega2560 main.go
// Connect pin A9 to 12V / Manufacturer ID: 0xDA / Devide ID: 0x08

package main

import (
	"device/avr"
	"machine"
	"time"
)

var address uint32
var delayPeriod = 1 * time.Microsecond // 1 us delay

func main() {

	machine.UART0.Configure(machine.UARTConfig{BaudRate: 115200}) // setup UART0 = 115200 Bps

	// Set PORT registers before DDR registers

	avr.PORTA.Set(0x00) // set Memory address A0 to A7 = 0
	avr.PORTC.Set(0x00) // set Memory address A8 to A15 = 0
	avr.PORTG.Set(0x0F) // set -CE(PG1) and -OE(PG0) = 1
	avr.PORTL.Set(0xFF) // set PORTL - DQ_0 to DQ_7 - Before change to input Pullup

	avr.DDRA.Set(0xFF) // configure PORTA pins as output - Memory Address pins A0 to A7
	avr.DDRC.Set(0xFF) // configure PORTC pins as output - Memory Address pins A8 to A15
	avr.DDRG.Set(0x0F) // configure PORT G0 to G3 pins as output - CE and OE pins
	avr.DDRL.Set(0xFF) // configure PORTL pins as output - DQ_0 to DQ_7 - Before change to input Pullup

	time.Sleep(delayPeriod) // delay

	avr.PORTL.Set(0xFF)     // configure PORTL pins as input PULLUP -  DQ_0 to DQ_7
	avr.DDRL.Set(0x00)      // configure PORTL pins as input -  DQ_0 to DQ_7
	time.Sleep(delayPeriod) // delay

	// EEPROM = 64 X 8 bits  64K=65536

	for address = 0; address < 2; address++ {

		avr.PORTC.Set(uint8(address >> 8)) // shift right 8 bits - Memory address A8 to A15
		avr.PORTA.Set(uint8(address))      // set Memory Address A0 to A7
		time.Sleep(delayPeriod)            // delay
		avr.PORTG.Set(0x00)                // set -CE(PG1) and -OE(PG0) = 0
		time.Sleep(delayPeriod)            // delay

		eepromData := (avr.PINL.Get())      // read Data Memory - I/O_0 to I/O_7
		avr.PORTG.Set(0x01)                 // set -CE(PG1) =0  and -OE(PG0) = 1
		machine.UART0.WriteByte(eepromData) // send byte to serial port 0
	}

	time.Sleep(delayPeriod) // delay
	avr.PORTG.Set(0x0F)     // set -CE(PG1) and -OE(PG0) = 1
	avr.PORTA.Set(0x00)     // set Memory address A0 to A7 = 0
	avr.PORTC.Set(0x00)     // set Memory Address A8 to A15 = 0
}
