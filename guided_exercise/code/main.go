package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// =====================================================================
// 1. STATE
// =====================================================================

// isEngineCalibrated tracks whether the engine has been calibrated.
// Starts as false so the first diagnostic always reports unaligned valves.
var isEngineCalibrated = false

// =====================================================================
// 2. INPUT / OUTPUT STRUCTS
// =====================================================================

// --- Meteorologist ---

type getWeatherArgs struct {
	Location string `json:"location" jsonschema:"The location to check weather for."`
}

type WeatherReport struct {
	WindSpeedKph  float64 `json:"wind_speed_kph"`
	LightningRisk string  `json:"lightning_risk"`
	Status        string  `json:"status"` // "GO" or "NO-GO"
}

// --- Chief Engineer ---

type targetVehicleArgs struct {
	TargetVehicle string `json:"target_vehicle" jsonschema:"The name of the vehicle."`
}

type DiagnosticsReport struct {
	EngineStatus string `json:"engine_status"`
	Status       string `json:"status"` // "GO" or "NO-GO"
}

type CalibrationReport struct {
	Message string `json:"message"`
}

// --- Flight Director ---

type launchRocketArgs struct {
	Vehicle     string `json:"vehicle"      jsonschema:"The rocket being launched."`
	MissionCode string `json:"mission_code" jsonschema:"The mission code."`
}

type LaunchReport struct {
	Message string `json:"message"`
}

// =====================================================================
// 3. TOOL IMPLEMENTATIONS
// =====================================================================

// getWeather is the Meteorologist's sole tool.
// It returns a WeatherReport with typed fields.
func getWeather(_ tool.Context, args getWeatherArgs) (WeatherReport, error) {
	fmt.Printf("\n[System] Executing Tool: GetWeather (Location: %s)\n", args.Location)

	switch strings.ToLower(args.Location) {
	case "cape canaveral":
		return WeatherReport{
			WindSpeedKph:  12.5,
			LightningRisk: "Low",
			Status:        "GO",
		}, nil
	case "storm base":
		return WeatherReport{
			WindSpeedKph:  85.0,
			LightningRisk: "High",
			Status:        "NO-GO",
		}, nil
	default:
		return WeatherReport{
			WindSpeedKph:  45.0,
			LightningRisk: "Medium",
			Status:        "NO-GO",
		}, nil
	}
}

// runDiagnostics checks the engine state via the isEngineCalibrated flag.
// Returns NO-GO with "Valves unaligned" until calibration has been performed.
func runDiagnostics(_ tool.Context, args targetVehicleArgs) (DiagnosticsReport, error) {
	fmt.Printf("\n[System] Executing Tool: RunDiagnostics (Vehicle: %s)\n", args.TargetVehicle)

	if !isEngineCalibrated {
		return DiagnosticsReport{
			EngineStatus: "Valves unaligned",
			Status:       "NO-GO",
		}, nil
	}
	return DiagnosticsReport{
		EngineStatus: "Nominal",
		Status:       "GO",
	}, nil
}

// calibrateEngine sets isEngineCalibrated to true so the next diagnostic passes.
func calibrateEngine(_ tool.Context, args targetVehicleArgs) (CalibrationReport, error) {
	fmt.Printf("\n[System] Executing Tool: CalibrateEngine (Vehicle: %s)\n", args.TargetVehicle)

	isEngineCalibrated = true
	return CalibrationReport{
		Message: "Calibration complete. Ready for re-test.",
	}, nil
}

// launchRocket is the Flight Director's sole tool, called only when all
// departments report GO.
func launchRocket(_ tool.Context, args launchRocketArgs) (LaunchReport, error) {
	fmt.Printf("\n[System] Executing Tool: LaunchRocket (Vehicle: %s, Mission: %s)\n",
		args.Vehicle, args.MissionCode)
	fmt.Println("🚀")
	return LaunchReport{
		Message: "Liftoff successful. Vehicle has cleared the tower.",
	}, nil
}

