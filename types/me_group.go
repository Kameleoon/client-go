package types

type MEGroup interface {
	FeatureFlags() []FeatureFlag
	GetFeatureFlagByHash(hash float64) FeatureFlag
}
