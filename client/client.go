package main

func main() {
	drv := device{}
	drv.connect("tcp", "10.0.40.199", 2000)
	drv.TakeVersion()
	drv.Stop()
}
