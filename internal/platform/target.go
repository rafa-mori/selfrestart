package platform

import (
	"fmt"
)

type TargetPlatform struct {
	Installed bool
	Version   string
	Platform  string
	Arch      []string
}

var FullPlatformMap = map[string]TargetPlatform{
	"darwin":  {Installed: false, Version: "", Platform: "darwin", Arch: []string{"amd64", "arm64"}},
	"windows": {Installed: false, Version: "", Platform: "windows", Arch: []string{"amd64", "i386"}},
	"linux":   {Installed: false, Version: "", Platform: "linux", Arch: []string{"amd64"}},
}

// GetPlatformTarget returns target platform information based on platform and architecture
func GetPlatformTarget(platform, arch string) (map[string]TargetPlatform, error) {
	hostInfo := GetHostPlatform()

	if platform == "" {
		platform = hostInfo.OS
	}
	if arch == "" {
		arch = hostInfo.Arch
	}

	platformTarget := make(map[string]TargetPlatform)

	if platform == "all" {
		// Return all supported platforms
		for name, target := range FullPlatformMap {
			platformTarget[name] = TargetPlatform{
				Installed: target.Installed,
				Version:   target.Version,
				Platform:  target.Platform,
				Arch:      make([]string, 0),
			}
		}
	} else {
		// Check if platform is supported
		target, exists := FullPlatformMap[platform]
		if !exists {
			return nil, fmt.Errorf("platform %s not supported", platform)
		}

		platformTarget[platform] = TargetPlatform{
			Installed: target.Installed,
			Version:   target.Version,
			Platform:  target.Platform,
			Arch:      make([]string, 0),
		}
	}

	// Validate architecture for each platform
	for pltfrm, target := range platformTarget {
		supportedArch := FullPlatformMap[pltfrm].Arch

		if arch == "all" {
			target.Arch = supportedArch
		} else {
			// Check if architecture is supported
			for _, supportedArch := range supportedArch {
				if supportedArch == arch {
					target.Arch = []string{arch}
					break
				}
			}
		}

		platformTarget[pltfrm] = target
	}

	// Remove platforms without supported architectures
	for pltfrm, target := range platformTarget {
		if len(target.Arch) == 0 {
			delete(platformTarget, pltfrm)
		}
	}

	if len(platformTarget) == 0 {
		return nil, fmt.Errorf("no supported platforms found for %s/%s", platform, arch)
	}

	return platformTarget, nil
}

// IsPlatformSupported checks if a platform/architecture combination is supported
func IsPlatformSupported(platform, arch string) bool {
	target, exists := FullPlatformMap[platform]
	if !exists {
		return false
	}

	for _, supportedArch := range target.Arch {
		if supportedArch == arch {
			return true
		}
	}

	return false
}
