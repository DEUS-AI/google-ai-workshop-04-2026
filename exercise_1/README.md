# Workshop Exercise: Operation Mars GO
## Building a Multi-Agent System with Google ADK in Go

## Objective
In this guided exercise, you will learn how to use the Google Agent Development Kit (ADK) for Go to build a simple multi-agent orchestration system.

## Prerequisites
- Go 1.26 or higher installed on your machine.
- A valid Gemini API Key (`GEMINI_API_KEY` set in your environment variables).
- Basic familiarity with Go syntax (functions, structs and primitive types).

## The Scenario
You are building the AI mission control system for the upcoming Mars GO rocket launch. The final decision to launch requires gathering real-time telemetry from different departments. If an issue is detected, your agents must work together to attempt a fix before aborting the mission.

## System Architecture & Tool Specifications
You will implement three agents, each with specific tools. To make the system testable, you will write standard Go functions for these tools and add `fmt.Println()` statements inside them so we can see the agents "thinking" in the terminal.

---

## The Meteorologist (Agent Tool)

### Description:
>“Specialist in launch-site meteorology.”

### Instruction:
>“You are a strict meteorologist. You analyze weather data and determine if conditions are safe for a rocket launch.  
>When asked about a location:
>1. Call `get_weather` with that location.
>2. Return a clear GO or NO-GO verdict, quoting the `wind_speed_kph` and `lightning_risk` values.
>3. Be concise and factual — no speculation beyond the data.”


### Tool: GetWeather
**Input**
- Location (`string`)

**Output**
- WindSpeedKph (`float64`)
- LightningRisk (`string`)
- Status (`string`)

**Logic:**

| Location Input      | WindSpeedKph | LightningRisk | Status |
|---------------------|--------------|----------------|--------|
| "Cape Canaveral"    | 12.5         | "Low"          | "GO"   |
| "Storm Base"        | 85.0         | "High"         | "NO-GO"|

---

## The Chief Engineer (Agent Tool)

### Description:
>“Specialist in rocket engine diagnostics and calibration. Owns the full repair loop and reports a final GO/NO-GO to the Flight Director.”

### Instruction:
>“You are the chief engineer. You run diagnostics on the rocket.  
>If issues are found, you must fix them using your available tools when instructed.
>
>Your procedure for every readiness check:
>1. Call `run_diagnostics` for the given vehicle.
>2. If the result is NO-GO due to "Valves unaligned":
>   a. Call `calibrate_engine` to fix the alignment.  
>   b. Call `run_diagnostics` again to confirm the fix.
>3. Report the final GO or NO-GO status with a brief summary of the `engine_status` field.
>
>Never report a final NO-GO before attempting at least one calibration cycle.”

### Tool: RunDiagnostics

**State Variable:**  
In your Go code, create a boolean variable `isEngineCalibrated` and set it to `false` initially.

**Input:** 
- TargetVehicle (`string`)  
**Output:** 
- EngineStatus (`string`)
- Status (`string`) 

**Logic:**
- If `isEngineCalibrated` is false:  
  Return `EngineStatus: "Valves unaligned"`, `Status: "NO-GO"`.
- If `isEngineCalibrated` is true:  
  Return `EngineStatus: "Nominal"`, `Status: "GO"`.

### Tool: CalibrateEngine
**Input Struct:** 
- TargetVehicle (`string`)  
**Output Struct:** 
- Message (`string`)

**Logic:**  
Sets `isEngineCalibrated` to true and returns:  
"Calibration complete. Ready for re-test."

---

## The Flight Director (Orchestrator Agent)

### Description:
>“Mission Control Orchestrator. Coordinates the meteorologist and chief engineer, then issues the final launch or abort call.”

### Instruction:
>“You are the Flight Director. You have final authority over the launch sequence.
>
>Your procedure:
>1. Call the 'meteorologist' tool with the launch site location to get the weather report.
>2. Call the 'chief_engineer' tool with the vehicle name to get the readiness report.  
>   The Chief Engineer handles their own calibration loop — you only need their final verdict.
>3. Evaluate both reports:
>   - If weather is NO-GO → abort immediately and explain the weather hazard.
>   - If the Chief Engineer's final report is NO-GO → abort and report the engineering issue.
>   - If BOTH report GO → call 'launch_rocket' to complete the sequence.
>
>Always summarise the status from each department before announcing your final decision.”

### Tool: LaunchRocket  
**Input Struct:** 
- Vehicle (`string`)
- MissionCode (`string`)

