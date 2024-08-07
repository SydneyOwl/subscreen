package disp_func

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/sys/windows/registry"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"sscreen/driver"
	"strings"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup
var ctx context.Context
var cancel context.CancelFunc

func InitScreen() {
	driver.SendOrder("DIR(1);BL(190);CLR(18);PL(0,30,600,30,0);DCV16(5,5,'WIFI CONNECTING',0);PL(360,30,360,150);")
	driver.SendOrder("PL(0,180,155,180);DCV16(22,190,'Outside Sensor',0);DCV16(0,220,'Temperature:--',0);DCV16(0,240,'Humidity:--%',0);DCV16(0,280,'Batt:--%',0);DCV16(0,260,'Signal:--dbm',0);")
	driver.SendOrder("DCV16(25,40,'Inside Sensor',0);DCV16(0,70,'Temperature:--',0);DCV16(0,90,'Humidity:--%',0);DCV16(0,130,'Batt:--%',0);DCV16(0,110,'Signal:--dbm',0);PL(155,150,480,150);")
	driver.SendOrder("PL(155,30,155,350);DCV24(255,170,'2024/02/11',0);DCV32(250,210,'11:22:33',0);DCV24(270,260,'-',0)")
	driver.SendOrder("DCV16(170,40,'System Status',0);DCV16(350,5,'Uptime: 11h22m',0);DCV24(170,70,'CPU:100% @ 120C',0);DCV24(170,110,'MEM:100%',0)")
	driver.SendOrder("DCV16(400,40,'U/R',0);DCV16(450,40,'D/W',0);DCV16(370,80,'DS',0);DCV16(370,120,'NT',0);")
	driver.SendOrder("DCV16(400,80,'100',0);DCV16(450,80,'110',0);DCV16(400,120,'1000',0);DCV16(450,120,'999',0);")
}

func UpdateNews(ctx context.Context, wag *sync.WaitGroup) {
	for {
		fmt.Println("Getting new news...")
		postValue := url.Values{"key": {"SECRET"}}
		res, err := http.PostForm("https://apis.tianapi.com/bulletin/index", postValue)
		if err != nil {
			fmt.Println("Failed and trying Getting new news...")
			time.Sleep(time.Second * 2)
			continue
		}
		defer res.Body.Close()
		tianapi_data, _ := io.ReadAll(res.Body)
		jsonVal := gjson.Parse(string(tianapi_data))
		total := jsonVal.Get("result.list.#").Int()
		for number := range total {
			dig := jsonVal.Get(fmt.Sprintf("result.list.%d.digest", number)).String()
			pData := fmt.Sprintf("DCV16(5,5,'%s                                           ',0);PL(0,30,600,30,0);", dig)
			driver.SendOrder(pData)
			select {
			case <-time.After(time.Minute * 2):
				continue
			case <-ctx.Done():
				fmt.Println("Exiting routine: UpdateNews")
				wag.Done()
				return
			}
		}
	}
}

func UpdateTime(ctx context.Context, wag *sync.WaitGroup) {
	for {
		curDate := time.Now().Format("2006-01-02")
		curTime := time.Now().Format("15:04:05")
		utcDate := time.Now().UTC().Format("2006-01-02 15:04:05")
		utcTime := time.Now().UTC().Format("")
		driver.SendOrder(fmt.Sprintf("DCV24(180,290,'UTC: %s %s',0);DCV24(255,170,'%s',0);DCV32(250,210,'%s',0);DCV24(270,260,'%s',0)", utcDate, utcTime, curDate, curTime, time.Now().Weekday().String()))
		select {
		case <-time.After(time.Second * 1):
			continue
		case <-ctx.Done():
			fmt.Println("Exiting routine: UpdateTime")
			wag.Done()
			return
		}
	}
}

