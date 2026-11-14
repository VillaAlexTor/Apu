package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestReadinessWithoutDatabaseConfiguration(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()

	readinessHandler(config{}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, recorder.Code)
	}
}

func TestReadinessWithReachableDatabaseAddress(t *testing.T) {
	originalDial := databaseDial
	databaseDial = func(network, address string, timeout time.Duration) (io.Closer, error) {
		if network != "tcp" {
			t.Fatalf("expected tcp network, got %q", network)
		}
		if address != "127.0.0.1:5433" {
			t.Fatalf("expected database address 127.0.0.1:5433, got %q", address)
		}
		return io.NopCloser(strings.NewReader("")), nil
	}
	defer func() { databaseDial = originalDial }()

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()
	databaseURL := "postgres://apu:local@127.0.0.1:5433/apu?sslmode=disable"

	readinessHandler(config{databaseURL: databaseURL}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestDatabaseReachableRejectsMissingHostname(t *testing.T) {
	if err := databaseReachable("postgres:///apu"); err == nil {
		t.Fatal("expected missing hostname to fail")
	}
}
