# IoT Predictive Maintenance Demo Template

An end-to-end PromptArena template that shows multi-tenant security, red-team self-play, and iterative prompt updates for an industrial IoT troubleshooter.

## Contents
- Arena config with two prompt versions (v1 baseline, v2 with engineer booking)
- Mock tool definitions for devices, telemetry, logs, maintenance, and booking
- Self-play personas (red-team attacker, plant operator) with conversation-level assertions
- Scenarios for normal ops, hardware faults, red-team pressure tests, and real-user flow
- Judge prompt + provider for safety scoring
- SDK demo stub (`sdk-demo/main.go`) for running against a compiled pack

## Render the template
```bash
promptarena templates fetch --template iot-maintenance-demo --version 1.0.0
promptarena templates render \
  --template iot-maintenance-demo \
  --version 1.0.0 \
  --values values.example.yaml \
  --out ./iot-maintenance-demo
```

## Run PromptArena scenarios
From the rendered output directory you can choose:
- `config.arena.yaml` (default/OpenAI + mock tools; needs OPENAI_API_KEY for self-play personas/judge)
- `mock.arena.yaml` (all-mock: personas/judge/providers mocked)

```bash
# Red-team self-play (all-mock)
PROMPTKIT_SCHEMA_SOURCE=local promptarena run --config mock.arena.yaml --scenario scenarios/redteam-selfplay.scenario.yaml --provider mock-provider

# Legitimate operator self-play (all-mock, judge disabled in this config)
PROMPTKIT_SCHEMA_SOURCE=local promptarena run --config mock.arena.yaml --scenario scenarios/realuser-selfplay.scenario.yaml --provider mock-provider

# Compare prompt iterations
#   v1 config only loads the baseline prompt
PROMPTKIT_SCHEMA_SOURCE=local promptarena run --config config-v1.arena.yaml --scenario scenarios/hardware-faults.scenario.yaml
#   v2 config loads the booking-enabled prompt
PROMPTKIT_SCHEMA_SOURCE=local promptarena run --config config.arena.yaml --scenario scenarios/hardware-faults.scenario.yaml

# Or run the fully mocked hardware flow
PROMPTKIT_SCHEMA_SOURCE=local promptarena run --config mock.arena.yaml --scenario scenarios/hardware-faults.scenario.yaml --provider mock-provider
```

## SDK demo
1) Compile a pack from the arena (requires packc on PATH):
```bash
packc compile -c config.arena.yaml -o iot-demo.pack.json
```
2) Run the Go demo (needs Go + OPENAI_API_KEY):
```bash
cd sdk-demo
OPENAI_MODEL=gpt-4o go run . ../iot-demo.pack.json troubleshooter-v2 "PUMP-002 is vibrating heavily"
```

Adjust `values.example.yaml` to change defaults (customer ID/name, model choices) before rendering.
