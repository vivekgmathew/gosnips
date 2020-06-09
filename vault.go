package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

var token = "s.Nt3vMPQLyR5KyWxVIPAd2hCr"
var vault_addr = "http://localhost:8200"

func main() {
	config := &api.Config{
		Address: vault_addr,
	}
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.SetToken(token)
	secret, err := client.Logical().Read("instantprov/bt")
	if err != nil {
		fmt.Println(err)
		return
	}
	m, ok := secret.Data["priv_key"].(map[string]interface{})
	if !ok {
		fmt.Printf("%T %#v\n", secret.Data["data"], secret.Data["data"])
		return
	}
	fmt.Printf("hello: %v\n", m["hello"])
}
