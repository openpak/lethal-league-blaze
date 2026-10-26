# Lethal League Blaze — the client's contract

Read from the Switch executable (Lethal League Blaze `01003AB00983C000`, base
v0 NSP, Unity 2017.4.17f1 IL2CPP). Every wire fact below comes from the
workspace recon at `scratch/decompile-lethal-league-blaze/RECON.md` and
`scratch/decompile-lethal-league-blaze/findings.json`; keep this file in step
with `cmd/lethal-league-blazed`.

## Stack

- NEX 4.6.2 + NEX_MM (`SDK MW+Nintendo+NEX-4_6_2-`,
  `NEX-for-NX-4_6_2-S640-20180920` include path; RECON.md "Container" note
  corrects the earlier gap line that claimed 4.6.3).
- Pia 5.14.0 full suite (PiaNex-forNEX-4_6_2, PiaSession, PiaTransport,
  PiaSync, PiaLan, PiaLocal, PiaCommon, PiaFramework, PiaClone) +
  PIA_UNITY-0.13.0 Unity bindings. 375 `nn::nex` + 359 `nn::pia` RTTI symbols,
  all in `main` (`exefs/main.dec`); `sdk.dec` has zero NEX refs.
- Transport: PRUDP v0/v1 + PRUDPLite (`PRUDPEndPoint`, `SecureEndPoint`,
  `PRUDPStream`, `NATTraversalEngine/Stream`, `NATRelayInterface`).
- StationURL forms: `prudp:/address=...;port=...`, `prudps:/`, `;stream=`,
  `;type=`.

## Hosts the client dials

- `g%08x-%%.s.n.srv.nintendo.net` — NEX game server. The numeric game-server
  ID (the `%08x` argument) is still unknown: a code constant, needs Ghidra on
  `exefs/main.dec` or a live capture (RECON.md open item 1). TODO-KEY.
- `%s.%%.p.srv.nintendo.net` — NEX portal host template.
- `nncs1/nncs2-%.n.n.srv.nintendo.net` — NAT check.
- `mii-secure-%%.cdn.nintendo.net` — Mii faces.

## Protocols and operations

NEX DDL protocols linked: MatchmakeExtension, NATTraversal,
NintendoAuthentication, AuthenticationFoundation, ProtocolFoundation,
NintendoProtocolFoundation.

- Login: `RendezVous::Login()`, `JobNNIDLogin`, `JobSigloLogin`,
  `LoginJob::{LoginFriendServer,PreLogin,PostLogin,...}`. NPDM services
  `baas`/`dauth`/`aauth` are absent — auth is RendezVous Login, not
  token-service mediated. Login params / auth-flow constants still outstanding
  (RECON.md open item 2). TODO-KEY.
- Matchmaking (`nn::pia::nex::NexMatch*`): create/join/browse/random-match/
  leave/destroy session, joint-session mesh, host migration, invitations +
  companion flow, session system passwords, community management incl.
  `FindOfficialCommunity`, NAT session start, monitoring-data send. NEX 4.6.2
  DDL method IDs for the calls used (browse/random/create/join/leave) still
  outstanding (RECON.md open item 3). TODO-KEY.
- Friends: `friend:u` + `nn::friends::Friend::GetPresence` + friend-server
  login (invites/presence path).
- Local play: full `LanMatch*` mesh stack over `ldn:u` (host candidate
  ranking, session keys, mesh restart).
- Netcode: deterministic rollback (FixMath.NET + `MSG_HEADER_ONLINE` +
  `MSG_OOS` desync + playerjoin/addplayer/spectate). Gameplay after
  matchmaking is peer-to-peer; nothing here simulates the game.

## What is not here

Verified absent (count 0 in the recon): NEX DataStore + Ranking protocols,
Photon, PlayFab, Demonware, GameSparks, brainCloud, ChilliConnect, mbedTLS,
embedded certs (TLS via the system `ssl` service), and any game HTTPS/API
hosts. The 263 `UnityEngine.Networking` (UNet) refs are engine-bundled with
no evidence of gameplay use — NEX/Pia is the transport (RECON.md open item 4,
low priority).

## Build facts (for the registry entry)

- Title ID `01003AB00983C000`, DisplayVersion `1.14`, publisher Team Reptile
  (control NACP; internal project name "LL2" from leaked dev-tree paths).
- NCA build SDK 6.6.0 (in-binary libs 6.4.2), MasterKey rev 7, KeyGen 8.
- No registry or SERVER entry exists yet in `switch-nex/` or `npln/`
  (RECON.md "Registry / server status").
