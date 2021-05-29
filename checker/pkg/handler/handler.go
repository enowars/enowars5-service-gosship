package handler

import (
	"bufio"
	"checker/pkg/checker"
	"checker/pkg/client"
	"checker/pkg/database"
	"checker/pkg/quotes"
	gsDatabase "checker/service/database"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var serviceInfo = &checker.InfoMessage{
	ServiceName:     "gosship",
	FlagVariants:    2,
	NoiseVariants:   1,
	HavocVariants:   1,
	ExploitVariants: 2,
}

var ErrVariantIdOutOfRange = errors.New("variantId out of range")
var ErrVariantNotFound = errors.New("variant not found")
var ErrInvalidVariant = errors.New("invalid variant database entry")
var ErrResponseNotFoundTimeout = errors.New("the response was not received after a certain timeout")

type Handler struct {
	log *logrus.Logger
	db  *database.Database
}

func New(log *logrus.Logger, db *database.Database) *Handler {
	return &Handler{
		log: log,
		db:  db,
	}
}

func (h *Handler) sendMessageAndCheckResponse(ctx context.Context, sessIo *client.SessionIO, message, check string) error {
	_, err := fmt.Fprintf(sessIo, "%s\n\r", message)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	scanner := bufio.NewScanner(sessIo)
	go func() {
		for scanner.Scan() {
			txt := stripansi.Strip(scanner.Text())
			if strings.Contains(txt, check) {
				time.Sleep(time.Millisecond * 100)
				break
			}
		}
		if err := scanner.Err(); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	// we just wait until the context deadline
	//case <-time.After(time.Second * 3):
	//	return ErrResponseNotFoundTimeout
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (h *Handler) sendDirectMessage(ctx context.Context, userA *client.User, userB *client.User, addr, msg string) error {
	_, _, chA, err := client.CreateSSHSession(ctx, userA.Name, addr, userA.PrivateKey)
	if err != nil {
		return err
	}
	chA.Execute()

	_, sessionIO, ch, err := client.CreateSSHSession(ctx, userB.Name, addr, userB.PrivateKey)
	if err != nil {
		return err
	}
	defer ch.Execute()

	directMessage := fmt.Sprintf("/dm %s %s", userA.Name, msg)
	checkStr := fmt.Sprintf("[dm][%s]:", userB.Name)
	return h.sendMessageAndCheckResponse(ctx, sessionIO, directMessage, checkStr)
}

func (h *Handler) sendPrivateRoomMessage(ctx context.Context, userA *client.User, addr, room, password, msg string) error {
	_, sessionIO, ch, err := client.CreateSSHSession(ctx, userA.Name, addr, userA.PrivateKey)
	if err != nil {
		return err
	}
	defer ch.Execute()

	createRoom := fmt.Sprintf("/create %s %s", room, password)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, createRoom, room+" was created")
	if err != nil {
		return err
	}

	err = h.sendMessageAndCheckResponse(ctx, sessionIO, msg, msg)
	if err != nil {
		return err
	}

	return h.sendMessageAndCheckResponse(ctx, sessionIO, "/j", "you are now in room default.")
}

func (h *Handler) putFlagDirectMessage(ctx context.Context, message *checker.TaskMessage) (*checker.HandlerInfo, error) {
	userA, err := client.GenerateNewUser()
	if err != nil {
		return nil, err
	}
	userB, err := client.GenerateNewUser()
	if err != nil {
		return nil, err
	}

	err = h.sendDirectMessage(ctx, userA, userB, message.Address, message.Flag)
	if err != nil {
		return nil, err
	}

	err = h.db.PutEntry(&database.Entry{
		Type:        "flag",
		Variant:     "dm",
		TaskMessage: message,
		UserA:       userA,
		UserB:       userB,
	})

	if err != nil {
		return nil, err
	}

	return checker.NewPutFlagInfo(userA.Name), nil
}

func (h *Handler) putFlagPrivateRoom(ctx context.Context, message *checker.TaskMessage) (*checker.HandlerInfo, error) {
	userA, err := client.GenerateNewUser()
	if err != nil {
		return nil, err
	}

	room, password := client.GenerateRoomAndPassword()
	err = h.sendPrivateRoomMessage(ctx, userA, message.Address, room, password, message.Flag)
	if err != nil {
		return nil, err
	}

	err = h.db.PutEntry(&database.Entry{
		Type:        "flag",
		Variant:     "room",
		TaskMessage: message,
		UserA:       userA,
		Room:        room,
		Password:    password,
	})

	if err != nil {
		return nil, err
	}

	return checker.NewPutFlagInfo(room), nil
}

func (h *Handler) PutFlag(ctx context.Context, message *checker.TaskMessage) (*checker.HandlerInfo, error) {
	if message.VariantId >= serviceInfo.FlagVariants {
		return nil, ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return h.putFlagDirectMessage(ctx, message)
	case 1:
		return h.putFlagPrivateRoom(ctx, message)
	}

	return nil, ErrVariantNotFound
}

func (h *Handler) getFlagDirectMessage(ctx context.Context, message *checker.TaskMessage) error {
	fi, err := h.db.GetEntry(message.TaskChainId)
	if err != nil {
		return err
	}
	if fi.Variant != "dm" {
		return ErrInvalidVariant
	}
	sshClient, _, ch, err := client.CreateSSHSession(ctx, fi.UserA.Name, message.Address, fi.UserA.PrivateKey)
	if err != nil {
		return err
	}
	defer ch.Execute()

	adminClient, chRpc, err := client.AttachRPCAdminClient(ctx, sshClient, false)
	if err != nil {
		return err
	}
	defer chRpc.Execute()

	found := false
	err = adminClient.DumpDirectMessages(fi.UserA.Name, func(entry *gsDatabase.MessageEntry) {
		found = found || strings.Contains(entry.Body, message.Flag)
	})
	if err != nil {
		return err
	}

	if !found {
		return checker.ErrFlagNotFound
	}

	return nil
}

func (h *Handler) getFlagPrivateRoom(ctx context.Context, message *checker.TaskMessage) error {
	fi, err := h.db.GetEntry(message.TaskChainId)
	if err != nil {
		return err
	}
	if fi.Variant != "room" {
		return ErrInvalidVariant
	}
	_, sessionIO, ch, err := client.CreateSSHSession(ctx, fi.UserA.Name, message.Address, fi.UserA.PrivateKey)
	if err != nil {
		return err
	}
	defer ch.Execute()

	joinCmd := fmt.Sprintf("/join %s %s", fi.Room, fi.Password)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, joinCmd, message.Flag)
	if err != nil {
		h.log.Error(err)
		return checker.ErrFlagNotFound
	}

	return nil
}

func (h *Handler) GetFlag(ctx context.Context, message *checker.TaskMessage) error {
	if message.VariantId >= serviceInfo.FlagVariants {
		return ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return h.getFlagDirectMessage(ctx, message)
	case 1:
		return h.getFlagPrivateRoom(ctx, message)
	}

	return ErrVariantNotFound
}

func (h *Handler) putNoiseDirectMessage(ctx context.Context, message *checker.TaskMessage) error {
	userA, err := client.GenerateNewUser()
	if err != nil {
		return err
	}
	userB, err := client.GenerateNewUser()
	if err != nil {
		return err
	}

	noise := client.GenerateNoise()

	err = h.sendDirectMessage(ctx, userA, userB, message.Address, noise)
	if err != nil {
		return err
	}

	err = h.db.PutEntry(&database.Entry{
		Type:        "noise",
		Variant:     "dm",
		TaskMessage: message,
		UserA:       userA,
		UserB:       userB,
		Noise:       noise,
	})

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) getNoiseDirectMessage(ctx context.Context, message *checker.TaskMessage) error {
	fi, err := h.db.GetEntry(message.TaskChainId)
	if err != nil {
		return err
	}
	if fi.Variant != "dm" {
		return ErrInvalidVariant
	}
	_, sessionIO, ch, err := client.CreateSSHSession(ctx, fi.UserA.Name, message.Address, fi.UserA.PrivateKey)
	if err != nil {
		return err
	}
	defer ch.Execute()
	historyCmd := fmt.Sprintf("/history %s", fi.UserB.Name)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, historyCmd, fi.Noise)
	if err != nil {
		h.log.Error(err)
		return checker.ErrNoiseNotFound
	}

	return nil
}

