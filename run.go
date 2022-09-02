package main

import (
	"github.com/olwolf/naabu-x6/pkg/runner"
	"github.com/projectdiscovery/goflags"
)

func PortScan(ip string, thread int, ports string) (error, map[string][]int) {
	host := goflags.StringSlice{}
	host.Set(ip)
	options := runner.ParseOptions()
	options.Host = host
	options.Ports = ports
	options.Threads = thread
	newRunner, err := runner.NewRunner(options)
	defer newRunner.Close()
	if err != nil {
		return err, nil
	}
	err, rets := newRunner.RunEnumeration()
	if err != nil {
		return err, nil
	}
	return nil, rets
}

//func main() {
//	err, rets := PortScan("127.0.0.1", 3, "135,445")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%v\n", rets)
//}
