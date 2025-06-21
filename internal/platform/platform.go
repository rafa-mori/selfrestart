package platform

import (
	"runtime"
)

type PlatformInfo struct {
	OS   string
	Arch string
}

func GetHostPlatform() PlatformInfo {
	return PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
