package handler

import (
	"bufio"
	"checker/pkg/checker"
	"checker/pkg/client"
	"checker/pkg/database"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/sirupsen/logrus"
)

var serviceInfo = &checker.InfoMessage{
	ServiceName:   "gosship",
	FlagVariants:  2,
	NoiseVariants: 1,
	HavocVariants: 1,
}

var ErrVariantIdOutOfRange = errors.New("variantId out of range")
var ErrVariantNotFound = errors.New("variant not found")
var ErrInvalidVariant = errors.New("invalid variant database entry")

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
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (h *Handler) putFlagDirectMessage(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	userA, err := client.GenerateNewUser()
	if err != nil {
		return nil, err
	}
	userB, err := client.GenerateNewUser()
	if err != nil {
		return nil, err
	}

	_, _, chA, err := client.CreateSSHSession(ctx, userA.Name, message.Address, userA.PrivateKey)
	if err != nil {
		return nil, err
	}
	chA.Execute()

	_, sessionIO, ch, err := client.CreateSSHSession(ctx, userB.Name, message.Address, userB.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ch.Execute()

	directMessage := fmt.Sprintf("/dm %s :wave: %s", userA.Name, message.Flag)
	checkStr := fmt.Sprintf("[dm][%s]:", userB.Name)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, directMessage, checkStr)

	if err != nil {
		return nil, err
	}

	err = h.db.PutFlagInfo(&database.FlagInfo{
		Variant:     "dm",
		TaskMessage: message,
		UserA:       userA,
		UserB:       userB,
	})

	if err != nil {
		return nil, err
	}

	return &checker.ResultMessage{
		Result: checker.ResultOk,
	}, nil
}

func (h *Handler) putFlagPrivateRoom(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	userA, err := client.GenerateNewUser()
	if err != nil {
		return nil, err
	}

	_, sessionIO, ch, err := client.CreateSSHSession(ctx, userA.Name, message.Address, userA.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ch.Execute()

	room, password := client.GenerateRoomAndPassword()
	createRoom := fmt.Sprintf("/create %s %s", room, password)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, createRoom, room+" was created")
	if err != nil {
		return nil, err
	}

	flagMessage := fmt.Sprintf(":clown: %s", message.Flag)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, flagMessage, message.Flag)
	if err != nil {
		return nil, err
	}

	err = h.sendMessageAndCheckResponse(ctx, sessionIO, "/j", "you are now in room default.")
	if err != nil {
		return nil, err
	}

	err = h.db.PutFlagInfo(&database.FlagInfo{
		Variant:     "room",
		TaskMessage: message,
		UserA:       userA,
		Room:        room,
		Password:    password,
	})

	if err != nil {
		return nil, err
	}

	return &checker.ResultMessage{
		Result: checker.ResultOk,
	}, nil
}

func (h *Handler) PutFlag(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
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

func (h *Handler) getFlagDirectMessage(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	fi, err := h.db.GetFlagInfo(message.TaskChainId)
	if err != nil {
		return nil, err
	}
	if fi.Variant != "dm" {
		return nil, ErrInvalidVariant
	}
	_, sessionIO, ch, err := client.CreateSSHSession(ctx, fi.UserA.Name, message.Address, fi.UserA.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ch.Execute()
	historyCmd := fmt.Sprintf("/history %s", fi.UserB.Name)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, historyCmd, message.Flag)
	if err != nil {
		h.log.Error(err)
		return &checker.ResultMessage{
			Result:  checker.ResultMumble,
			Message: "flag not found",
		}, err
	}

	return &checker.ResultMessage{
		Result: checker.ResultOk,
	}, nil
}

func (h *Handler) getFlagPrivateRoom(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	fi, err := h.db.GetFlagInfo(message.TaskChainId)
	if err != nil {
		return nil, err
	}
	if fi.Variant != "room" {
		return nil, ErrInvalidVariant
	}
	_, sessionIO, ch, err := client.CreateSSHSession(ctx, fi.UserA.Name, message.Address, fi.UserA.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ch.Execute()

	historyCmd := fmt.Sprintf("/join %s %s", fi.Room, fi.Password)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, historyCmd, message.Flag)
	if err != nil {
		h.log.Error(err)
		return &checker.ResultMessage{
			Result:  checker.ResultMumble,
			Message: "flag not found",
		}, err
	}

	return &checker.ResultMessage{
		Result: checker.ResultOk,
	}, nil
}

func (h *Handler) GetFlag(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	if message.VariantId >= serviceInfo.FlagVariants {
		return nil, ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return h.getFlagDirectMessage(ctx, message)
	case 1:
		return h.getFlagPrivateRoom(ctx, message)
	}

	return nil, ErrVariantNotFound
}

func (h *Handler) PutNoise(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	if message.VariantId >= serviceInfo.NoiseVariants {
		return nil, ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return nil, errors.New("not implemented")
	}

	return nil, ErrVariantNotFound
}

func (h *Handler) GetNoise(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	if message.VariantId >= serviceInfo.NoiseVariants {
		return nil, ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return nil, errors.New("not implemented")
	}

	return nil, ErrVariantNotFound
}

func (h *Handler) Havoc(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	if message.VariantId >= serviceInfo.HavocVariants {
		return nil, ErrVariantIdOutOfRange
	}
	switch message.VariantId {
	case 0:
		return nil, errors.New("not implemented")
	}

	return nil, ErrVariantNotFound
}

func (h *Handler) GetServiceInfo() *checker.InfoMessage {
	return serviceInfo
}
