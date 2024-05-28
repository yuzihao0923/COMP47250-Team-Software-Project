package ip

import (
	"fmt"
	"net"
	"os"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// Check the ip address to determine if it is a loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addr := ipnet.IP.String()
				return addr, nil
				// fmt.Println(ipnet.IP.String())
			}
		}
	}
	return "", fmt.Errorf("no non-loopback IPv4 address found")
}
