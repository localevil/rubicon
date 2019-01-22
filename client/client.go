package main

import "time"

func main() {
	sh := 1
	// addrDev := 18
	// addrSub := 1

	drv := device{}
	drv.connect("tcp", "10.0.40.199", 2000)
	drv.Start()
	go drv.Tick()
	drv.SendCommandToUnit(*newSnid(32, 1975, sh, nil, nil), basicCodesConst|0x1E, 0, []byte{0xFF, 0xFF, 0xFF, 0xFF})
	// //drv.TakeVersion()
	// //drv.TakeFileList()
	// //drv.SendCommandToUnit(snid{sn: 0x2207B700, id: id}, basicCodesConst|0x0A, 0, nil)
	time.Sleep(20 * time.Second)
	// //drv.SendCommandToUnit(snid{sn: 0x2207B700, id: id}, basicCodesConst|0x0B, 0, nil)
	drv.SendCommandToUnit(snid{sn: 0x2207B700, id: 0xA000}, basicCodesConst|0x21, 0, nil)
	drv.TakeFindSnOnAS()
	drv.Stop()
}
