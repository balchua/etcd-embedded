package util

import (
	"net"

	"go.uber.org/zap"
)

func GetDefaultIPV4() (string, error) {
	lg, _ := zap.NewProduction()
	ietfs, _ := net.Interfaces()

	for _, ietf := range ietfs { // get ipv4 address

		if !isInterfaceUp(&ietf) {
			lg.Info("Interface is down", zap.String("name", ietf.Name))
			continue
		}

		if isLoopbackOrPointToPoint(&ietf) {
			lg.Info("Interface is a loopback", zap.String("name", ietf.Name))
			continue
		}

		if isInExcludedInterfaces(&ietf) {
			lg.Info("Interface is in excluded list", zap.String("name", ietf.Name))
			continue
		}

		lg.Info("Interface details", zap.String("name", ietf.Name), zap.Bool("isUp", isInterfaceUp(&ietf)))
		addrs, _ := ietf.Addrs()

		for _, addr := range addrs {
			lg.Info("Interface name")
			ipv4Addr := addr.(*net.IPNet).IP.To4()
			if ipv4Addr != nil {
				lg.Info("IP found", zap.String("ip", ipv4Addr.String()))
				return ipv4Addr.String(), nil
			}

		}
	}
	return "", nil

}

func isLoopbackOrPointToPoint(intf *net.Interface) bool {
	return intf.Flags&(net.FlagLoopback|net.FlagPointToPoint) != 0
}

func isInterfaceUp(intf *net.Interface) bool {
	if intf == nil {
		return false
	}
	if intf.Flags&net.FlagUp != 0 {
		return true
	}
	return false
}

func isIp4(ip net.IP) bool {
	if ip.To4() != nil {
		return true
	}
	return false
}

func isInExcludedInterfaces(intf *net.Interface) bool {
	if intf.Name == "vxlan" {
		return true
	}
	return false
}
