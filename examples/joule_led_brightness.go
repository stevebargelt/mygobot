// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/drivers/gpio"
	"github.com/stevebargelt/mygobot/platforms/intel-iot/joule"
)

func main() {
	e := joule.NewAdaptor()
	led := gpio.NewLedDriver(e, "J12_26")

	work := func() {
		brightness := uint8(0)
		fadeAmount := uint8(15)

		gobot.Every(100*time.Millisecond, func() {
			err := led.Brightness(brightness)
			if err != nil {
				fmt.Println(err)
			}
			brightness = brightness + fadeAmount
			if brightness == 0 || brightness == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("pwmBot",
		[]gobot.Connection{e},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
