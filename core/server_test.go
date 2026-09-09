package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/coder/websocket"
)

func TestIsBenignClose(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"normal closure", websocket.CloseError{Code: websocket.StatusNormalClosure}, true},
		{"going away", websocket.CloseError{Code: websocket.StatusGoingAway}, true},
		{"no status rcvd", websocket.CloseError{Code: websocket.StatusNoStatusRcvd}, true},
		{"abnormal closure", websocket.CloseError{Code: websocket.StatusAbnormalClosure}, true},
		// The shape readPump actually sees: coder/websocket wraps the close
		// error with "failed to get reader" on its way out of Read.
		{"wrapped no status rcvd", fmt.Errorf("failed to get reader: received close frame: %w",
			websocket.CloseError{Code: websocket.StatusNoStatusRcvd}), true},
		{"conn already closed", net.ErrClosed, true},
		{"eof", io.EOF, true},
		{"context canceled", context.Canceled, true},
		{"wrapped context canceled", fmt.Errorf("read: %w", context.Canceled), true},

		{"policy violation", websocket.CloseError{Code: websocket.StatusPolicyViolation}, false},
		{"protocol error", websocket.CloseError{Code: websocket.StatusProtocolError}, false},
		{"message too big", websocket.CloseError{Code: websocket.StatusMessageTooBig}, false},
		{"rate limited", ErrRateLimited, false},
		{"plain error", errors.New("boom"), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isBenignClose(c.err); got != c.want {
				t.Errorf("isBenignClose(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}
