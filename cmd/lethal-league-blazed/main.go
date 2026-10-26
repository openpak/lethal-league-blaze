// Command lethal-league-blazed is the OpenPak front for Lethal League Blaze's
// online stack on Switch (01003AB00983C000).
//
// The game is a straightforward NEX title (NEX 4.6.2 + NEX_MM + Pia 5.14.0,
// PRUDP v0/v1 + PRUDPLite, RendezVous Login, MatchmakeExtension + NATTraversal;
// no DataStore, no Ranking), so gameplay traffic is served by switch-nexd once
// the registry entry in docs/nex-registry.md lands. This daemon is the
// title's own HTTP surface: readiness reporting over the NEX auth/secure
// listeners via a small middleware client stub (nexclient.go), so deploys can
// gate on the title actually answering before the console is pointed at it.
//
// Endpoints (all plain HTTP on -addr):
//
//	GET /healthz  liveness: always 200 when the process is up
//	GET /readyz   readiness: 200 only when the configured NEX auth and secure
//	              listeners answer; otherwise 503 with the failing leg
//
// The wire facts behind the config defaults live in docs/protocol.md.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	addr       = flag.String("addr", ":23132", "HTTP listen address (/healthz, /readyz)")
	authAddr   = flag.String("nex-auth", "", "NEX auth (RendezVous) UDP address, host:port (env LLB_NEX_AUTH)")
	secureAddr = flag.String("nex-secure", "", "NEX secure UDP address, host:port (env LLB_NEX_SECURE)")
)

func main() {
	flag.Parse()
	cfg := configFromEnv(*authAddr, *secureAddr)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { io.WriteString(w, "ok\n") })
	mux.HandleFunc("/readyz", readyHandler(cfg))
	log.Printf("[lethal-league-blaze] serving HTTP on %s (auth=%s secure=%s)", *addr, cfg.authAddr, cfg.secureAddr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// config names the NEX listeners this title's switch-nexd instance serves.
// Defaults follow the requested port base 23130 (auth = base, secure =
// base+1, this daemon's HTTP = base+2); flags win over env, env over defaults.
type config struct {
	authAddr   string
	secureAddr string
}

func configFromEnv(authFlag, secureFlag string) config {
	return config{
		authAddr:   firstNonEmpty(authFlag, os.Getenv("LLB_NEX_AUTH"), "127.0.0.1:23130"),
		secureAddr: firstNonEmpty(secureFlag, os.Getenv("LLB_NEX_SECURE"), "127.0.0.1:23131"),
	}
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

func readyHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		c := newClient(cfg.authAddr, cfg.secureAddr)
		if err := c.ping(); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "not ready: %s\n", err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "ready\n")
	}
}

// portBase parses the requested NEX port base from the environment for
// operators who run several titles off one compose file. Unused by the
// daemon itself; kept so the documented base has one canonical reader.
func portBase() int {
	if v := os.Getenv("LLB_PORT_BASE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 23130
}
