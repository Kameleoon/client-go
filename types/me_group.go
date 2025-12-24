package types

type MEGroup interface {
	FeatureFlags() []IFeatureFlag
	GetFeatureFlagByHash(hash float64) IFeatureFlag
}