func UpdateSensor(ctx context.Context, wag *sync.WaitGroup) {
	for {
		fmt.Println("Getting new sensorData")
		res, err := http.Get("http://192.168.1.60:3649/getSensor")
		if err != nil {
			driver.SendOrder("DCV16(22,190,'Outside Sensor',1);DCV16(25,40,'Inside Sensor',1);")
			time.Sleep(time.Second * 2)
			continue
		}
		defer res.Body.Close()
		data, _ := io.ReadAll(res.Body)
		jsonVal := gjson.Parse(string(data))
		//in->pout
		driver.SendOrder(fmt.Sprintf("DCV16(0,70,'Temperature:%s  ',0);DCV16(0,90,'Humidity:%s%%  ',0);DCV16(0,130,'Batt:%d%%  ',0);DCV16(0,110,'Signal:%ddbm  ',0);DCV16(0,150,'%s',0);",
			jsonVal.Get("0.temp").String(), jsonVal.Get("0.humi").String(), jsonVal.Get("0.batt").Int(), jsonVal.Get("0.rssi").Int(), strings.ReplaceAll(jsonVal.Get("0.time").String(), " ", "")))
		driver.SendOrder(fmt.Sprintf("DCV16(0,220,'Temperature:%s  ',0);DCV16(0,240,'Humidity:%s%%  ',0);DCV16(0,280,'Batt:%d%%  ',0);DCV16(0,260,'Signal:%ddbm  ',0);DCV16(0,300,'%s',0);",
			jsonVal.Get("1.temp").String(), jsonVal.Get("1.humi").String(), jsonVal.Get("1.batt").Int(), jsonVal.Get("1.rssi").Int(), strings.ReplaceAll(jsonVal.Get("1.time").String(), " ", "")))
		select {
		case <-time.After(time.Minute):
			continue
		case <-ctx.Done():
			fmt.Println("Exiting routine: UpdateSensor")
			wag.Done()
			return
		}
	}
}

func UpdateSysInfo(ctx context.Context, wag *sync.WaitGroup) {
	for {
		key, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\FinalWire\\AIDA64\\SensorValues", registry.QUERY_VALUE)
		if err != nil {
			log.Fatal(err)
			return
		}
		valUptime, _, _ := key.GetStringValue("Value.SUPTIMENS")
		cmd := exec.Command("powershell.exe", "wmic cpu get loadpercentage")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
		// 执行命令，并返回结果
		output, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		res := strings.TrimSpace(strings.Split(string(output), "\n")[1])
		valMem, _, _ := key.GetStringValue("Value.SMEMUTI")
		valCpuTemp, _, _ := key.GetStringValue("Value.TCPUPKG")
		valDiskReadSpeed, _, _ := key.GetStringValue("Value.SDSK1READSPD")
		valDiskWriteSpeed, _, _ := key.GetStringValue("Value.SDSK1WRITESPD")
		valNetDownloadSpeed, _, _ := key.GetStringValue("Value.SNIC3DLRATE")
		valNetUploadSpeed, _, _ := key.GetStringValue("Value.SNIC3ULRATE")
		key.Close()
		driver.SendOrder(fmt.Sprintf("DCV16(350,5,'Uptime: %s  ',0);DCV24(170,70,'CPU: %s%% @ %sC  ',0);DCV24(170,110,'MEM: %s%%   ',0);", valUptime, res, valCpuTemp, valMem))
		// up-dsk do-dsk up-net do-net
		driver.SendOrder(fmt.Sprintf("PL(360,30,360,150);DCV16(400,80,'%s  ',0);DCV16(450,80,'%s   ',0);DCV16(400,120,'%s',0);DCV16(450,120,'%s',0);", valDiskReadSpeed, valDiskWriteSpeed, valNetUploadSpeed, valNetDownloadSpeed))
		//
		select {
		case <-time.After(time.Second * 2):
			continue
		case <-ctx.Done():
			fmt.Println("Exiting routine: UpdateSysInfo")
			wag.Done()
			return
		}
	}
}

func StartAll() {
	ctx, cancel = context.WithCancel(context.Background())
	wg.Add(4)
	go UpdateTime(ctx, &wg)
	go UpdateSysInfo(ctx, &wg)
	go UpdateSensor(ctx, &wg)
	go UpdateNews(ctx, &wg)
}

func StopAll() {
	cancel()
	wg.Wait()
}
