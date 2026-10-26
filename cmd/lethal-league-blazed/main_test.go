package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthz(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { io.WriteString(w, "ok\n") })
	srv := httptest.NewServer(mux)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || string(body) != "ok\n" {
		t.Fatalf("GET /healthz = %d %q", resp.StatusCode, body)
	}
}

func TestReadyz(t *testing.T) {
	// Two silent UDP listeners stand in for the NEX auth/secure servers.
	auth, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer auth.Close()
	secure, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer secure.Close()

	h := readyHandler(config{authAddr: auth.LocalAddr().String(), secureAddr: secure.LocalAddr().String()})
	req := httptest.NewRequest("GET", "/readyz", nil)

	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != 200 {
		t.Fatalf("ready legs: GET /readyz = %d, want 200 (%s)", rec.Code, rec.Body.String())
	}

	// Nothing listening on either leg: must be 503. Bind-then-close to get
	// ports that are (momentarily) free.
	free := func() string {
		l, err := net.ListenPacket("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		return l.LocalAddr().String()
	}
	h = readyHandler(config{authAddr: free(), secureAddr: free()})
	rec = httptest.NewRecorder()
	done := make(chan struct{})
	go func() { h(rec, req); close(done) }()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("GET /readyz hung")
	}
	if rec.Code != 503 {
		t.Fatalf("dead legs: GET /readyz = %d, want 503", rec.Code)
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("LLB_NEX_AUTH", "10.0.0.1:10000")
	t.Setenv("LLB_NEX_SECURE", "10.0.0.1:10001")
	cfg := configFromEnv("", "")
	if cfg.authAddr != "10.0.0.1:10000" || cfg.secureAddr != "10.0.0.1:10001" {
		t.Fatalf("env config = %+v", cfg)
	}
	cfg = configFromEnv("flag:1", "flag:2")
	if cfg.authAddr != "flag:1" || cfg.secureAddr != "flag:2" {
		t.Fatalf("flag config = %+v", cfg)
	}
}
