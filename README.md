<div align="center">

# 🌐 FreeLLM

**Free, unlimited access to large language models — for everyone.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)

---

*Democratizing access to AI by intelligently routing requests across free tiers of LLM providers.*

</div>

## ✨ What is FreeLLM?

FreeLLM is an open-source proxy that lets you use large language models **completely free** by leveraging the free tiers offered by various LLM providers. It tracks usage in real time and automatically routes requests to models that still have available quota — so you never hit a rate limit wall.

## 🔑 Key Features

- **Automatic model rotation** — Tracks token usage, requests per minute, and daily quotas. When one model's limits are reached, requests are seamlessly routed to the next available model.
- **OpenAI-compatible API** — FreeLLM exposes standard OpenAI-compatible endpoints. Any tool, library, or agent that speaks the OpenAI API can connect directly.
- **Coding agent ready** — Works out of the box with coding agents like [PI](https://github.com/pi). PI is the recommended agent, but any OpenAI-compatible client will work.
- **Multi-provider support** — Configure multiple LLM providers and models. FreeLLM maximizes your combined free-tier capacity.

## 🚀 Quick Start

Connect any OpenAI-compatible client to FreeLLM:

```bash
export OPENAI_API_BASE=http://localhost:8080
```

That's it. Your requests will be automatically distributed across configured free-tier models.

## 🏗️ Architecture

```
┌─────────────┐     ┌──────────┐     ┌──────────┐     ┌──────────────┐
│  Your App   │────▶│ FreeLLM  │────▶│ LiteLLM  │────▶│ LLM Provider │
│  (PI, etc.) │     │  (proxy) │     │          │     │  (Gemini...) │
└─────────────┘     └──────────┘     └──────────┘     └──────────────┘
```

FreeLLM sits as a proxy layer on top of [LiteLLM](https://github.com/BerriAI/litellm). LiteLLM handles the connection management, API key rotation, and provider-specific API translations. FreeLLM adds:

- **Usage tracking** — Persists per-model token and request counters in PostgreSQL
- **Rate limit awareness** — Compares real-time usage against configured quotas
- **Intelligent routing** — Round-robin model selection filtered to only models with remaining capacity

## ⚙️ Configuration

Models are defined in `config.yml` with a `provider/model` format:

```yaml
models:
  - model: gemini/gemini-2.5-flash
  - model: gemini/gemma-3-27b-it
```

Rate limits for each provider are stored in the `defaults/` directory:

```yaml
# defaults/gemini.yml
- name: gemini-2.5-flash
  TPM: 250000
  RPM: 5
  RPD: 20
```

FreeLLM loads these at startup, populates the rate limits database, and begins routing immediately.

## 📦 Deployment

<!-- TODO -->

## 📄 License

MIT

