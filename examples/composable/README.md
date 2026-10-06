# Composable Hello Agent

Minimal reference agent for Leamout's composable execution path.

```text
Caller audio
    ↓
Deepgram STT
    ↓
Groq LLM
    ↓
Cartesia TTS
    ↓
Caller audio
```

This example demonstrates the engine boundary only. The Agent Runtime owns call/session lifecycle, tool execution, conversation persistence, and media orchestration.

The provider bindings in `agent.json` identify roles and providers; credentials are supplied by the tenant through Leamout and are never committed with an agent definition.
