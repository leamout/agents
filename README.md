# Leamout Agents

Reference autonomous voice agents, templates, and examples built on Leamout.

This repository contains things built **with** Leamout. It does not contain runtime infrastructure required **by** Leamout.

```text
leamout/contracts
    ↓
leamout/ai-providers
    ↓
leamout/leamout
    ↓
leamout/agents
```

The main Leamout repository owns the Agent Runtime, media runtime, telephony, tools, sessions, and persistence. `leamout/ai-providers` owns vendor integrations. This repository shows how those pieces are composed into usable voice agents.

## Agent Manifest v1

Every Agent definition uses the portable Leamout Agent Manifest format. Reference examples keep a single `agent.json`; business templates may contain multiple Agent manifests under `agents/`.

```text
Agent manifest
├── schema_version
├── name
├── engine
├── language
├── instructions
├── voice
├── providers
├── tools
├── interruption_policy
└── recording_policy
```

`schema_version` versions the portable manifest format. It is separate from a deployed Voice Agent's configuration revision.

`voice` is a first-class Agent property. Composable sessions pass it to the selected TTS provider, while realtime sessions pass it to the selected realtime provider.

Each provider binding may include a provider-specific `config` object. These values override the adapter's built-in defaults for that Agent, allowing model and tuning choices such as Groq temperature or the selected Deepgram Flux model without changing global provider configuration.

Provider credentials, tenant integration IDs, and deployment-specific secrets are intentionally excluded. Deployment resolves the manifest's provider names against the tenant's configured integrations.

Template Agent `tools` entries reference tool definitions from the package's `tools.json` by name.

## Routing Manifest v1

`routing.json` is an optional portable manifest for selecting which Agent receives an inbound call before an Agent session begins.

```text
routing.json
├── schema_version
├── name
├── description
├── default_agent
├── routes
└── schedules
```

Routing stays intentionally small. `default_agent` is the fallback, while each route may select an Agent for a named schedule. Schedules use explicit timezone/day/time fields rather than arbitrary condition strings.

```text
Incoming call
    │
    ├── matching route ──→ selected Agent
    │
    └── no match ────────→ default Agent
```

Inside a business template, routing Agent aliases are the filenames under `agents/` without the `.json` suffix. A deployer can therefore validate every routing reference locally before creating runtime resources.

Routing ends once an Agent receives the call. Intent handling, tool execution, transfers, hangup, and other call-control behavior remain inside the Agent Runtime and Agent tools. `routing.json` therefore does not contain SIP destinations, workflow nodes, `on_transfer` transitions, provider credentials, or tenant integration IDs.

See `examples/scheduled-routing/` for the reference routing format.

## Execution models

Leamout supports two alternative AI execution paths.

```text
                         Agent Runtime
                              │
                ┌─────────────┴─────────────┐
                │                           │
                ▼                           ▼
        Composable Engine            Realtime Engine
                │                           │
       ┌────────┼────────┐                  │
       ▼        ▼        ▼                  ▼
      STT      LLM      TTS              Realtime

   Deepgram    Groq    Cartesia            OpenAI
   AssemblyAI  OpenAI  ElevenLabs          Gemini
```

A composable agent selects STT, LLM, and TTS independently. A realtime agent selects one realtime speech-to-speech provider.

## Reference examples

The reference examples use concrete provider bindings so developers can see a complete Agent manifest.

```text
examples/
├── composable/
│   ├── README.md
│   └── agent.json
├── openai-realtime/
│   ├── README.md
│   └── agent.json
├── gemini-realtime/
│   ├── README.md
│   └── agent.json
└── scheduled-routing/
    ├── README.md
    └── routing.json
```

- `composable` demonstrates Deepgram + Groq + Cartesia.
- `openai-realtime` demonstrates OpenAI Realtime as the complete live AI path.
- `gemini-realtime` demonstrates Gemini Live as the complete live AI path.
- `scheduled-routing` demonstrates schedule-based inbound-call Agent selection.

## Deployable templates

Business templates are self-contained portable packages. Every Agent alias referenced by `routing.json` must resolve to an Agent manifest shipped in the same template.

```text
templates/<template>/
├── agents/
│   ├── <primary>.json
│   └── after-hours.json
├── routing.json
├── tools.json
└── README.md
```

The package owns Agent behavior, provider defaults, tool contracts, routing, and schedules. Deployment supplies environment-specific values such as provider credentials/integration IDs, webhook endpoints, transfer destinations, routing timezone overrides, and phone-number bindings.

A deployer should be able to validate the package without external Agent definitions, then create all Agent resources, tools, and routing from the package in one deployment operation.

Changing providers does not require rewriting the Agent's instructions or tool contracts. For example, a package can replace its composable provider bindings with a realtime binding while preserving the business behavior.

## Repository boundary

This repository should not contain:

- provider SDK implementations,
- telephony infrastructure,
- media runtime code,
- tenant credentials,
- control-plane persistence,
- vendor-specific orchestration logic.

Those concerns belong in the Leamout runtime or provider repositories.
