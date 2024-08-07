package driver

import (
	"go.bug.st/serial"
	"golang.org/x/text/encoding/simplifiedchinese"
	"sync"
	"time"
)

var (
	port serial.Port
	lock sync.Mutex
)

func SendOrder(order string) error {
	order = order + "\r\n"
	gbkData, _ := simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(order))
	lock.Lock()
	_, err := port.Write(gbkData)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 200)
	lock.Unlock()
	return nil
}
func InitSerial() error {
	portd, err := serial.Open("COM9", &serial.Mode{
		BaudRate: 115200,
	})
	if err != nil {
		return err
	}
	port = portd
	return nil
}
