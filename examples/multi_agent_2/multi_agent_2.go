package main

import (
	"context"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

// =============================================================================
// SLIDE A — THE SPECIALIST: GIVING AN AGENT A BUILT-IN CAPABILITY
// =============================================================================

// geminitool.GoogleSearch replaces getPopulation from multi_agent_1.
// Key difference: no Go function, no input/output structs.
// The tool is a native Gemini capability — the model invokes it internally
// to retrieve live web results. Just instantiate the struct and pass it as a tool.

func main() {
	ctx := context.Background()

	// =============================================================================
	// SLIDE B — THE ORCHESTRATOR: AGENTS CALLING AGENTS
	// =============================================================================

	// 1. Initialize the model — shared across all agents.
	//    GoogleSearch requires a Gemini 2+ model.
	m, _ := gemini.NewModel(ctx, "gemini-3.1-flash-lite-preview", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})

	// 2. Attach the built-in Google Search tool to an LLM agent.
	//    No functiontool.New() needed — geminitool.GoogleSearch{} is ready to use.
	analyst, _ := llmagent.New(llmagent.Config{
		Name:        "city_analyst",
		Model:       m,
		Description: "Searches the web for city population data and classifies it as large or small.",
		Instruction: `You are a city analyst. When given a city:
1. Use google_search to find its current population.
2. Return "large city" if population is over 300,000, or "small city" if not.`,
		Tools: []tool.Tool{geminitool.GoogleSearch{}},
	})

	// 3. Wrap the specialist agent as a callable tool.
	//    agenttool.New() is the core multi-agent primitive:
	//    the orchestrator retains control after each specialist call
	//    and acts on the result before deciding what to do next.
	analystTool := agenttool.New(analyst, nil)

	orchestrator, _ := llmagent.New(llmagent.Config{
		Name:        "orchestrator",
		Model:       m,
		Description: "Coordinates city analysis and summarises findings.",
		Instruction: `You are an orchestrator. Call city_analyst for each city
the user asks about, then provide a final summary comparing all results.`,
		Tools: []tool.Tool{analystTool},
	})

	// 4. Launch.
	l := full.NewLauncher()
	l.Execute(ctx, &launcher.Config{
		AgentLoader: agent.NewSingleLoader(orchestrator),
	}, os.Args[1:])
}
