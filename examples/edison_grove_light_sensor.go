// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/drivers/aio"
	"github.com/stevebargelt/mygobot/platforms/intel-iot/edison"
)

func main() {
	board := edison.NewAdaptor()
	sensor := aio.NewGroveLightSensorDriver(board, "0")

	work := func() {
		sensor.On(sensor.Event("data"), func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
