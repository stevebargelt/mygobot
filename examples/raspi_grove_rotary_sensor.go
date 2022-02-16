// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/drivers/aio"
	"github.com/stevebargelt/mygobot/drivers/i2c"
	"github.com/stevebargelt/mygobot/platforms/raspi"
)

func main() {
	board := raspi.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(board)
	sensor := aio.NewGroveRotaryDriver(ads1015, "0")

	work := func() {
		sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{ads1015, sensor},
		work,
	)

	robot.Start()
}
