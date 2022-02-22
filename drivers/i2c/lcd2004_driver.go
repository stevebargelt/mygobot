package i2c

import (
	"time"

	gobot "github.com/stevebargelt/mygobot"
)

// Followed:
// https://sourcegraph.com/github.com/dotnet/iot@12a9297e0dda6593f0479543da0f9cd5f3768c4f/-/blob/src/devices/CharacterLcd/Hd44780.cs?L104
// and
// https://sourcegraph.com/github.com/dotnet/iot@12a9297e0dda6593f0479543da0f9cd5f3768c4f/-/blob/src/devices/CharacterLcd/Flags.cs?L77

const (
	COMMAND_MODE				= 0x80
	DATA_MODE					= 0x40

	LCD2004_CLEARDISPLAY        = 0x01
	LCD2004_RETURNHOME          = 0x02
	LCD2004_SETCGRAMADDR        = 0x40
	LCD2004_SETDDRAMADDR        = 0x80
	
	//Display Entry Mode
	LCD2004_DISPLAYSHIFT		= 0x01
	LCD2004_ENTRYSHIFTINCREMENT = 0x02
	LCD2004_ENTRYMODESET        = 0x04
	
	//Display Function
	LCD2004_EXTINSTRUCTION		= 0x01
	LCD2004_FONT5x10			= 0x04
	LCD2004_TWOLINE				= 0x08
	LCD2004_EIGHTBIT			= 0x10
	LCD2004_FUNCCOMMAND			= 0x20	
	
	//Display Control
	LCD2004_BLINKON             = 0x01
	LCD2004_CURSORON            = 0x02
	LCD2004_DISPLAYON           = 0x04
	LCD2004_DISPLAYCONTROLSET   = 0x08
	
	//Display Shift
	LCD2004_MOVERIGHT           = 0x04
	LCD2004_DISPLAYMOVE         = 0x08
	LCD2004_DISPLAYSHIFTCOMMAND = 0x10
	
	//Offsets
	LCD2004_2NDLINEOFFSET = 0x40
)

// LCD2004Driver is a driver for the LCD2004 LCD display
type LCD2004Driver struct {
	name      string
	connector Connector
	Config
	gobot.Commander
	lcdAddress    int
	lcdConnection Connection
	rows int
	cols int
	rowOffsets [4]int
}

// NewLCD2004Driver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
func NewLCD2004Driver(a Connector, options ...func(Config)) *LCD2004Driver {
	j := &LCD2004Driver{
		name:       gobot.DefaultName("LCD2004"),
		connector:  a,
		Config:     NewConfig(),
		Commander:  gobot.NewCommander(),
		lcdAddress: 0x27,
		rows: 4,
		cols: 20,
	}

	j.rowOffsets[0] = 0x00
	j.rowOffsets[1] = LCD2004_2NDLINEOFFSET
	j.rowOffsets[2] = 0x00 + j.cols
	j.rowOffsets[3] = LCD2004_2NDLINEOFFSET + j.cols

	for _, option := range options {
		option(j)
	}

	j.AddCommand("Clear", func(params map[string]interface{}) interface{} {
		return j.Clear()
	})
	j.AddCommand("Home", func(params map[string]interface{}) interface{} {
		return j.Home()
	})
	j.AddCommand("Write", func(params map[string]interface{}) interface{} {
		msg := params["msg"].(string)
		return j.Write(msg)
	})
	// j.AddCommand("SetPosition", func(params map[string]interface{}) interface{} {
	// 	pos, _ := strconv.Atoi(params["pos"].(string))
	// 	return j.SetPosition(pos)
	// })
	// j.AddCommand("Scroll", func(params map[string]interface{}) interface{} {
	// 	lr, _ := strconv.ParseBool(params["lr"].(string))
	// 	return j.Scroll(lr)
	// })

	return j
}

// Name returns the name the LCD2004 Driver was given when created.
func (h *LCD2004Driver) Name() string { return h.name }

// SetName sets the name for the LCD2004 Driver.
func (h *LCD2004Driver) SetName(n string) { h.name = n }

// Connection returns the driver connection to the device.
func (h *LCD2004Driver) Connection() gobot.Connection {
	return h.connector.(gobot.Connection)
}

