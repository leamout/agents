# OpenAI Realtime Hello Agent

Minimal reference agent for Leamout's realtime execution path using OpenAI Realtime.

```text
Caller audio
    ↓
OpenAI Realtime
    ↕
Live speech-to-speech session
    ↓
Caller audio
```

The realtime provider replaces the composable STT → LLM → TTS chain for this session. Leamout still owns telephony, session lifecycle, tools, interruption handling, and durable conversation state around the provider session.

Credentials are tenant-owned and are not stored in this repository.
