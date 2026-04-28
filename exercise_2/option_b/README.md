# Option B: The Ad Compliance Analyzer

Build an automated multi-agent pipeline that screens advertisement transcripts for legal and regulatory compliance using **Google ADK in Go**.

---

## The Goal

You are given **5 advertisement descriptions** (in [`ads.json`](ads.json)). Your system must analyze each ad against a compliance checklist and either clear it or flag it for review — autonomously, without a human reviewing each step.

---

## System Architecture

The structure mirrors what you built in Exercise 1. Replace the Flight Director with a **Final Reviewer**, and the specialists with **Compliance Agents**.

```
ads.json
    │
    ├──► Specialist Agent A  ──┐
    ├──► Specialist Agent B  ──┼──► Final Reviewer Agent ──► Output
    └──► Specialist Agent C  ──┘
```

Your pipeline must have **at least 3 agents**:

1. **Specialist Compliance Agents** — Each checks one or more rules from the checklist. They receive the ad description and return a `Compliant` or `Non-Compliant` verdict with supporting evidence.
2. **Final Reviewer Agent** — Orchestrates the specialists (via `agenttool`), collects their verdicts, and makes the final call.

**Every agent must use at least one tool** (use `functiontool` to define them, exactly as in Exercise 1).

### ADK Agent Type hint

This pipeline is a natural fit for **LLM Agents** wired together via `agenttool` — the same pattern as the Flight Director. You could also explore a `Sequential` workflow agent if you want the specialist calls to be fully deterministic.

---

## Output Logic

| Final Verdict | Action |
|---|---|
| **Compliant** | Print `"This ad is compliant"` and terminate. |
| **Non-Compliant** | Trigger a tool that **sends an email** to `eduardo.carvalho@deus.ai`, naming the specific ad and the rules it violated. |

---

## Ad Compliance Checklist

Your agents must check the following rules:

| Rule | Description |
|---|---|
| **Shows Minors** | The ad features children or individuals under the legal age of majority. |
| **Shows Alcohol** | The ad features alcoholic beverages, alcohol branding, or the act of consuming alcohol. |
| **Shows Gambling** | The ad features betting, casino environments, or games of chance involving monetary stakes. |
| **Shows Tobacco** | The ad features cigarettes, cigars, tobacco packaging, or the act of smoking. |
| **Shows Weapons** | The ad features firearms, ammunition, tactical gear, or bladed weaponry. |
| **Sexually Suggestive** | The ad uses provocative posing, revealing clothing, or sexualized imagery to market the product. |
| **Contains Medical Claims** | The ad makes health-related assertions or uses medical authority (doctors, stethoscopes, clinical settings) to imply a product is safe or beneficial. |

---

## Input Data

The 5 ads are in [`ads.json`](ads.json). Each entry contains:

| Field | Description |
|---|---|
| `ad_name` | The brand/campaign name |
| `VISUAL SWEEP & COMPOSITION` | Overall layout and visual elements |
| `OBJECTS & PROPS` | Items present in the ad |
| `TEXT & SYMBOLS` | All visible text, slogans, and copy |
| `PEOPLE` *(when present)* | Descriptions of people featured |
| `ENVIRONMENT & CONTEXT` | Setting and framing |

---

## Technical Reference

Use the same packages from Exercise 1:

```go
import (
    "google.golang.org/adk/agent/llmagent"    // LLM-powered agents
    "google.golang.org/adk/tool/functiontool" // wrap Go functions as tools
    "google.golang.org/adk/tool/agenttool"    // wrap agents as tools for the orchestrator
    "google.golang.org/adk/model/gemini"      // Gemini model
)

// Same model and API key as Exercise 1
model, _ := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})
```

---

## Success Criteria

- [ ] Agents communicate with each other to reach a conclusion (inter-agent communication)
- [ ] At least **3 distinct agents** are implemented
- [ ] Every agent is equipped with **at least one functional tool**
- [ ] Non-compliant ads trigger an email notification to `eduardo.carvalho@deus.ai` with the ad name and violated rules
- [ ] Compliant ads print `"This ad is compliant"` and stop

---

## Tips

- **Same pattern, new domain.** The Final Reviewer calling specialists via `agenttool` is identical to the Flight Director calling the Meteorologist and Chief Engineer. Start from that skeleton.
- **Use `fmt.Printf()` inside your tools** — if you can't see the agents talking, you can't debug.
- **The email tool doesn't need to send a real email.** A `fmt.Printf("Sending email to: ...")` implementation is sufficient.
- **Group the checklist rules.** You don't need one agent per rule — a "Content Safety Agent" could cover minors, weapons, and tobacco; a "Regulatory Claims Agent" could handle medical claims and alcohol.
- **Think about what passes between agents.** A simple struct with `Verdict string` and `Evidence string` is enough for specialists to report back to the Final Reviewer.