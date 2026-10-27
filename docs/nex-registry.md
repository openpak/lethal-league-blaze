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
- Requested port base: **23170** (23130 was taken by `tmnt-shredders-revenge`, 23160 requested by `rivals-of-aether`; verified against switch-nex `bin/title-ports.json` 2026-10-26; after `full-metal-furies` 23120
  in `switch-nex/bin/title-ports.json`; auth = base, secure = base+1,
  healthz = base+2 per `switch-nex/titles/README.md`)
- Protocols: MatchmakeExtension + NATTraversal; NO DataStore, NO Ranking
- Auth shape: RendezVous Login (NNID + Siglo jobs); no BAAS/DAuth/AAuth

```go
{"lethal-league-blaze", 23170, "a3267987", func(e *protocols.Env) {
	e.StationScheme = "prudps"
	e.Pia = protocols.LegacyPia() // Pia 5.14.0, pre-5.19 cutover
	e.Match = protocols.NewMatchmakingStore()
	e.FinalizeCreateSession = true
	e.JoinEchoSession = true
	e.JoinRecapParticipate = true
	// NO e.Ranking, NO e.DataStore: verified absent from the binary.
}}
```

## Key recovery (2026-10-27) — TODO-KEY closed

- Game server ID `21117E00`, access key `a3267987`: same-struct
  native construction in `PiaSettings::GetInternetInitSettings`
  (`0x10cfb80`) — ID into `w21` at `0x10cfc28/2c`, stored
  `[x19+0x1c]`; key string from literal slot `0x33461F0`
  (→ `0x35AFB70` = `"a3267987"`) at `[x19+0x20]`.
- Registry status: **REGISTERED** — entry + `bin/title-ports.json`
  + `switch-nex/docs/lethal-league-blaze.md`.

## Still outstanding

1. RendezVous port + NEX login params (NNID/Siglo derivation).
2. NEX 4.6.2 DDL method IDs for the MatchmakeExtension calls used
   (browse/random/create/join/leave) — first routed boot will show.

## Notes for the integrator

- Straightforward NEX title: P2P rollback gameplay after matchmaking, no
  developer web API, no extra middleware servers.
- `e.Pia` choice follows the era rule used by the existing registry (e.g.
  Mario Tennis Aces comment in `cmd/switch-nexd/main.go`): Pia 5.14.0 is
  pre-5.19, so the legacy shape applies unless console measurement says
  otherwise.
- NAT check traffic (`nncs1/nncs2`) belongs to `nn-nncs`, not this entry.
