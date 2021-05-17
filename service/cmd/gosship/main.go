package main

import (
	"gosship/pkg/chat"
	"gosship/pkg/database"
	"gosship/pkg/logger"
	"gosship/pkg/rpc"
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

	//TODO: cleanup
	log.Println("resetting database...")
	db.ResetExceptConfig()

	c := make(chan os.Signal, 1)
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
	log.Printf("loaded key with fingerprint: %s", gossh.FingerprintSHA256(signer.PublicKey()))
	log.Println("setting up host...")
	roomConfig, err := utils.GetRoomConfig(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("rooms: %v", roomConfig.Rooms)
	h := chat.NewHost(log, db, roomConfig.Rooms)
	go h.Serve()

	log.Println("starting grpc server...")
	rpcServer := rpc.NewGRPCServer(log, db, h)
	go rpcServer.Serve()

	log.Println("starting ssh server...")
	srv := &ssh.Server{
		Addr:             ":2222",
		Handler:          h.HandleNewSession,
		HostSigners:      []ssh.Signer{signer},
		Version:          "gosship",
		PublicKeyHandler: h.HandlePublicKey,
		ChannelHandlers: map[string]ssh.ChannelHandler{
			"session": ssh.DefaultSessionHandler,
			"rpc":     rpcServer.Handle,
		},
	}
	log.Printf("listening on %s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
