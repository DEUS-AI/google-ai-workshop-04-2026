# Challenge 2: The "Director's Cut" (Free-Form)

Now that you have a working mission control, it’s time to take off the training wheels. In this phase, you have the freedom to build something entirely new.

## Option A: Build Your Own Use Case
Use the ADK patterns you just learned to solve a different problem:
- Travel Concierge: An orchestrator that talks to a "Flight Agent," a "Hotel Agent," and a "Budget Agent."
- Code Reviewer System: An "Architect Agent" that reviews code structure and a "Security Agent" that hunts for vulnerabilities.
- Smart Home Manager: An agent system that coordinates lighting, temperature, and security based on user "mood" prompts.

## Option B: The Ad Compliance Analyzer
- Follow the instructions in the corresponding [README](option_b/README.md).


## Success Criteria
To complete this challenge successfully, your system must meet the following benchmarks:

- **Inter-Agent Communication:** Agents must effectively transmit data to one another to reach a conclusion.
- **Team Density:** You must implement at least three or more distinct agents.
- **Tool Integration:** Every agent in the system must be equipped with at least one functional tool.
- **Output Integrity:** The system’s final output must match the expected value and format (e.g., if the logic dictates a "Non-Compliant" result, the final action, like an email, must reflect that exactly)


**Pro-Tip:** Remember to use `fmt.Printf()` inside your Go tools. In a multi-agent system, visibility is everything—if you can't see the agents talking, you can't debug the mission!
