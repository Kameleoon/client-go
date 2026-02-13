package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/logging"
)

// Version represents a semantic version with major, minor, and patch components.
type Version struct {
	Major int
	Minor int
	Patch int
}

// VersionFromString parses a version string (e.g., "1.2.3") and returns a Version if no parsing fails.
func NewVersionFromString(versionString string) (*Version, error) {
	versions := []int{0, 0, 0}
	parts := strings.Split(versionString, ".")
	for i := 0; i < len(versions) && i < len(parts); i++ {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			logging.Error("Invalid version component, index: %d, value: '%s'", i, parts[i])
			return nil, errs.NewInternalError("Parsing error")
		}
		versions[i] = val
	}
	return &Version{Major: versions[0], Minor: versions[1], Patch: versions[2]}, nil

}

// CompareTo compares this version with another.
// Returns: -1 if this < other, 0 if equal, 1 if this > other.
func (v *Version) CompareTo(other *Version) int {
	if cmp := compareInt(v.Major, other.Major); cmp != 0 {
		return cmp
	}
	if cmp := compareInt(v.Minor, other.Minor); cmp != 0 {
		return cmp
	}
	return compareInt(v.Patch, other.Patch)
}

// ToFloat returns version as a float32 (major, minor)
func (v *Version) ToFloat() (float32, error) {
	versionString := fmt.Sprintf("%d.%d", v.Major, v.Minor)
	if versionFloat, err := strconv.ParseFloat(versionString, 32); err == nil {
		return float32(versionFloat), nil
	}
	logging.Error(fmt.Sprintf("Failed to parse Version: %v", v))
	return 0.0, errs.NewInternalError(fmt.Sprintf("ToFloat parsing failed: %v", v))
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