**Logic:**  
Returns:
>"Liftoff successful."

---

## Step-by-Step Milestones

### Milestone 1: Project Setup & Tools
- Initialize a new Go module and import the Google ADK Go SDK.
```bash
go mod init mars-go
go get google.golang.org/adk@v1.0.0
go get google.golang.org/genai
```

Export you gemini token in as env var in your terminal:
```bash
export GOOGLE_API_KEY="YOUR_TOKEN"
```

Create a `main.go` file with a main function.

Instantiate the model LLM as:

```go
import (
  "google.golang.org/adk/model/gemini"
)

func main() {
  ctx := context.Background()
  model, err := gemini.NewModel(ctx, "gemini-3.1-flash-lite-preview", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
  })
  if err != nil {
    log.Fatalf("Failed to create model: %v", err)
  }
}
```

- Write the Go functions for GetWeather, RunDiagnostics, CalibrateEngine, and LaunchRocket.
- Remember to add `fmt.Printf("[System] Executing Tool: %s\n", toolName)` inside each function!

Here's the GetWeather example:
```go

type getWeatherArgs struct {
	Location string `json:"location" jsonschema:"The location to check weather for."`
}

type WeatherReport struct {
	WindSpeedKph  float64 `json:"wind_speed_kph"`
	LightningRisk string  `json:"lightning_risk"`
	Status        string  `json:"status"` // "GO" or "NO-GO"
}

func getWeather(_ tool.Context, args getWeatherArgs) (WeatherReport, error) {
  fmt.Printf("\n[System] Executing Tool: GetWeather (Location: %s)\n", args.Location)

  // add implementation
}
```

### Milestone 2: Agent Assembly

- Instantiate the Meteorologist and Chief Engineer agents, attaching their respective tools and system instructions using the ADK `functiontool` package.

Example:
```go
func newMeteorologist(model model.LLM) (agent.Agent, error) {
	weatherTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_weather",
			Description: "Retrieves current weather conditions (wind speed, lightning risk, GO/NO-GO status) for a given launch-site location.",
		},
		getWeather,
	)
	if err != nil {
		return nil, fmt.Errorf("meteorologist: weather tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:  "meteorologist",
		Model: model,
		Description: "...",
		Instruction: `...`,
		Tools: []tool.Tool{weatherTool},
	})
}
```

- Instantiate the Flight Director agent and link it to the agent tools using the `agenttool` package.

```go
func newFlightDirector(model model.LLM, meteorologist, chiefEngineer agent.Agent) (agent.Agent, error) {
	// Wrap sub-agents as callable tools so the Flight Director retains
	// control after each call and can act on the returned result.
	metTool := agenttool.New(meteorologist, nil)
	engTool := agenttool.New(chiefEngineer, nil)

	launchTool, err := functiontool.New(
		functiontool.Config{
			Name:        "launch_rocket",
			Description: "Executes the physical rocket launch sequence. Only call this when ALL departments have reported GO.",
		},
		launchRocket,
	)
	if err != nil {
		return nil, fmt.Errorf("flight director: launch tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:  "flight_director",
		Model: model,
		Description: "...",
		Instruction: `...`,
		Tools: []tool.Tool{metTool, engTool, launchTool},
	})
}
```

### Milestone 3: Launch
- Setup the adk launcher using the flight director as root agent

```go
func main() {
  ... // Instantiate agents

  config := &launcher.Config{
    AgentLoader: agent.NewSingleLoader(flightDirector),
  }

  l := full.NewLauncher()

  fmt.Println("--- STARTING MISSION CONTROL ---")
  fmt.Println("To test, type: 'Initiate launch sequence for the Ares-1 vehicle at Cape Canaveral.'")

  if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
    log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
  }
}
```
- Execute the program

```bash
go run agent.go web api webui
```
The args web api webui will launch a web UI interface that can be used optionally.
You can also chat directly in the command line.

---

## Testing & Expected Output

Once built, trigger the Flight Director with this prompt:  
"Initiate launch sequence for the Ares-1 vehicle at Cape Canaveral."

### What you should see in your terminal:
You will see the agents collaborating in real-time through your log statements:

1. `[System] Executing Tool: GetWeather (Cape Canaveral)`
2. `[System] Executing Tool: RunDiagnostics`
3. `[System] Executing Tool: CalibrateEngine`
4. `[System] Executing Tool: RunDiagnostics`
5. `[System] Executing Tool: LaunchRocket`