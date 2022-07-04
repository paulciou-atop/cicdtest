#!/bin/bash
./snmpscan scan --range=192.168.100.0/24
ping -c 10 192.168.100.10
ping -c 10 192.168.100.11
ping -c 10 192.168.100.12
ping -c 10 192.168.100.13
ping -c 10 192.168.100.20
ping -c 10 192.168.100.21
./snmpscan scan --range=192.168.5.0/24
ping -c 10 192.168.5.43

