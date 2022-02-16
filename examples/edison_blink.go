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
	led := gpio.NewLedDriver(e, "13")

	// Uncomment the line below if you are using a Sparkfun
	// Edison board with the GPIO block
	// e.SetBoard("sparkfun")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}