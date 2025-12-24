package types

import "fmt"

type DataFile struct {
	FeatureFlags map[string]FeatureFlag
}

func (df DataFile) String() string {
	return fmt.Sprintf("DataFile{FeatureFlags:%v}", df.FeatureFlags)
}
