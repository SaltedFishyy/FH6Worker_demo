package telemetry

import "net"

func ListNetworkInterfaces() []NetworkInterface {
	result := []NetworkInterface{
		{
			Name:        "all",
			DisplayName: "All interfaces",
			Address:     "0.0.0.0",
			IsLoopback:  false,
			IsPrivate:   true,
			IsUp:        true,
		},
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return result
	}

	for _, iface := range interfaces {
		isUp := iface.Flags&net.FlagUp != 0
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := interfaceIP(addr)
			if ip == nil || ip.To4() == nil {
				continue
			}
			result = append(result, NetworkInterface{
				Name:        iface.Name,
				DisplayName: iface.Name,
				Address:     ip.String(),
				IsLoopback:  ip.IsLoopback(),
				IsPrivate:   ip.IsPrivate(),
				IsUp:        isUp,
			})
		}
	}

	return result
}

func interfaceIP(addr net.Addr) net.IP {
	switch value := addr.(type) {
	case *net.IPNet:
		return value.IP
	case *net.IPAddr:
		return value.IP
	default:
		return nil
	}
}
