# Example 1 — Multi-Agent with Function Tool

Demonstrates the base structure of a multi-agent system in Google ADK (Go).

A **specialist agent** (`city_analyst`) uses a plain Go function wrapped as a `functiontool` to return hardcoded population data for a city. An **orchestrator** calls the specialist via `agenttool` and summarises the results.

## Key concepts

| Concept | Implementation |
|---|---|
| Custom tool | `functiontool.New()` wrapping a typed Go function |
| Specialist agent | `llmagent.New()` with one tool attached |
| Orchestration | `agenttool.New()` wrapping the specialist for the orchestrator to call |

## When to use this pattern

Use `functiontool` when your tool needs to run local Go logic — database lookups, API calls, calculations, or any deterministic operation you control.

## Run

```bash
export GOOGLE_API_KEY="your-key-here"
go run multi_agent_example.go console
```

Then prompt: `What tier is Lisbon?`