package network

import "github.com/Kameleoon/client-go/v2/utils"

const NonceLength = 16

func GetNonce() string {
	return utils.GetRandomString(NonceLength)
}
