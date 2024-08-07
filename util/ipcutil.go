package util

import (
	ipc "github.com/james-barrow/golang-ipc"
)

func Send(data string) error {
	c, err := ipc.StartClient("sscreenPipe", nil)
	if err != nil {
		return err
	}
	for {
		message, err := c.Read()
		if err == nil {
			if message.MsgType == -1 {
				if c.Status() == "Connected" {
					if err = c.Write(5, []byte(data)); err != nil {
						return err
					}
				}
			} else {
				if string(message.Data) == "R" {
					return nil
				}
			}
		} else {
			return err
		}
	}
}
