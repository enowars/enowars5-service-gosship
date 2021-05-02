package main

import (
	"gosship/pkg/chat"
	"gosship/pkg/database"
	"gosship/pkg/logger"
	"gosship/pkg/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	log := logger.New()
	log.Println("starting...")
	log.Println("opening database...")
	db, err := database.NewDatabase(log)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("closing database...")
		db.Close()
		os.Exit(0)
	}()

	log.Println("loading/generating server key...")
	signer, err := utils.GetHostSigner(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded key with fingerprint: %s\n", gossh.FingerprintLegacyMD5(signer.PublicKey()))
	log.Println("setting up host...")
	h := chat.NewHost(log, db)
	go h.Serve()

	log.Println("starting ssh server...")
	srv := &ssh.Server{
		Addr:             ":2222",
		Handler:          h.HandleNewSession,
		HostSigners:      []ssh.Signer{signer},
		Version:          "gosship",
		PublicKeyHandler: h.HandlePublicKey,
	}
	log.Printf("listening on %s\n", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
