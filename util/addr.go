package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidPortRange = errors.New("invalid port range")
)

type Address struct {
	Scheme  string
	Host    string
	MinPort int
	MaxPort int
	Path    string
}

func (self *Address) HostPortString(port int) string {
	return fmt.Sprintf("%s:%d", self.Host, port)
}

func (self *Address) String(port int) string {
	if self.Scheme == "" {
		return self.HostPortString(port)
	}
	return fmt.Sprintf("%s://%s:%d%s", self.Scheme, self.Host, port, self.Path)
}

// 格式 scheme://host:minPort~maxPort/path
func ParseAddress(addr string) (addrObj *Address, err error) {
	addrObj = new(Address)

	schemePos := strings.Index(addr, "://")
	if schemePos != -1 {
		addrObj.Scheme = addr[:schemePos]
		addr = addr[schemePos+3:]
	}

	colonPos := strings.Index(addr, ":")
	if colonPos != -1 {
		addrObj.Host = addr[:colonPos]
	}
	addr = addr[colonPos+1:]
	rangePos := strings.Index(addr, "~")
	var minStr, maxStr string
	if rangePos != -1 {
		minStr = addr[:rangePos]
		slashPos := strings.Index(addr, "/")
		if slashPos != -1 {
			maxStr = addr[rangePos+1 : slashPos]
			addrObj.Path = addr[slashPos:]
		} else {
			maxStr = addr[rangePos+1:]
		}
	} else {
		slashPos := strings.Index(addr, "/")
		if slashPos != -1 {
			addrObj.Path = addr[slashPos:]
			minStr = addr[rangePos+1 : slashPos]
		} else {
			minStr = addr[rangePos+1:]
		}
	}
	addrObj.MinPort, err = strconv.Atoi(minStr)
	if err != nil {
		return nil, ErrInvalidPortRange
	}

	if maxStr != "" {
		// extract max port
		addrObj.MaxPort, err = strconv.Atoi(maxStr)
		if err != nil {
			return nil, ErrInvalidPortRange
		}
	} else {
		addrObj.MaxPort = addrObj.MinPort
	}
	return
}

// 将ip和端口合并为地址
func JoinAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func DetectPort(addr string, fn func(a *Address, port int) (interface{}, error)) (interface{}, error) {
	addrObj, err := ParseAddress(addr)
	if err != nil {
		return nil, err
	}

	for port := addrObj.MinPort; port <= addrObj.MaxPort; port++ {
		ln, err := fn(addrObj, port)
		if err == nil {
			return ln, nil
		}
		if port == addrObj.MaxPort {
			return nil, err
		}
	}
	return nil, fmt.Errorf("unable to bind to %s", addr)
}
