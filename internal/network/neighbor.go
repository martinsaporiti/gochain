package network

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

// FindNeighbors - find neighbors in the network.
func FindNeighbors(myHost string, myPort uint16, startIp uint16, endIp uint16, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)
	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])

	neighbors := make([]string, 0)

	wg := sync.WaitGroup{}
	hosts := make(chan string)
	for port := startPort; port <= endPort; port++ {
		for ip := startIp - 1; ip <= endIp-1; ip++ {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address {
				wg.Add(1)
				go func(guessTarget, guessHost string, port uint16) {
					defer wg.Done()
					if IsFoundHost(guessHost, port) {
						hosts <- guessTarget
					}
				}(guessTarget, guessHost, port)
			}
		}
	}

	go func() {
		wg.Wait()
		close(hosts)
	}()

	for host := range hosts {
		neighbors = append(neighbors, host)
	}

	return neighbors

}

// IsFoundHost checks if a host is found
func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	return err == nil
}

// GetHost - returns the hostname of the machine.
func GetHost() string {
	// TODO: update this function to support docker ips and real ips
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1"
	}
	address, err := net.LookupHost(hostname)
	if err != nil || len(address) <= 2 {
		return "127.0.0.1"
	}
	return address[2]
}
