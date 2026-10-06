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

The examples stay intentionally small and exercise the engine boundary rather than business logic.

```text
examples/
├── composable/
│   ├── README.md
│   └── agent.json
├── openai-realtime/
│   ├── README.md
│   └── agent.json
└── gemini-realtime/
    ├── README.md
    └── agent.json
```

- `composable` demonstrates STT + LLM + TTS composition.
- `openai-realtime` demonstrates OpenAI Realtime as the complete live AI path.
- `gemini-realtime` demonstrates Gemini Live as the complete live AI path.

The manifests are reference agent definitions. Provider credentials remain tenant-owned and are not stored in this repository.

## Business templates

Templates focus on behavior, domain workflows, and tool contracts rather than a specific AI vendor.

```text
templates/
├── receptionist/
│   ├── README.md
│   ├── agent.json
│   └── tools.json
├── customer-support/
│   ├── README.md
│   ├── agent.json
│   └── tools.json
├── appointment-booking/
│   ├── README.md
│   ├── agent.json
│   └── tools.json
└── sales-qualifier/
    ├── README.md
    ├── agent.json
    └── tools.json
```

Template `agent.json` files are deliberately engine-neutral partial definitions. At deployment time, choose either the composable or realtime execution path and bind the required providers in Leamout. Template `tools.json` files use the Leamout tool shape and keep deployment-specific webhook URLs as placeholders.

This separation lets the same receptionist or support workflow run with different provider combinations without rewriting business behavior.

## Repository boundary

This repository should not contain:

- provider SDK implementations,
- telephony infrastructure,
- media runtime code,
- tenant credentials,
- control-plane persistence,
- vendor-specific orchestration logic.

Those concerns belong in the Leamout runtime or provider repositories.
