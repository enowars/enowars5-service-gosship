package main

import (
	"context"
	"gosship/pkg/chat"
	"gosship/pkg/database"
	"gosship/pkg/logger"
	"gosship/pkg/rpc"
	"gosship/pkg/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
)

func run(log *logrus.Logger) error {
	log.Println("starting...")
	log.Println("opening database...")
	db, err := database.NewDatabase(log)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("loading/generating server key...")
	signer, err := utils.GetHostSigner(db)
	if err != nil {
		return err
	}

	log.Printf("loaded key with fingerprint: %s", gossh.FingerprintSHA256(signer.PublicKey()))

	log.Println("setting up host...")
	roomConfig, err := utils.GetRoomConfig(db)
	if err != nil {
		return err
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

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err = srv.ListenAndServe(); err != ssh.ErrServerClosed {
			log.Error(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	log.Println("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.Announcement(":no_entry_sign:server is shutting down :no_entry_sign:")
	err = srv.Shutdown(ctx)
	if err == context.DeadlineExceeded {
		log.Println("force closing all active connections...")
		if err := srv.Close(); err != nil {
			log.Error(err)
		}
		log.Println("finishing pending database writes...")
		<-time.After(time.Second)
	} else if err != nil {
		log.Error(err)
	}

	return nil
}

func main() {
	log := logger.New()
	if err := run(log); err != nil {
		log.Fatal(err)
	}
}
