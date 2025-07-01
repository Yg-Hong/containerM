#!/bin/bash
set -e

CONTAINER_PID=$1
VETH_HOST=veth_host_$CONTAINER_PID
BRIDGE=br0_cm

echo "[host] Cleaning up network..."
ip link delete $VETH_HOST || true
ip link delete $BRIDGE || true