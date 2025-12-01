#!/usr/bin/env bash
set -euo pipefail

# Quick demo helper for the IoT maintenance template.
# Requires: promptarena CLI, packc, OPENAI_API_KEY.

CONFIG=${CONFIG:-config.arena.yaml}
PACK_OUT=${PACK_OUT:-iot-demo.pack.json}

if [[ ! -f "$CONFIG" ]]; then
  echo "Config $CONFIG not found. Run from the rendered template directory." >&2
  exit 1
fi

echo "== Red-team self-play =="
promptarena run --config "$CONFIG" --scenario scenarios/redteam-selfplay.scenario.yaml

echo "\n== Hardware fault before prompt update (v1) =="
promptarena run --config "$CONFIG" --scenario scenarios/hardware-faults.scenario.yaml --prompt troubleshooter-v1

echo "\n== Hardware fault after prompt update (v2) =="
promptarena run --config "$CONFIG" --scenario scenarios/hardware-faults.scenario.yaml --prompt troubleshooter-v2

echo "\n== Compile pack for SDK demo =="
packc compile -c "$CONFIG" -o "$PACK_OUT"

echo "\nRun SDK demo (adjust prompt/message as needed):"
echo "OPENAI_MODEL=gpt-4o go run ./sdk-demo $PACK_OUT troubleshooter-v2 'PUMP-002 is vibrating heavily'"
