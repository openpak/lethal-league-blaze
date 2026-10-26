# NEX registry spec — lethal-league-blaze

Proposed `switch-nexd` registry entry for Lethal League Blaze
(`01003AB00983C000`). This file is the deliverable for the serial NEX
registry integration — it specifies the entry; the edit to
`platforms/switch/games/switch-nex/cmd/switch-nexd/main.go` happens later,
not in this repo.

Wire facts sourced from `scratch/decompile-lethal-league-blaze/RECON.md` and
`scratch/decompile-lethal-league-blaze/findings.json` (see
[protocol.md](protocol.md)).

## Entry

- Slug: `lethal-league-blaze`
- Title ID: `01003AB00983C000`
- NEX version: 4.6.2 (+ NEX_MM); Pia 5.14.0; PIA_UNITY 0.13.0
- Requested port base: **23130** (next free after `full-metal-furies` 23120
  in `switch-nex/bin/title-ports.json`; auth = base, secure = base+1,
  healthz = base+2 per `switch-nex/titles/README.md`)
- Protocols: MatchmakeExtension + NATTraversal; NO DataStore, NO Ranking
- Auth shape: RendezVous Login (NNID + Siglo jobs); no BAAS/DAuth/AAuth

```go
{"lethal-league-blaze", 23130, "<game-server-id>", func(e *protocols.Env) {
	e.StationScheme = "prudps"
	e.Pia = protocols.LegacyPia() // Pia 5.14.0, pre-5.19 cutover
	e.Match = protocols.NewMatchmakingStore()
	e.FinalizeCreateSession = true
	e.JoinEchoSession = true
	e.JoinRecapParticipate = true
	// NO e.Ranking, NO e.DataStore: verified absent from the binary.
}}
```

## Outstanding keys (TODO-KEY)

1. Numeric NEX game-server ID (the `%08x` arg in
   `g%08x-%%.s.n.srv.nintendo.net`) + RendezVous port: code constants, need
   Ghidra on `exefs/main.dec` or a live capture.
2. NEX login params / auth-flow constants (username/password derivation for
   the RendezVous NNID/Siglo flow).
3. NEX 4.6.2 DDL method IDs for the MatchmakeExtension calls used
   (browse/random/create/join/leave).

## Notes for the integrator

- Straightforward NEX title: P2P rollback gameplay after matchmaking, no
  developer web API, no extra middleware servers.
- `e.Pia` choice follows the era rule used by the existing registry (e.g.
  Mario Tennis Aces comment in `cmd/switch-nexd/main.go`): Pia 5.14.0 is
  pre-5.19, so the legacy shape applies unless console measurement says
  otherwise.
- NAT check traffic (`nncs1/nncs2`) belongs to `nn-nncs`, not this entry.
