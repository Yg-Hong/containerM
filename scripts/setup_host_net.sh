#!/bin/bash
set -e

CONTAINER_PID=$1
VETH_HOST=veth_host_$CONTAINER_PID
VETH_CONT=veth_cont_$CONTAINER_PID
BRIDGE=br0_cm

# Create bridge if it does not exist
if ! ip link show $BRIDGE > /dev/null 2>&1; then
  ip link add name $BRIDGE type bridge
  ip addr add 10.0.0.1/24 dev $BRIDGE
  ip link set $BRIDGE up
fi

# Create veth pair
ip link add $VETH_HOST type veth peer name $VETH_CONT

# Move container side veth to container network namespace
ip link set $VETH_CONT netns $CONTAINER_PID

# Bring up host side veth and add to bridge
ip link set $VETH_HOST up
ip link set $VETH_HOST master $BRIDGE

# Enable IP forwarding
sysctl -w net.ipv4.ip_forward=1

# Setup NAT using iptables
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 ! -o br0 -j MASQUERADE
