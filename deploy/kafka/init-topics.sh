#!/usr/bin/env bash
# Creates AEGIS Kafka topics for local development
# Usage: bash init-topics.sh [BOOTSTRAP_SERVER]
set -euo pipefail

BOOTSTRAP="${1:-localhost:9092}"

echo "Creating AEGIS Kafka topics on $BOOTSTRAP..."

topics=(
  "aegis.events.raw:12:1"
  "aegis.events.detections:12:1"
  "aegis.events.rag:6:1"
  "aegis.control.policy-reload:1:1"
  "aegis.redteam.jobs:3:1"
)

for topic_spec in "${topics[@]}"; do
  IFS=':' read -r name partitions replication <<< "$topic_spec"
  kafka-topics.sh \
    --bootstrap-server "$BOOTSTRAP" \
    --create \
    --if-not-exists \
    --topic "$name" \
    --partitions "$partitions" \
    --replication-factor "$replication"
  echo "  Created: $name"
done

echo "Done."
