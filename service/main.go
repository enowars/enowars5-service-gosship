package main

import (
	"github.com/gliderlabs/ssh"
)

func main() {
	log.Println("starting...")
	signer, err := GetHostSigner()
	if err != nil {
		log.Fatal(err)
	}
	srv := &ssh.Server{
		Addr: ":2222",
		Handler: func(session ssh.Session) {
			err := session.Exit(0)
			if err != nil {
				log.Error(err)
			}
		},
		HostSigners: []ssh.Signer{signer},
		Version:     "gosship",
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			// allow all public keys
			log.Printf("new connection (%s) with key type: %s\n", ctx.RemoteAddr(), key.Type())
			return true
		},
	}
	log.Println("starting ssh server...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
