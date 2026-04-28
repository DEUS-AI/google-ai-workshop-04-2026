# Agentic AI — What It Is and How to Build Your First Agentic System

**DEUS.ai // Workshop 2026**

A hands-on workshop that takes you from the theory of agentic AI to building a real multi-agent system using **Google ADK in Go**.

---

## Workshop Structure

The session is divided into a theory block followed by two progressive coding exercises.

### Theory
- Evolution of AI: from early ML to agentic systems
- Core concepts: LLMs, Transformer architecture
- Generative AI vs. Agentic AI: Prompt & Response vs. Reason & Act
- Key agentic concepts: Perceive → Reason → Plan → Act → Reflect
- Single-agent vs. multi-agent systems
- Google ADK: agent types, tool primitives, orchestration

### Exercises

| | Exercise | Description |
|---|---|---|
| 1 | [Operation Mars GO](exercise_1/README.md) | Guided. Build a multi-agent mission control system: a Flight Director orchestrating a Meteorologist and a Chief Engineer to decide whether to launch a rocket. |
| 2 | [The Director's Cut](exercise_2/README.md) | Free-form. Apply the same patterns independently to build your own agentic service. Choose a use case or follow Option B. |

---

## Exercise 2 — Options

**Option A — Build Your Own Use Case**
Use the ADK patterns from Exercise 1 to solve a problem of your choice. Examples from the session: Travel Concierge, Code Reviewer System, Smart Home Manager.

**Option B — The Ad Compliance Analyzer**
A structured challenge: build a multi-agent pipeline that screens advertisement transcripts for legal and regulatory compliance.
→ [Full brief and checklist](exercise_2/option_b/README.md)

---

## Examples

The [`examples/`](examples/) directory contains reference code for the core ADK patterns covered in the theory block.

| Example | Description |
|---|---|
| [`multi_agent_1/`](examples/multi_agent_1/multi_agent_example.go) | Base structure: specialist agent with a `functiontool` (plain Go function), orchestrator wiring via `agenttool`. Maps directly to the two code slides in the presentation. |
| [`multi_agent_2/`](examples/multi_agent_2/multi_agent_2.go) | Same structure as example 1, but the specialist uses `geminitool.GoogleSearch` instead of a custom function — no structs or Go logic required. |

---

## Prerequisites

- Go 1.26 or higher
- A Gemini API key from [Google AI Studio](https://aistudio.google.com) set as `GOOGLE_API_KEY`
- Basic familiarity with Go (functions, structs, interfaces)

---

## Getting Started

```bash
# Clone the repo
git clone <repo-url>
cd google-ai-workshop-04-2026

# Set your API key
export GOOGLE_API_KEY="your-key-here"

# Run Exercise 1
cd exercise_1/code
go run main.go console
```

---

## Key ADK Concepts Cheat Sheet

| Concept | Package | What it does |
|---|---|---|
| LLM Agent | `google.golang.org/adk/agent/llmagent` | Reasoning unit — one agent, one instruction, one set of tools |
| Function Tool | `google.golang.org/adk/tool/functiontool` | Wraps a plain Go function as an agent-callable tool |
| Agent Tool | `google.golang.org/adk/tool/agenttool` | Wraps a specialist agent so an orchestrator can call it |
| Gemini Model | `google.golang.org/adk/model/gemini` | Connects the agent to Gemini as its LLM backend |
| Gemini Tool | `google.golang.org/adk/tool/geminitool` | Built-in Gemini capabilities (e.g. `GoogleSearch{}`) — no Go function needed |