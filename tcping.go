/*
* @Author:   mian
* @Date:     2019/10/31 11:43
 */
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func simpleTcp(addrPort string) error {
	conn, err := net.Dial("tcp", addrPort)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}
func simpleUdp(addrPort string) error {
	conn, err := net.Dial("udp", addrPort)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// 配连参数 bind 避免使用匿名函数
func bindErrFunc1(f func(string) error, addPort string) func() error {
	return func() error {
		err := f(addPort)
		return err
	}
}
func getTime(f func() error) (t float32, lost bool) {
	start := time.Now()
	err := f()
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("       ∞ms 丢失\n")
		return 0, true
	} else {
		i := float32(elapsed.Nanoseconds()) / 1e6
		fmt.Printf("%7.2fms 未丢失\n", i)
		return i, false
	}
}
func getNTime(f func() error, n int) (avg float32, min float32, max float32, lostRate float32) {
	var sum float32 = 0
	var reachNum int = 0
	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 500)
		fmt.Printf("%2d: ", i+1)
		t, l := getTime(f)
		if !l {
			reachNum++
			sum += t
			if i == 0 {
				min = t
				max = t
			} else {
				if t > max {
					max = t
				}
				if t < min {
					min = t
				}
			}
		}
	}
	avg = sum / float32(reachNum)
	lostRate = 1.0 - float32(reachNum)/float32(n)
	return
}
func help(s string) {
	println(s)
	println("tcping (-a addr [-p port] [-n testNum] [-pr (tcp||udp)]) || (-h)")
	println("OPTION:")
	flag.PrintDefaults()
	os.Exit(0)
}
func main() {
	h := flag.Bool("h", false, "帮助")
	num := flag.Int("n", 5, "测试次数")
	addr := flag.String("a", "", "地址")
	port := flag.String("p","80", "端口")
	protocol := flag.String("pr","tcp", "协议(tcp/udp)")

	flag.Parse()
	if *h {
		flag.PrintDefaults()
		return
	}
	if *addr == "" {
		help("无地址参数")
	}
	var proF func(string) error
	switch *protocol {
	case "tcp":
		proF = simpleTcp
	case "udp":
		proF = simpleUdp
	case "ping":
		help("请用系统的ping")
	default:
		help("错误协议")
	}

	addrPort := *addr + ":" + *port



	fmt.Printf("对%s的网络测试结果：\n", addrPort)
	a, min, max, lost := getNTime(bindErrFunc1(proF, addrPort), *num)
	fmt.Printf("avg: %.2fms min: %.2fms max: %.2fms lost: %.2f%%\n", a, min, max, lost*100)
}
