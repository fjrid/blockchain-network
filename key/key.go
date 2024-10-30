package key

import (
	"bytes"
	"crypto/sha256"
	"flag"

	"github.com/libp2p/go-libp2p/core/crypto"
)

type ServerPrivateKey struct {
	challengeKey string
}

func NewPrivateKey() *ServerPrivateKey {
	serverPrivateKey := &ServerPrivateKey{}

	flag.StringVar(&serverPrivateKey.challengeKey, "challenge-key", "BLOCKCHAIN_NETWORK_PRIV_KEY", "Used only for server mode")
	flag.Parse()

	return serverPrivateKey
}

func (spk *ServerPrivateKey) GeneratePrivateKey() (crypto.PrivKey, error) {
	hash := sha256.Sum256([]byte(spk.challengeKey))
	privKey, _, err := crypto.GenerateEd25519Key(bytes.NewReader(hash[:]))
	return privKey, err
}
