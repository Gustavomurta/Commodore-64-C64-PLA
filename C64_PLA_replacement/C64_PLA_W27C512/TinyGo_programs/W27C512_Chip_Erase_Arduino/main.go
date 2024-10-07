// W27C512 CHIP Erase- Arduino MEGA
// Gustavo Murta   2024_10_06
// tinygo version 0.32.0 windows/amd64 (using go version go1.22.4 and LLVM version 18.1.2)
// C:\Users\jgust\tinygo\programas\programador_eprom\EEPROM WC27C512 Programmer\Version V2
// tinygo flash -target=arduino-mega2560 main.go
// Connect pin OE/VPP to 14V
// Connect pin A9 to 14V

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
	avr.PORTL.Set(0xFF) // set PORTL - DQ_0 to DQ_7 = 0xFF

	avr.DDRA.Set(0xFF) // configure PORTA pins as output - Memory Address pins A0 to A7
	avr.DDRC.Set(0xFF) // configure PORTC pins as output - Memory Address pins A8 to A15
	avr.DDRG.Set(0x0F) // configure PORT G0 to G3 pins as output - CE and OE pins
	avr.DDRL.Set(0xFF) // configure PORTL pins as output - DQ_0 to DQ_7

	time.Sleep(delayPeriod) // delay

	// EEPROM = 64 X 8 bits  64K=65536

	// CHIP Erase pulse 100 ms -CE
	address = 0

	avr.PORTC.Set(uint8(address >> 8)) // shift right 8 bits - Memory address A8 to A15
	avr.PORTA.Set(uint8(address))      // set Memory Address A0 to A7
	avr.PORTL.Set(0xFF)                // write Data Memory - D0 to D7 = 0xFF
	time.Sleep(delayPeriod)            // delay 1 us

	avr.PORTG.Set(0x00)                // set -CE(PG1) = 0
	time.Sleep(100 * time.Millisecond) // delay 100 ms ERASE Pulse
	avr.PORTG.Set(0x0F)                // set -CE(PG1) = 1
}
