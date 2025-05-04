package netutil

import "net"

func FindActiveNics() ([]string, error) {
	nics, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	activeNics := []string{}
	for _, nic := range nics {
		if (nic.Flags&net.FlagUp != 0) && (nic.Flags&net.FlagLoopback == 0) {
			activeNics = append(activeNics, nic.Name)
		}
	}

	return activeNics, nil
}

func IsMyActiveHostIp(ip string) (bool, error) {
	nics, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, nic := range nics {
		if (nic.Flags&net.FlagUp != 0) && (nic.Flags&net.FlagLoopback == 0) {
			addrs, err := nic.Addrs()
			if err != nil {
				return false, err
			}
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.String() == ip {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
