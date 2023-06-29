package models

import (
	"net"

	"strings"
)

func InSliceStr(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

func InSliceStrFuzzy(str string, slice []string) bool {
	for _, s := range slice {
		if strings.Contains(str, s) || strings.Contains(s, str) {
			return true
		}
	}
	return false
}

// IPCheck 检查 ip 是否在特定的 ip 地址范围内
func IPCheck(thisip string, ips []string) bool {
	for _, ip := range ips {
		ip = strings.TrimRight(ip, "/")
		if strings.Contains(ip, "/") {
			if ipCheckMask(thisip, ip) {
				return true
			}
		} else {
			if thisip == ip {
				return true
			}
		}
	}
	return false
}

func ipCheckMask(ip, ipMask string) bool {
	_, subnet, _ := net.ParseCIDR(ipMask)

	thisIP := net.ParseIP(ip)
	return subnet.Contains(thisIP)
}
