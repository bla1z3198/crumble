#!/bin/bash
echo "setup tun..."

sudo ip addr add 10.10.0.1/24 dev tun0
sudo ip link set tun0 up

echo "routing..."

sudo ip route add 10.10.0.0/24 dev tun0

echo "ping test"

ping -I tun0 10.10.0.2
