// +build example
//
// Do not build by default.

package main

import (
	"time"

	"github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/drivers/gpio"
	"github.com/stevebargelt/mygobot/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewAdaptor()
	led := gpio.NewRgbLedDriver(e, "3", "5", "6")

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			led.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("rgbBot",
		[]gobot.Connection{e},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