// =====================================================================
// 4. AGENT CONSTRUCTORS
// =====================================================================

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
		Description: "Specialist in launch-site meteorology. " +
			"Returns a GO or NO-GO recommendation based on wind speed and lightning risk.",
		Instruction: `You are a strict meteorologist. You analyze weather data and determine 
if conditions are safe for a rocket launch.

When asked about a location:
1. Call get_weather with that location.
2. Return a clear GO or NO-GO verdict, quoting the wind_speed_kph and lightning_risk values.
3. Be concise and factual — no speculation beyond the data.`,
		Tools: []tool.Tool{weatherTool},
	})
}

func newChiefEngineer(model model.LLM) (agent.Agent, error) {
	diagTool, err := functiontool.New(
		functiontool.Config{
			Name:        "run_diagnostics",
			Description: "Runs engine diagnostics for a vehicle. Returns engine_status and a GO/NO-GO status.",
		},
		runDiagnostics,
	)
	if err != nil {
		return nil, fmt.Errorf("chief engineer: diagnostics tool: %w", err)
	}

	calibTool, err := functiontool.New(
		functiontool.Config{
			Name:        "calibrate_engine",
			Description: "Fixes engine valve alignment. Use this when run_diagnostics reports unaligned valves, then re-run diagnostics to confirm.",
		},
		calibrateEngine,
	)
	if err != nil {
		return nil, fmt.Errorf("chief engineer: calibration tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:  "chief_engineer",
		Model: model,
		Description: "Specialist in rocket engine diagnostics and calibration. " +
			"Owns the full repair loop and reports a final GO/NO-GO to the Flight Director.",
		Instruction: `You are the chief engineer. You run diagnostics on the rocket. 
If issues are found, you must fix them using your available tools when instructed.

Your procedure for every readiness check:
1. Call run_diagnostics for the given vehicle.
2. If the result is NO-GO due to "Valves unaligned":
   a. Call calibrate_engine to fix the alignment.
   b. Call run_diagnostics again to confirm the fix.
3. Report the final GO or NO-GO status with a brief summary of the engine_status field.

Never report a final NO-GO before attempting at least one calibration cycle.`,
		Tools: []tool.Tool{diagTool, calibTool},
	})
}

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
		Description: "Mission Control Orchestrator. Coordinates the meteorologist and chief engineer, " +
			"then issues the final launch or abort call.",
		Instruction: `You are the Flight Director. You have final authority over the launch sequence.

Your procedure:
1. Call the 'meteorologist' tool with the launch site location to get the weather report.
2. Call the 'chief_engineer' tool with the vehicle name to get the readiness report.
   The Chief Engineer handles their own calibration loop — you only need their final verdict.
3. Evaluate both reports:
   - If weather is NO-GO → abort immediately and explain the weather hazard.
   - If the Chief Engineer's final report is NO-GO → abort and report the engineering issue.
   - If BOTH report GO → call 'launch_rocket' to complete the sequence.

Always summarise the status from each department before announcing your final decision.`,
		Tools: []tool.Tool{metTool, engTool, launchTool},
	})
}

// =====================================================================
// 5. MAIN
// =====================================================================

func main() {
	ctx := context.Background()

	// Shared Gemini model — swap individual agents to a different model as needed.
	model, err := gemini.NewModel(ctx, "gemini-3.1-flash-lite-preview", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	meteorologist, err := newMeteorologist(model)
	if err != nil {
		log.Fatalf("Failed to create meteorologist: %v", err)
	}

	chiefEngineer, err := newChiefEngineer(model)
	if err != nil {
		log.Fatalf("Failed to create chief engineer: %v", err)
	}

	flightDirector, err := newFlightDirector(model, meteorologist, chiefEngineer)
	if err != nil {
		log.Fatalf("Failed to create flight director: %v", err)
	}

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
