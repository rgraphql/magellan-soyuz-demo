package main

import (
	"context"
	"io"
	"strings"
)

// isNormalCloseError checks if the error is normal and should be ignored.
func isNormalCloseError(err error) bool {
	if err == nil {
		return true
	}
	errStr := err.Error()
	return err == context.Canceled ||
		err == io.EOF ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "ws closed: 1001") ||
		strings.Contains(errStr, "ws closed: 1005") ||
		strings.Contains(errStr, "websocket: close 1001 (going away)")
}

// isNormalHandshakeError checks if the error is normal and should be ignored.
func isNormalHandshakeError(err error) bool {
	if err == nil {
		return true
	}
	if isNormalCloseError(err) {
		return true
	}

	// ignore all handshake errors with bad headers, protocols, or request methods.
	return strings.HasPrefix(err.Error(), "handshake error: bad ")
}
