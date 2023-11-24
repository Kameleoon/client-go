package types

import "github.com/Kameleoon/client-go/v3/utils"

type Sendable interface {
	QueryEncodable

	Nonce() string

	Sent() bool
	MarkAsSent()
}

type sendableBase struct {
	nonce string
	sent  bool
}

func (sb *sendableBase) Sent() bool {
	return sb.sent
}
func (sb *sendableBase) MarkAsSent() {
	sb.sent = true
	sb.nonce = ""
}

type duplicationSafeSendableBase struct {
	sendableBase
}

func (sb *duplicationSafeSendableBase) initSendale() {
	sb.nonce = utils.GetNonce()
}

func (sb *duplicationSafeSendableBase) Nonce() string {
	return sb.nonce
}

type duplicationUnsafeSendableBase struct {
	sendableBase
}

func (sb *duplicationUnsafeSendableBase) Nonce() string {
	if !sb.sent && (len(sb.nonce) == 0) {
		sb.nonce = utils.GetNonce()
	}
	return sb.nonce
}
