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

	"github.com/sirupsen/logrus"
)

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
			if strings.Contains(scanner.Text(), check) {
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

func (h *Handler) PutFlag(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
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

	flagMessage := fmt.Sprintf(":wave: %s", message.Flag)
	directMessage := fmt.Sprintf("/dm %s %s", userA.Name, flagMessage)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, directMessage, flagMessage)

	if err != nil {
		return nil, err
	}

	err = h.db.PutFlagInfo(&database.FlagInfo{
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

func (h *Handler) GetFlag(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	fi, err := h.db.GetFlagInfo(message.TaskChainId)
	if err != nil {
		return nil, err
	}
	_, sessionIO, ch, err := client.CreateSSHSession(ctx, fi.UserA.Name, message.Address, fi.UserA.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ch.Execute()
	historyCmd := fmt.Sprintf("/history %s", fi.UserB.Name)
	err = h.sendMessageAndCheckResponse(ctx, sessionIO, historyCmd, message.Flag)
	if err != nil {
		return nil, err
	}

	return &checker.ResultMessage{
		Result: checker.ResultOk,
	}, nil
}

func (h *Handler) PutNoise(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) GetNoise(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) Havoc(ctx context.Context, message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) GetServiceInfo() *checker.InfoMessage {
	return &checker.InfoMessage{
		ServiceName:   "gosship",
		FlagVariants:  1,
		NoiseVariants: 1,
		HavocVariants: 1,
	}
}
