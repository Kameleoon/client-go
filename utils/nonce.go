package utils

const NonceLength = 16

func GetNonce() string {
	return GetRandomString(NonceLength)
}
