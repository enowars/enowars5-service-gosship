package checker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCheckerHandler struct {
}

func (m *mockCheckerHandler) PutFlag(ctx context.Context, message *TaskMessage) (*HandlerInfo, error) {
	return NewPutFlagInfo("user-123"), nil
}

func (m *mockCheckerHandler) GetFlag(ctx context.Context, message *TaskMessage) error {
	return ErrFlagNotFound
}

func (m *mockCheckerHandler) PutNoise(ctx context.Context, message *TaskMessage) error {
	return nil
}

func (m *mockCheckerHandler) GetNoise(ctx context.Context, message *TaskMessage) error {
	return nil
}

func (m *mockCheckerHandler) Havoc(ctx context.Context, message *TaskMessage) error {
	panic("implement me")
}

func (m *mockCheckerHandler) Exploit(ctx context.Context, message *TaskMessage) (*HandlerInfo, error) {
	return NewExploitInfo("FLAG123"), nil
}

func (m *mockCheckerHandler) GetServiceInfo() *InfoMessage {
	return &InfoMessage{
		ServiceName:     "mock",
		FlagVariants:    1,
		NoiseVariants:   1,
		HavocVariants:   0,
		ExploitVariants: 1,
	}
}

type nullLogger struct{}

func (m *nullLogger) Debugf(format string, args ...interface{}) {}
func (m *nullLogger) Infof(format string, args ...interface{})  {}
func (m *nullLogger) Warnf(format string, args ...interface{})  {}
func (m *nullLogger) Errorf(format string, args ...interface{}) {}
func (m *nullLogger) Debug(args ...interface{})                 {}
func (m *nullLogger) Info(args ...interface{})                  {}
func (m *nullLogger) Warn(args ...interface{})                  {}
func (m *nullLogger) Error(args ...interface{})                 {}

func newMockChecker() *Checker {
	return NewChecker(&nullLogger{}, &mockCheckerHandler{})
}

func TestGetServiceInfo(t *testing.T) {
	c := NewChecker(&nullLogger{}, &mockCheckerHandler{})
	req := httptest.NewRequest("GET", "http://checker/service", nil)
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	resp := w.Result()
	var info InfoMessage
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&info))
	assert.EqualValues(t, c.info, &info)
}

func newMockTaskMessage(method TaskMessageMethod, variantId uint64) *TaskMessage {
	return &TaskMessage{
		Method:    method,
		TeamName:  "test-team",
		VariantId: variantId,
		Timeout:   10000,
	}
}

func TestHandlers(t *testing.T) {
	c := newMockChecker()
	testCases := []struct {
		method    TaskMessageMethod
		variantId uint64
		result    Result
	}{
		{TaskMessageMethodPutFlag, 0, ResultOk},
		{TaskMessageMethodGetFlag, 0, ResultMumble},
		{TaskMessageMethodPutNoise, 0, ResultOk},
		{TaskMessageMethodGetNoise, 0, ResultOk},
		{TaskMessageMethodHavoc, 0, ResultError},
		{TaskMessageMethodExploit, 1, ResultError},
		{TaskMessageMethodExploit, 0, ResultOk},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s %d", testCase.method, testCase.variantId), func(t *testing.T) {
			taskMessage := newMockTaskMessage(testCase.method, testCase.variantId)
			payload, err := json.Marshal(taskMessage)
			require.NoError(t, err)
			req := httptest.NewRequest("POST", "http://checker/", bytes.NewReader(payload))
			w := httptest.NewRecorder()
			c.ServeHTTP(w, req)
			resp := w.Result()
			var resMessage ResultMessage
			err = json.NewDecoder(resp.Body).Decode(&resMessage)
			require.NoError(t, err)
			assert.EqualValues(t, testCase.result, resMessage.Result)
		})
	}
}
