package blockchain

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Blockchain struct{}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getPrivateKey(index uint32) (*ecdsa.PrivateKey, error) {
	// Get the mnemonic from the environment
	mnemonic := goDotEnvVariable("MNEMONIC")
	if mnemonic == "" {
		return nil, errors.New("MNEMONIC not set in environment")
	}

	// Generate the seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate a master private key from the seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Derive a new private key for the user using the provided index
	childKey, err := masterKey.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	// Convert the child key to an ECDSA private key
	privateKey, err := crypto.ToECDSA(childKey.Key)
	if err != nil {
		return nil, err
	}

	return privateKey, nil

}

func (b *Blockchain) CreateDepositAddress(userID string, index uint32) (string, error) {
	privateKey, err := getPrivateKey(index)

	// Generate the public address from the private key
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Return the hexadecimal representation of the public address
	depositAddress := publicAddress.Hex()
	return depositAddress, err
}

func (b *Blockchain) GetPrivateKey(index uint32) (*ecdsa.PrivateKey, error) {
	// Get the mnemonic from the environment
	mnemonic := goDotEnvVariable("MNEMONIC")
	if mnemonic == "" {
		return nil, errors.New("MNEMONIC not set in environment")
	}

	// Generate the seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate a master private key from the seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Derive a new private key for the user using the provided index
	childKey, err := masterKey.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	// Convert the child key to an ECDSA private key
	privateKey, err := crypto.ToECDSA(childKey.Key)
	if err != nil {
		return nil, err
	}

	return privateKey, nil

}