func (h *Handler) PutNoise(ctx context.Context, message *checker.TaskMessage) error {
	if message.VariantId >= serviceInfo.NoiseVariants {
		return ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return h.putNoiseDirectMessage(ctx, message)
	}

	return ErrVariantNotFound
}

func (h *Handler) GetNoise(ctx context.Context, message *checker.TaskMessage) error {
	if message.VariantId >= serviceInfo.NoiseVariants {
		return ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return h.getNoiseDirectMessage(ctx, message)
	}

	return ErrVariantNotFound
}

func (h *Handler) havocRPC(ctx context.Context, message *checker.TaskMessage) error {
	userA, err := client.GenerateNewUser()
	if err != nil {
		return err
	}
	signer, err := ssh.NewSignerFromSigner(userA.PrivateKey)
	if err != nil {
		return err
	}
	sshClient, err := client.GetSSHClient(ctx, "quote-bot", message.Address, signer)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	adminClient, chRpc, err := client.AttachRPCAdminClient(ctx, sshClient, false)
	if err != nil {
		return err
	}
	defer chRpc.Execute()

	quote := quotes.GetRandom()
	msg := fmt.Sprintf("[quote-bot]: \"%s\" - %s", quote.Text, quote.From)
	err = adminClient.SendMessageToRoom("default", msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Havoc(ctx context.Context, message *checker.TaskMessage) error {
	if message.VariantId >= serviceInfo.HavocVariants {
		return ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return h.havocRPC(ctx, message)
	}

	return ErrVariantNotFound
}

func (h *Handler) Exploit(ctx context.Context, message *checker.TaskMessage) (*checker.HandlerInfo, error) {
	if message.VariantId >= serviceInfo.ExploitVariants {
		return nil, ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		flag, err := h.hijackAdminSession(ctx, message.AttackInfo, message.Address, message.FlagRegex)
		if err != nil {
			return nil, err
		}
		return checker.NewExploitInfo(flag), nil
	case 1:
		flag, err := h.hijackPrivateRoom(ctx, message.AttackInfo, message.Address, message.FlagRegex)
		if err != nil {
			return nil, err
		}
		return checker.NewExploitInfo(flag), nil
	}
	return nil, ErrVariantNotFound
}

func (h *Handler) GetServiceInfo() *checker.InfoMessage {
	return serviceInfo
}
