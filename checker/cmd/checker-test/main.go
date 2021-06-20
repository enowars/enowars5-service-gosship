package main

import (
	"checker/pkg/client"
	"checker/service/database"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func run(signer ssh.Signer) error {
	sshClient, pubKey, err := client.GetSSHClient(context.Background(), "client", "localhost", signer)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	log.Printf("server public key fingerprint: %s", ssh.FingerprintSHA256(pubKey))

	adminClient, ch, err := client.AttachRPCAdminClient(context.Background(), sshClient, pubKey, false)
	if err != nil {
		return err
	}
	defer ch.Execute()

	log.Printf("logged in with %s", adminClient.SessionToken)
	err = adminClient.DumpDirectMessages("chris", func(entry *database.MessageEntry) {
		fmt.Println(entry)
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	user, err := client.GenerateNewUser()
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.NewSignerFromSigner(user.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = run(signer)
	if err != nil {
		log.Fatal(err)
	}
}
