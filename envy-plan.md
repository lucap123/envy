# envy — project plan

> Zero-config environment variable manager for developers.
> One tool. Every project. No cloud. No bullshit.

---

## Het probleem (waarom envy bestaat)

Elke developer kent dit:
- Je hebt 5 projecten, elk met een `.env` file
- Je copy-past tokens tussen projecten
- Je commit per ongeluk een secret (en dan begint de paniek)
- Je wisselt van project en vergeet welke vars er moeten zijn
- Je onboardt een nieuwe collega en stuurt de `.env` via Slack (slecht)

Bestaande tools (direnv, dotenv-vault, infisical) lossen dit op maar zijn:
- Te complex om op te zetten
- Cloud-afhankelijk
- Niet developer-first

envy is anders: **je typt één commando en het werkt gewoon.**

---

## Wat envy doet (core features)

### v0.1 — basis
- `envy init` — detecteert automatisch je project type (Node, Python, Go, etc.) en maakt een `.env.example` aan
- `envy set DATABASE_URL postgres://...` — sla een var op, encrypted lokaal
- `envy get DATABASE_URL` — haal een var op
- `envy list` — toon alle vars voor huidig project (waarden gemaskeerd)
- `envy run npm start` — run een commando met alle vars geïnjecteerd
- `envy export` — genereer een `.env` file van je opgeslagen vars

### v0.2 — git integratie
- `envy hook install` — installeert automatisch een git pre-commit hook
- Blokkeert een commit als er een secret in de staged files zit
- Detecteert patronen: API keys, tokens, passwords, private keys
- Toont welke file en welke regel het probleem is

### v0.3 — team sharing (geen cloud)
- `envy share --output team.env.enc` — exporteer encrypted bundle
- `envy import team.env.enc` — importeer met gedeeld wachtwoord
- Werkt via git repo, USB, email — jij kiest
- Geen server, geen account, geen SaaS

### v0.4 — profiles
- `envy profile add staging` — maak omgevingsprofiel aan
- `envy use staging` — switch naar staging vars
- `envy run --profile production npm start`

---

## Hoe het verschilt van alternatieven

| Feature | envy | direnv | dotenv-vault | infisical |
|---|---|---|---|---|
| Zero config | ✅ | ❌ | ❌ | ❌ |
| Git hook integratie | ✅ | ❌ | ❌ | ✅ maar complex |
| Geen cloud nodig | ✅ | ✅ | ❌ | ❌ |
| Encrypted lokaal | ✅ | ❌ | ✅ | ✅ |
| Team sharing offline | ✅ | ❌ | ❌ | ❌ |
| Install in 10 sec | ✅ | ❌ | ❌ | ❌ |

---

## Tech stack

- **Taal:** Go (zelfde als Scope — je kent het al)
- **Encryptie:** AES-256-GCM voor opgeslagen vars
- **Opslag:** `~/.envy/` — plain JSON encrypted, geen database
- **Git hooks:** gewone shell scripts die envy aanroepen
- **Auto-detectie:** herkent `package.json`, `go.mod`, `requirements.txt`, `Cargo.toml`, etc.

---

## Projectstructuur

```
envy/
├── cmd/
│   ├── init.go
│   ├── set.go
│   ├── get.go
│   ├── list.go
│   ├── run.go
│   ├── export.go
│   ├── hook.go
│   └── share.go
├── pkg/
│   ├── store/       # encrypted opslag
│   ├── detect/      # project type detectie
│   ├── secrets/     # secret scanning patronen
│   └── crypto/      # encryptie helpers
├── main.go
├── Makefile
└── README.md
```

---

## Roadmap

| Versie | Features | Tijdsinschatting |
|---|---|---|
| v0.1 | init, set, get, list, run, export | 1-2 weken |
| v0.2 | git hook, secret scanning | 1 week |
| v0.3 | team sharing encrypted | 1 week |
| v0.4 | profiles / omgevingen | 1 week |

---

## Promotie strategie (leer van Scope)

### Bij launch van v0.1
- **Hacker News Show HN** — geen karma nodig, perfecte doelgroep
- **README first** — schrijf de README voor je een regel code schrijft, zodat het verhaal klopt
- **Demo gif** — zelfde als Scope, toont in 10 seconden de waarde

### Bij launch van v0.2 (git hooks)
- **Reddit r/programming, r/golang, r/devops** — git hook feature is universeel, geen niche
- **Dev.to artikel** — "I built a tool that prevents secret leaks in 1 command" — schrijft zichzelf
- Tagline voor dit moment: *"The git commit hook that saves your job"*

### Verschil met Scope promotie
- Scope = niche (bug bounty hunters)
- envy = universeel (elke developer)
- Je kan gewoon posten zonder "promoting" gevoel want het lost een probleem op dat iedereen kent

---

## De killer tagline

```
envy — set it once. run everything.
```

Of voor de git hook angle:

```
envy — the last time you accidentally commit a secret.
```

---

## Eerste stap

Schrijf de README. Dan `envy init` en `envy run`. Die twee commando's zijn genoeg voor een eerste Show HN post.
