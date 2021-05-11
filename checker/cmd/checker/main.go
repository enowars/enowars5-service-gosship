package main

import (
	"checker/pkg/client"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func run(signer ssh.Signer) error {
	sshClient, err := client.GetSSHClient(signer)
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
	err = adminClient.SendMessageToRoom("default", "hello from rpc :wave:")
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

	//_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// session stuff
	//session, err := client.NewSession()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//out, err := session.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//in, err := session.StdinPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = session.RequestPty("xterm", 40, 80, ssh.TerminalModes{
	//	ssh.ECHO:  0,
	//	ssh.IGNCR: 1,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = session.Shell()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Fprintf(in, ":clown:\n\r")
	//fmt.Fprintf(in, "/history chris\n\r")
	//done := make(chan struct{})
	//go func() {
	//	outScanner := bufio.NewScanner(out)
	//	for outScanner.Scan() {
	//		fmt.Println(outScanner.Text())
	//	}
	//	close(done)
	//}()
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//<-c
	//log.Println("exiting...")
	//session.Close()
	//<-done
}
