package conditions

import (
	"fmt"
	"strconv"
	"strings"
)

func GetMajorMinorPatch(version string) (major, minor, patch int, err error) {
	versions := []int{major, minor, patch}
	parts := strings.Split(version, ".")
	for i := 0; i < len(parts) && i < len(versions); i++ {
		versions[i], err = strconv.Atoi(parts[i])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid version component: %v", err)
		}
	}

	return versions[0], versions[1], versions[2], nil
}

func GetMajorMinorAsFloat(version string) (float32, error) {
	major, minor, _, err := GetMajorMinorPatch(version)
	if err != nil {
		return 0.0, err
	}
	versionString := fmt.Sprintf("%d.%d", major, minor)
	versionFloat, err := strconv.ParseFloat(versionString, 32)
	if err != nil {
		return 0.0, fmt.Errorf("failed to convert version to float: %v", err)
	}
	return float32(versionFloat), nil
}
