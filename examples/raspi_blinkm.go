// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/drivers/i2c"
	"github.com/stevebargelt/mygobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	blinkm := i2c.NewBlinkMDriver(r)

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := byte(gobot.Rand(255))
			g := byte(gobot.Rand(255))
			b := byte(gobot.Rand(255))
			blinkm.Rgb(r, g, b)
			color, _ := blinkm.Color()
			fmt.Println("color", color)
		})
	}

	robot := gobot.NewRobot("blinkmBot",
		[]gobot.Connection{r},
		[]gobot.Device{blinkm},
		work,
	)

	robot.Start()
}
