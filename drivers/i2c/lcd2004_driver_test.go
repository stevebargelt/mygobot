package i2c

import (
	"errors"
	"strings"
	"testing"

	gobot "github.com/stevebargelt/mygobot"
	"github.com/stevebargelt/mygobot/gobottest"
)

var _ gobot.Driver = (*LCD2004Driver)(nil)

// --------- HELPERS
func initTestLCD2004Driver() (driver *LCD2004Driver) {
	driver, _ = initTestLCD2004DriverWithStubbedAdaptor()
	return
}

func initTestLCD2004DriverWithStubbedAdaptor() (*LCD2004Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewLCD2004Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewLCD2004Driver(t *testing.T) {
	// Does it return a pointer to an instance of LCD2004Driver?
	var mpl interface{} = NewLCD2004Driver(newI2cTestAdaptor())
	_, ok := mpl.(*LCD2004Driver)
	if !ok {
		t.Errorf("NewLCD2004Driver() should have returned a *LCD2004Driver")
	}
}

// Methods
func TestLCD2004Driver(t *testing.T) {
	jhd := initTestLCD2004Driver()

	gobottest.Refute(t, jhd.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(jhd.Name(), "LCD2004"), true)
}

func TestLCD2004DriverSetName(t *testing.T) {
	d := initTestLCD2004Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestLCD2004DriverOptions(t *testing.T) {
	d := NewLCD2004Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestLCD2004DriverStart(t *testing.T) {
	d := initTestLCD2004Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestLCD2004StartConnectError(t *testing.T) {
	d, adaptor := initTestLCD2004DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestLCD2004DriverStartWriteError(t *testing.T) {
	d, adaptor := initTestLCD2004DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestLCD2004DriverHalt(t *testing.T) {
	d := initTestLCD2004Driver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestLCD2004DriverClear(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Clear(), nil)
}

func TestLCD2004DriverClearError(t *testing.T) {
	d, a := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Clear(), errors.New("write error"))
}

func TestLCD2004DriverHome(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Home(), nil)
}

func TestLCD2004DriverWrite(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Write("Hello"), nil)
}

func TestLCD2004DriverWriteError(t *testing.T) {
	d, a := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	gobottest.Assert(t, d.Write("Hello"), errors.New("write error"))
}

func TestLCD2004DriverWriteTwoLines(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Write("Hello\nthere"), nil)
}

func TestLCD2004DriverWriteTwoLinesError(t *testing.T) {
	d, a := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Write("Hello\nthere"), errors.New("write error"))
}

// func TestLCD2004DriverSetPosition(t *testing.T) {
// 	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
// 	d.Start()
// 	gobottest.Assert(t, d.SetPosition(2), nil)
// }

// func TestLCD2004DriverSetSecondLinePosition(t *testing.T) {
// 	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
// 	d.Start()
// 	gobottest.Assert(t, d.SetPosition(18), nil)
// }

// func TestLCD2004DriverSetPositionInvalid(t *testing.T) {
// 	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
// 	d.Start()
// 	gobottest.Assert(t, d.SetPosition(-1), ErrInvalidPosition)
// 	gobottest.Assert(t, d.SetPosition(32), ErrInvalidPosition)
// }

func TestLCD2004DriverScroll(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Scroll(true), nil)
}

func TestLCD2004DriverReverseScroll(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Scroll(false), nil)
}

// func TestLCD2004DriverSetCustomChar(t *testing.T) {
// 	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
// 	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
// 	d.Start()
// 	gobottest.Assert(t, d.SetCustomChar(0, data), nil)
// }

// func TestLCD2004DriverSetCustomCharError(t *testing.T) {
// 	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
// 	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
// 	d.Start()
// 	gobottest.Assert(t, d.SetCustomChar(10, data), errors.New("can't set a custom character at a position greater than 7"))
// }

// func TestLCD2004DriverSetCustomCharWriteError(t *testing.T) {
// 	d, a := initTestLCD2004DriverWithStubbedAdaptor()
// 	d.Start()

// 	a.i2cWriteImpl = func([]byte) (int, error) {
// 		return 0, errors.New("write error")
// 	}
// 	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
// 	gobottest.Assert(t, d.SetCustomChar(0, data), errors.New("write error"))
// }

func TestLCD2004DriverCommands(t *testing.T) {
	d, _ := initTestLCD2004DriverWithStubbedAdaptor()
	d.Start()

	err := d.Command("Clear")(map[string]interface{}{})
	gobottest.Assert(t, err, nil)

	err = d.Command("Home")(map[string]interface{}{})
	gobottest.Assert(t, err, nil)

	err = d.Command("Write")(map[string]interface{}{"msg": "Hello"})
	gobottest.Assert(t, err, nil)

	// err = d.Command("SetPosition")(map[string]interface{}{"pos": "1"})
	// gobottest.Assert(t, err, nil)

	// err = d.Command("Scroll")(map[string]interface{}{"lr": "true"})
	// gobottest.Assert(t, err, nil)
}
