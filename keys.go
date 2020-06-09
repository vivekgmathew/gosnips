package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
)

var VClient *api.Client // Global variables !

func InitVault(token string) error {
	conf := &api.Config{
		Address: "http://127.0.0.1:8200",
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return err
	}
	VClient = client

	VClient.SetToken(token)
	return nil
}

type Key struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New() (Key, error) {
	var k Key

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return k, err
	}

	k.publicKey = &privateKey.PublicKey
	k.privateKey = privateKey

	return k, nil
}

func (k Key) PublicKeyToPemString() string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(k.publicKey),
			},
		),
	)
}

func (k Key) PrivateKeyToPemString() string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(k.privateKey),
			},
		),
	)
}

func main() {
	key, err := New()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(key.PublicKeyToPemString())
	fmt.Println(key.PrivateKeyToPemString())

	err = InitVault("s.Nt3vMPQLyR5KyWxVIPAd2hCr")
	if err != nil {
		log.Println(err)
	}
	c := VClient.Logical()

	/*
	secret, err := c.Write("instantprov/new",
		map[string]interface{}{
			"pub":  key.PublicKeyToPemString(),
			"priv": key.PrivateKeyToPemString(),
		})
	if err != nil {
		log.Println(err)
	}
	log.Println(secret)
   */

	log.Println("Reading values")

	secret, err := c.Read("instantprov/new")
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("PRIVATE KEY")
	log.Println("############")
	log.Println(secret.Data["priv"])

	log.Println("PUBLIC KEY")
	log.Println("#############")
	log.Println(secret.Data["pub"])
}