// Start starts the backlit and the screen and initializes the states.
func (h *LCD2004Driver) Start() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())

	// While the chip supports 5x10 pixel characters for one line displays they
    // don't seem to be generally available. Supporting 5x10 would require extra
    // support for CreateCustomCharacter
	if h.lcdConnection, err = h.connector.GetConnection(h.lcdAddress, bus); err != nil {
		return err
	}

	// SEE PAGE 45/46 FOR INITIALIZATION SPECIFICATION!
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way befer 4.5V so we'll wait 50
	time.Sleep(50 * time.Millisecond)

	// Function must be set first to ensure that we always have the basic
	// instruction set selected. (See PCF2119x datasheet Function_set note
	// for one documented example of where this is necessary.
	init_payload := []byte{COMMAND_MODE, LCD2004_FUNCCOMMAND | LCD_2LINE}
	if _, err := h.lcdConnection.Write(init_payload); err != nil {
		return err
	}
	
	time.Sleep(1 * time.Millisecond)	
	if _, err := h.lcdConnection.Write([]byte{COMMAND_MODE, LCD_DISPLAYCONTROL | LCD_DISPLAYON}); err != nil {
		return err
	}

	time.Sleep(1 * time.Millisecond)
	if _, err := h.lcdConnection.Write([]byte{COMMAND_MODE, LCD_ENTRYMODESET | LCD_ENTRYSHIFTINCREMENT}); err != nil {
		return err
	}

	time.Sleep(1 * time.Millisecond)
	if err := h.Clear(); err != nil {
		return err
	}

	return nil
}

// Clear clears the text on the lCD display.
func (h *LCD2004Driver) Clear() error {
	err := h.command([]byte{LCD_CLEARDISPLAY})
	return err
}

// Home sets the cursor to the origin position on the display.
func (h *LCD2004Driver) Home() error {
	err := h.command([]byte{LCD_RETURNHOME})
	// This wait fixes a race condition when calling home and clear back to back.
	time.Sleep(2 * time.Millisecond)
	return err
}

// Write displays the passed message on the screen.
func (h *LCD2004Driver) Write(message string) error {
	// This wait fixes an odd bug where the clear function doesn't always work properly.
	time.Sleep(1 * time.Millisecond)
	
	for _, val := range message {
		// TODO: Does \n really not work? 
		// if val == '\n' {
		// 	if err := h.SetPosition(16); err != nil {
		// 		return err
		// 	}
		// 	continue
		// }
		if _, err := h.lcdConnection.Write([]byte{DATA_MODE, byte(val)}); err != nil {
			return err
		}
	}
	return nil
}

// func (h *LCD2004Driver) SetPosition(left int, top int) (err error) {

// 	//h.rowOffsets := []int{ 0, 64, 20, 84 }

// 	if top < 0 || top >= h.rows {
// 		err = ErrInvalidPosition
// 		return
// 	}
// 	var newAddress = left + h.rowOffsets[top];

// 	if pos < 0 || pos > 31 {
// 		err = ErrInvalidPosition
// 		return
// 	}
// 	offset := byte(pos)
// 	if pos >= 16 {
// 		offset -= 16
// 		offset |= LCD_2NDLINEOFFSET
// 	}
// 	err = h.command([]byte{LCD_SETDDRAMADDR | offset})
// 	return
// }

// Scroll sets the scrolling direction for the display, either left to right, or
// right to left.
func (h *LCD2004Driver) Scroll(lr bool) error {
	if lr {
		_, err := h.lcdConnection.Write([]byte{COMMAND_MODE, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT})
		return err
	}

	_, err := h.lcdConnection.Write([]byte{COMMAND_MODE, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT})
	return err
}

func (h *LCD2004Driver) Halt() error { return nil }

func (h *LCD2004Driver) command(buf []byte) error {
	_, err := h.lcdConnection.Write(append([]byte{COMMAND_MODE}, buf...))
	return err
}

func (h *LCD2004Driver) commandAndWait(buf []byte) error {
	err := h.command(buf)
	time.Sleep(1 * time.Millisecond)
	return err
}

