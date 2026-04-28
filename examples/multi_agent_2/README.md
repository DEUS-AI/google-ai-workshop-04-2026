# Example 2 — Multi-Agent with Google Search

Same multi-agent structure as Example 1, but the specialist uses `geminitool.GoogleSearch` instead of a custom Go function.

A **specialist agent** (`city_analyst`) searches the web in real time to find a city's population. An **orchestrator** calls the specialist via `agenttool` and summarises the results.

## Key concepts

| Concept | Implementation |
|---|---|
| Built-in tool | `geminitool.GoogleSearch{}` — no Go function or structs needed |
| Specialist agent | `llmagent.New()` with the search tool attached |
| Orchestration | `agenttool.New()` wrapping the specialist for the orchestrator to call |

## Difference from Example 1

| | Example 1 | Example 2 |
|---|---|---|
| Tool type | Custom Go function (`functiontool`) | Native Gemini capability (`geminitool`) |
| Data source | Hardcoded local logic | Live web search |
| Structs required | Yes — input/output types | No |

## When to use this pattern

Use `geminitool.GoogleSearch` when your agent needs live, up-to-date information from the web that cannot be hardcoded or pre-fetched.

## Run

```bash
export GOOGLE_API_KEY="your-key-here"
go run multi_agent_2.go console
```

Then prompt: `What tier is Lisbon?`