package main

import (
	"gosship/pkg/host"
	"gosship/pkg/logger"
	"gosship/pkg/utils"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	log := logger.New()
	log.Println("starting...")
	log.Println("loading/generating server key...")
	signer, err := utils.GetHostSigner()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded key with fingerprint: %s\n", gossh.FingerprintLegacyMD5(signer.PublicKey()))
	log.Println("setting up host...")
	h := host.New(log)

	log.Println("creating ssh-server...")
	srv := &ssh.Server{
		Addr:             ":2222",
		Handler:          h.HandleNewSession,
		HostSigners:      []ssh.Signer{signer},
		Version:          "gosship",
		PublicKeyHandler: h.HandlePublicKey,
	}
	log.Println("starting ssh server...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
