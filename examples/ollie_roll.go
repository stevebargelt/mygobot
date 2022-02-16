// +build example
//
// Do not build by default.

package main

import (
	"os"
	"time"

	"github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/platforms/ble"
	"github.com/stevebargelt/mygobot/platforms/sphero/ollie"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ollie := ollie.NewDriver(bleAdaptor)

	work := func() {
		ollie.SetRGB(255, 0, 255)
		gobot.Every(3*time.Second, func() {
			ollie.Roll(40, uint16(gobot.Rand(360)))
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollie},
		work,
	)

	robot.Start()
}
