package main

import (
	"checker/pkg/client"
	"context"
	"fmt"
	"gosship/pkg/database"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func run(signer ssh.Signer) error {
	sshClient, err := client.GetSSHClient(context.Background(), "client", "localhost", signer)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	rpcChannel, err := client.OpenRPCChannel(sshClient)
	if err != nil {
		return err
	}
	defer rpcChannel.Close()

	grpcConn, err := client.CreateNewGRPCClient(rpcChannel)
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	adminClient := client.NewAdminClient(grpcConn)

	token, err := adminClient.Auth()
	if err != nil {
		return err
	}
	log.Printf("logged in with %s", token)
	err = adminClient.DumpDirectMessages("chris", func(entry *database.MessageEntry) {
		fmt.Println(entry)
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.Println("reading private key...")
	data, err := os.ReadFile("./client-key")
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		log.Fatal(err)
	}

	err = run(signer)
	if err != nil {
		log.Fatal(err)
	}
}
