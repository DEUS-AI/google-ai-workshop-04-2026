package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// =============================================================================
// SLIDE A — THE SPECIALIST: GIVING AN AGENT A CAPABILITY
// =============================================================================

//  1. Define the tool's input/output as Go structs.
//     ADK infers the JSON schema automatically from the json/jsonschema tags.
type cityArgs struct {
	City string `json:"city" jsonschema:"The city to check."`
}

type PopulationReport struct {
	Population int    `json:"population"`
	City       string `json:"city"`
	Tier       string `json:"tier"` // "large" or "small"
}

//  2. Implement the tool as a plain Go function.
//     No framework-specific logic inside — just regular Go.
func getPopulation(_ tool.Context, args cityArgs) (PopulationReport, error) {
	fmt.Printf("[Tool] GetPopulation called for: %s\n", args.City)

	if args.City == "Lisbon" {
		ret := PopulationReport{Population: 500000, City: args.City, Tier: "large"}
		fmt.Printf("[Tool] Returning: %v\n", ret)
		return ret, nil
	}

	ret := PopulationReport{Population: 50000, City: args.City, Tier: "small"}
	fmt.Printf("[Tool] Returning: %v\n", ret)
	return ret, nil
}

func main() {
	ctx := context.Background()

	// =============================================================================
	// SLIDE B — THE ORCHESTRATOR: AGENTS CALLING AGENTS
	// =============================================================================

	// 3. Initialize the model — shared across all agents.
	//    Individual agents can be swapped to a different model if needed.
	m, _ := gemini.NewModel(ctx, "gemini-3.1-flash-lite-preview", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})

	// 4. Register the Go function as a tool and attach it to an LLM agent.
	popTool, _ := functiontool.New(
		functiontool.Config{
			Name:        "get_population",
			Description: "Returns population data for a given city.",
		},
		getPopulation,
	)

	analyst, _ := llmagent.New(llmagent.Config{
		Name:        "city_analyst",
		Model:       m,
		Description: "Looks up city population and classifies it as large or small.",
		Instruction: `You are a city analyst. When given a city:
1. Call get_population.
2. Return "large city" or "small city" based on the tier field.`,
		Tools: []tool.Tool{popTool},
	})

	// 5. Wrap the specialist agent as a callable tool.
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

	// 6. Launch.
	l := full.NewLauncher()
	l.Execute(ctx, &launcher.Config{
		AgentLoader: agent.NewSingleLoader(orchestrator),
	}, os.Args[1:])
}
