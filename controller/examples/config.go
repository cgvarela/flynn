package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type config struct {
	controllerDomain string
	controllerKey    string
	ourAddr          string
	ourPort          string
}

func loadConfigFromEnv() (*config, error) {
	c := &config{}
	c.controllerDomain = os.Getenv("CONTROLLER_DOMAIN")
	if c.controllerDomain == "" {
		return nil, fmt.Errorf("CONTROLLER_DOMAIN is required")
	}
	c.controllerKey = os.Getenv("CONTROLLER_KEY")
	if c.controllerKey == "" {
		return nil, fmt.Errorf("CONTROLLER_KEY is required")
	}
	c.ourAddr = os.Getenv("ADDR")
	if c.ourAddr == "" {
		err := c.discoverAddr()
		if err != nil {
			return nil, fmt.Errorf("Discovery failed, ADDR is required")
		}
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4456"
	}
	c.ourPort = port
	return c, nil
}

func (c *config) discoverAddr() error {
	controllerAddrs, err := net.LookupHost(c.controllerDomain)
	if err != nil {
		return err
	}
	ints, err := net.Interfaces()
	if err != nil {
		return err
	}
	addrs := make([]string, 0, len(ints))
	for _, i := range ints {
		iAddrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range iAddrs {
			addrs = append(addrs, addr.String())
		}
	}
	var ourAddr string
	for _, cAddr := range controllerAddrs {
		cAddrParts := strings.Split(cAddr, ".")
		if len(cAddrParts) != 4 {
			continue
		}
		for _, iAddr := range addrs {
			iAddrParts := strings.Split(iAddr, ".")
			if len(iAddrParts) != 4 || cAddrParts[0] != iAddrParts[0] || cAddrParts[1] != iAddrParts[1] || cAddrParts[2] != iAddrParts[2] {
				continue
			}
			ourAddr = iAddr
			break
		}
	}
	if ourAddr == "" {
		return fmt.Errorf("No interface found")
	}
	ourAddr = strings.Split(ourAddr, "/")[0]
	c.ourAddr = ourAddr
	return nil
}
