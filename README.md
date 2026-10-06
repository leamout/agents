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

Every `agent.json` uses the portable Leamout Agent Manifest format.

```text
agent.json
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

Provider credentials, tenant integration IDs, and deployment-specific secrets are intentionally excluded. A future importer can resolve the manifest's provider names against the tenant's configured integrations when the Agent is installed.

Template `tools` entries reference tool definitions from the sibling `tools.json` file by name.

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

Agent references in a routing manifest are portable aliases that are resolved to deployed Voice Agent resources during installation or deployment.

Routing ends once an Agent receives the call. Intent handling, tool execution, transfers, hangup, and other call-control behavior remain inside the Agent Runtime and Agent tools. `routing.json` therefore does not contain SIP destinations, workflow nodes, `on_transfer` transitions, provider credentials, or tenant integration IDs.

See `examples/scheduled-routing/` for the reference format.

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

## Templates

Business-oriented templates include sensible provider defaults while keeping all provider bindings and config values replaceable at deployment time.

```text
templates/
├── receptionist/
├── customer-support/
├── appointment-booking/
└── sales-qualifier/
```

Each template contains:

```text
agent.json
routing.json
README.md
tools.json
```

The `agent.json` defines the Agent behavior, voice, execution engine, provider bindings, provider-specific overrides, and tool references. `tools.json` defines the built-in or webhook tool contracts required by the template. `routing.json` defines the initial inbound-call assignment for the template.

The included template routing manifests intentionally start with a single `default_agent` and empty `routes`/`schedules`. Deployments can add schedule-based exceptions without changing the Agent's behavior or tool contracts.

Changing providers does not require rewriting the Agent's instructions or tool contracts. For example, a template can replace its composable bindings with a realtime binding while preserving the business behavior.

## Repository boundary

This repository should not contain:

- provider SDK implementations,
- telephony infrastructure,
- media runtime code,
- tenant credentials,
- control-plane persistence,
- vendor-specific orchestration logic.

Those concerns belong in the Leamout runtime or provider repositories.
