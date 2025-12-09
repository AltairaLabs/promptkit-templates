# IoT Predictive Maintenance Demo

An industrial IoT troubleshooter with multi-tenant security, red-team testing, and tool calls.

## Quick Start

### 1. Inspect the configuration
```bash
promptarena config-inspect
```

### 2. Run a test scenario
```bash
promptarena run --scenario scenarios/hardware-faults.scenario.yaml
```

### 3. Red-team security testing
```bash
promptarena run --scenario scenarios/redteam-selfplay.scenario.yaml
```

### 4. View past results
```bash
promptarena view
```

### 5. Compile & deploy with the SDK
```bash
# Compile prompts into a portable pack
packc compile -c config.arena.yaml -o iot-demo.pack.json

# Run with Go SDK
cd sdk-demo
go run . ../iot-demo.pack.json troubleshooting "PUMP-002 is vibrating heavily"
```

## What's Inside

- **Prompts** - IoT troubleshooter with tool access (v1 baseline, v2 with booking)
- **Tools** - Device listing, sensor data, error logs, maintenance scheduling
- **Scenarios** - Hardware faults, normal ops, red-team attacks
- **Personas** - Adversarial attacker, legitimate plant operator
- **SDK Demo** - Go application using compiled pack file
