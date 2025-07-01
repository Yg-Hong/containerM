#!/bin/bash
set -e

# Wait up to 5 seconds for veth interface to appear
for i in {1..50}; do
  VETH=$(ip -o link | awk -F': ' '{print $2}' | grep '@' | cut -d@ -f1 | head -n 1)
  if [[ -n "$VETH" ]]; then
    break
  fi
  sleep 0.1
done

# Exit if no veth found
if [[ -z "$VETH" ]]; then
  echo "[ERROR] VETH interface not found"
  exit 1
fi

# Rename interface and configure
ip link set "$VETH" name eth0
ip link set lo up
ip link set eth0 up
ip addr add 10.0.0.2/24 dev eth0
ip route add default via 10.0.0.1
