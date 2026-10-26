# Lethal League Blaze

OpenPak's Lethal League Blaze Server.

| | |
| --- | --- |
| Game | Lethal League Blaze (Team Reptile) |
| Network | Nintendo Switch |
| Title ID | `01003AB00983C000` |
| Licence | AGPL-3.0-only |

## What it serves

| Host | What | Served by |
| --- | --- | --- |
| `g%08x-%%.s.n.srv.nintendo.net` | NEX game server: RendezVous login, matchmaking, NAT traversal, P2P mesh | [`switch-nex`](https://github.com/openpak/switch-nex) (registry entry proposed in [docs/nex-registry.md](docs/nex-registry.md), not yet landed) |
| `nncs1/nncs2-%.n.n.srv.nintendo.net` | NAT check | [`nn-nncs`](https://github.com/openpak/nn-nncs) |
| this repo (`lethal-league-blazed`) | `/healthz` liveness + `/readyz` readiness over the title's NEX listeners | this repo |

The game is a straightforward NEX title: NEX 4.6.2 + NEX_MM + Pia 5.14.0 over
PRUDP v0/v1 + PRUDPLite. Gameplay after matchmaking is peer-to-peer
(FixMath.NET deterministic rollback). No DataStore, no Ranking, no
developer web API.

## Run

```sh
go run ./cmd/lethal-league-blazed
```

| Setting | Default | Meaning |
| --- | --- | --- |
| `-addr` | `:23132` | HTTP listen address (`/healthz`, `/readyz`) |
| `-nex-auth` / `LLB_NEX_AUTH` | `127.0.0.1:23130` | NEX auth (RendezVous) UDP listener to probe |
| `-nex-secure` / `LLB_NEX_SECURE` | `127.0.0.1:23131` | NEX secure UDP listener to probe |
| `LLB_PORT_BASE` | `23130` | Requested NEX port base for this title |

## Docs

- [Protocol](docs/protocol.md)
- [NEX registry spec](docs/nex-registry.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
