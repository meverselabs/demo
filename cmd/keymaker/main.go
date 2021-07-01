package main

import (
	"crypto/elliptic"
	"encoding/hex"
	"log"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/sha3"

	"github.com/fletaio/fleta/common"
	ecrypto "github.com/fletaio/fleta/common/crypto/ethereum/crypto"
	"github.com/fletaio/fleta/common/hash"
	"github.com/fletaio/fleta/common/key"
	"github.com/fletaio/fleta/core/backend"
	_ "github.com/fletaio/fleta/core/backend/buntdb_driver"
	"github.com/fletaio/fleta/core/chain"
	"github.com/fletaio/fleta/core/pile"
)

// Config is a configuration for the cmd
type Config struct {
	SeedNodeMap     map[string]string
	NodeKeyHex      string
	ObserverKeys    []string
	InitGenesisHash string
	InitHash        string
	InitHeight      uint32
	InitTimestamp   uint64
	Port            int
	APIPort         int
	StoreRoot       string
	RLogHost        string
	RLogPath        string
	UseRLog         bool
}

func main() {
	ChainID := uint8(0x01)
	Symbol := "DEMO"
	Usage := "Mainnet"
	Version := uint16(0x0001)

	back, err := backend.Create("buntdb", "./context")
	if err != nil {
		panic(err)
	}
	cdb, err := pile.Open("./chain", hash.Hash256{}, 0, 0)
	if err != nil {
		panic(err)
	}
	st, err := chain.NewStore(back, cdb, ChainID, Symbol, Usage, Version)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		sk, _ := GenKey()
		memkey, _ := key.NewMemoryKeyFromString(sk)
		memkey.PublicKey()
		ph := common.NewPublicHash(memkey.PublicKey())
		log.Println(sk, memkey.PublicKey(), ph)
	}
	for i := uint16(0); i < 10; i++ {
		log.Println(st.NewAddress(0, i))
	}

}

func GenKey() (string, string) {
	privKey, _ := ecrypto.GenerateKey()

	pub := &privKey.PublicKey

	if pub == nil || pub.X == nil || pub.Y == nil {
		panic("nil")
	}
	pubBytes := elliptic.Marshal(btcec.S256(), pub.X, pub.Y)

	b := Keccak256(pubBytes[1:])[12:]

	a := make([]byte, 20)
	if len(b) > 20 {
		b = b[len(b)-20:]
	}
	copy(a[20-len(b):], b)

	return hex.EncodeToString(privKey.D.Bytes()), Hex(a)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func Hex(a []byte) string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}
