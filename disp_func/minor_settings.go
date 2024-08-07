package disp_func

import (
	"fmt"
	"sscreen/driver"
)

func SetBrightness(brightness int) {
	driver.SendOrder(fmt.Sprintf("BL(%d);", brightness))
}
