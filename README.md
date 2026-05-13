<div align="center">

#  FreeLLM

**Free and easy access to large language models — for everyone.**

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)

---

*Providing access to AI by intelligently routing requests across free tiers of LLM providers.*

</div>

## ✨ What is FreeLLM?

FreeLLM is an open-source proxy that lets you use large language models **completely free** by leveraging the free tiers offered by various LLM providers. It tracks usage in real time and automatically routes requests to models that still have available quota — so you never hit a rate limit wall.

## 🔑 Key Features

- **Automatic model rotation** — Tracks token usage, requests per minute, and daily quotas. When one model's limits are reached, requests are seamlessly routed to the next available model.
- **OpenAI-compatible API** — FreeLLM exposes standard OpenAI-compatible endpoints. Any tool, library, or agent that speaks the OpenAI API can connect directly.
- **Coding agent ready** — Works out of the box with coding agents like [PI](https://github.com/pi). PI is the recommended agent, but any OpenAI-compatible client will work.
- **Multi-provider support** — Configure multiple LLM providers and models. FreeLLM maximizes your combined free-tier capacity.

## 📋 Available Models

| Provider | Model | Tokens/min | Requests/min | Requests/day |
|----------|-------|-----------|--------------|--------------|
| Gemini | gemini-2.5-flash | 250,000 | 5 | 20 |
| Gemini | gemini-3-flash | 250,000 | 5 | 20 |
| Gemini | gemini-3.1-flash-lite | 250,000 | 15 | 500 |
| Gemini | gemini-2.5-flash-lite | 250,000 | 10 | 20 |
| Gemini | gemma-4-26b | unlimited | 15 | 1,500 |
| Gemini | gemma-4-31b | unlimited | 15 | 1,500 |
| Groq | llama-3.1-8b-instant | 6,000 | 30 | 14,400 |
| Groq | llama-3.3-70b-versatile | 12,000 | 30 | 1,000 |

**Combined capacity:** 1,018,000 tokens/min • 125 requests/min • 18,960 requests/day


## 🐳 Quick Start with Docker

The fastest way to get FreeLLM running is with the pre-built Docker image from [Docker Hub](https://hub.docker.com/repository/docker/cezarychodun/freellm/general).

```bash
# Clone the repo and navigate to the example
git clone git@github.com:CezaryChodun/FreeLLM.git
cd FreeLLM/examples/docker-compose

# Configure your environment
cp .env.example .env
# Edit .env with your API keys and passwords

# Start all services
docker compose up -d
```

FreeLLM will be available at `http://localhost:3000`. Point any OpenAI-compatible client at this URL.

> **Note:** For non-local deployments, change the default database password in `docker-compose.yml` to a strong, unique password.

### What's included

The example spins up four containers:
- **freellm** — the proxy (port 3000)
- **litellm** — LLM gateway (port 4000)
- **postgres** — usage tracking database
- **prometheus** — metrics collection (port 9090)

### Customization

- **Choose models** — edit `config.yml`
- **Adjust rate limits** — edit `defaults/gemini.yml`
- **Add providers** — add new entries to `litellm-config.yaml`, `config.yml`, and corresponding defaults

### Environment Variables

All sensitive configuration is provided via `.env`. Copy `.env.example` and fill in your values:

| Variable | Description |
|----------|-------------|
| `LITELLM_MASTER_KEY` | Admin key for LiteLLM API access |
| `LITELLM_SALT_KEY` | Encryption key for stored API keys |
| `LITELLM_DB_PASSWORD` | Password for the LiteLLM PostgreSQL user |
| `GEMINI_API_KEY` | Your Google AI API key for Gemini/Gemma models |
| `GEMINI_API_BASE` | Google AI API base URL |
| `GROQ_API_KEY` | Your Groq API key for Llama/other Groq models |


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


## 🏗️ Architecture

```
┌─────────────┐     ┌──────────┐     ┌──────────┐     ┌──────────────┐
│  Your App   │────▶│ FreeLLM  │────▶│ LiteLLM  │────▶│ LLM Provider │
│  (PI, etc.) │     │  (proxy) │     │          │     │  (Gemini...) │
└─────────────┘     └──────────┘     └──────────┘     └──────────────┘
```

FreeLLM sits as a proxy layer on top of [LiteLLM](https://github.com/BerriAI/litellm). LiteLLM handles the connection management and provider-specific API translations. FreeLLM adds:

- **Usage tracking** — Persists per-model token and request counters in PostgreSQL
- **Rate limit awareness** — Compares real-time usage against configured quotas
- **Intelligent routing** — Model selection filtered to only models with remaining capacity


## 📄 License

This project is licensed under the GNU General Public License v3.0 — see the [LICENSE](LICENSE) file for details. [Learn more about GPL v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).

