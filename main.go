package main

import (
	"fmt"
	"github.com/Tinkerforge/go-api-bindings/dual_button_bricklet"
	"github.com/Tinkerforge/go-api-bindings/dual_button_v2_bricklet"
	"github.com/Tinkerforge/go-api-bindings/io16_bricklet"
	"github.com/Tinkerforge/go-api-bindings/ipconnection"
	"github.com/Tinkerforge/go-api-bindings/stepper_brick"
	"math"
	"strconv"
	"strings"
)

const ADDR string = "localhost:4223"
const DBUID = "mvi"
const IO16UID = "gqN"
const HOURUID = "5Wqtru"

const STEP = -187

const ZERO uint8 = 0b11000011
const ONE uint8 = 0b11100001
const TWO uint8 = 0b11100000
const FOUR uint8 = 0b11111110
const SIX uint8 = 0b11111101
const BLUE uint8 = 0b11010101

var hourRunning = false

func main() {

	fmt.Println(hour["b"])

	ipcon := ipconnection.New()
	defer ipcon.Close()
	db, _ := dual_button_bricklet.New(DBUID, &ipcon)
	io, _ := io16_bricklet.New(IO16UID, &ipcon)
	hourstepper, _ := stepper_brick.New(HOURUID, &ipcon)
	_ = hourstepper.SetMaxVelocity(2000)
	_ = hourstepper.SetSpeedRamping(50000, 50000)

	ipcon.Connect(ADDR)
	defer ipcon.Disconnect()

	db.RegisterStateChangedCallback(func(buttonL dual_button_v2_bricklet.ButtonState, buttonR dual_button_v2_bricklet.ButtonState, ledL dual_button_v2_bricklet.LEDState, ledR dual_button_v2_bricklet.LEDState) {

		if buttonL == dual_button_v2_bricklet.ButtonStatePressed {
			fmt.Println("Left Button: Pressed")
			if !hourRunning {
				toggleHourRunning()
				go rotateHourToPosition(&io, &hourstepper, BLUE)
			}

		} else if buttonL == dual_button_v2_bricklet.ButtonStateReleased {
			fmt.Println("Left Button: Released")
		}

		if buttonR == dual_button_v2_bricklet.ButtonStatePressed {
			fmt.Println("Right Button: Pressed")
			if !hourRunning {
				// toggleHourRunning()
				//go rotateHourToPosition(&io, &hourstepper, SIX)

				// Check Hour on button press
				go rotateToNextPosition('a', &io, &hourstepper)

			}
		} else if buttonR == dual_button_v2_bricklet.ButtonStateReleased {
			fmt.Println("Right Button: Released")
		}

		fmt.Println()
	})

	go io16reader(&io)

	fmt.Print("Press enter to exit.\n")
	_, _ = fmt.Scanln()
}

func toggleHourRunning() {
	hourRunning = !hourRunning
}

func btod(mask uint8) int {

	bits := fmt.Sprintf("%b", mask)
	bit5, _ := strconv.ParseFloat(string(bits[5]), 64)
	bit4, _ := strconv.ParseFloat(string(bits[4]), 64)
	bit3, _ := strconv.ParseFloat(string(bits[3]), 64)
	bit2, _ := strconv.ParseFloat(string(bits[2]), 64)
	bit1, _ := strconv.ParseFloat(string(bits[1]), 64)
	bit0, _ := strconv.ParseFloat(string(bits[0]), 64)

	return int(bit5*math.Pow(2, 5) +
		bit4*math.Pow(2, 4) +
		bit3*math.Pow(2, 3) +
		bit2*math.Pow(2, 2) +
		bit1*math.Pow(2, 1) +
		bit0*math.Pow(2, 0))
}

func rotateToNextPosition(port rune, io *io16_bricklet.IO16Bricklet, stepper *stepper_brick.StepperBrick) {
	valueMask, _ := io.GetPort(port)

	_ = stepper.Enable()

	for valueMask == 255 {
		if 255 == valueMask {
			_ = stepper.DriveBackward()
			valueMask, _ = io.GetPort(port)
		} else {
			_ = stepper.Stop()
			break
		}
	}

	fmt.Printf("Value Mask (Port %s): %b\n", strings.ToUpper(string(port)), valueMask)
	stepper.Disable()
}

func rotateHourToPosition(io *io16_bricklet.IO16Bricklet, stepper *stepper_brick.StepperBrick, position uint8) {

	// Make the mask
	valueMaskA, _ := io.GetPort('a')

	fmt.Println(position == valueMaskA)
	_ = stepper.Enable()

	for true {

		fmt.Printf("%d %d\n", position, valueMaskA)

		if position != valueMaskA {
			//_ = stepper.SetSteps(-50)
			_ = stepper.DriveBackward()
			valueMaskA, _ = io.GetPort('a')
		} else {
			_ = stepper.Stop()
			break
		}
	}

	_ = stepper.Disable()
	toggleHourRunning()
}

func io16reader(io *io16_bricklet.IO16Bricklet) {

	blueMaskHour := 0b11101101

	valueMaskA, _ := io.GetPort('a')
	fmt.Printf("Value Mask (Port A): %b\n", valueMaskA)
	// fmt.Println(valueMaskA == blueMaskHour)

	valueMaskAString := fmt.Sprintf("%b", valueMaskA)
	blueMaskHourString := fmt.Sprintf("%b", blueMaskHour)

	fmt.Println(valueMaskAString == blueMaskHourString)

	// Get current value from port B as bitmask.
	valueMaskB, _ := io.GetPort('b')
	fmt.Printf("Value Mask (Port B): %b\n", valueMaskB)

	// Block that forever
	(chan int)(nil) <- 0

}
