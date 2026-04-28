# The Ad Compliance Analyzer

Build an automated system to screen advertisement transcripts for legal and regulatory compliance.

## The Workflow:
- The Transcripts: You are provided with 5 `json` formatted ad descriptions in the file `ads.md`.
- The Compliance Agents: Define multiple specialized agents. Each agent should check the text against a set of **"legal rules"** defined above.
- The Final Verdict: Each analyzer agent passes a "Compliant" or "Non-Compliant" status to a Final Reviewer Agent.
The Output Logic:
- If Compliant: The system prints **"This ad is compliant"** and terminates.
- If Non-Compliant: The system must trigger a tool that sends an email to the address `eduardo.carvalho@deus.ai` flagging the specific ad and marking it as rejected.

## Ad Compliance Checklist (Legal Rules)
- Shows Minors: The ad features depictions of children or individuals under the legal age of majority.
- Shows Alcohol: The ad features alcoholic beverages, alcohol branding, or the act of consuming alcohol.
- Shows Gambling: The ad features betting activities, casino environments, or games of chance involving monetary stakes.
- Shows Tobacco: The ad features cigarettes, cigars, tobacco packaging, or the act of smoking.
- Shows Weapons: The ad features firearms, ammunition, tactical gear, or bladed weaponry.
- Sexually Suggestive: The ad uses provocative posing, revealing clothing (e.g., crop tops/mini-shorts), or sexualized imagery to market the product.
- Contains Medical Claims: The ad makes health-related assertions or uses medical authority (e.g., doctors, stethoscopes, clinical settings) to imply a product is safe or beneficial.
