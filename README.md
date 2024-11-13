# Overview
https://github.com/user-attachments/assets/49c6f8ca-bb9b-4c5b-965c-85c3a415e4fb

# Requirements
An AMD64 CPU is recommended, but you can still run this on other architectures.

## Running on Raspberry PI 4 (aarch64)
1. Install Docker if you haven’t already.
2. Enable AMD64 emulation:
```shell
docker run --privileged --rm tonistiigi/binfmt --install amd64
```
If this doesn’t work, try:
```shell
sudo apt install qemu-user-static -y
```
### Optional: Port Forwarding
To open ports on the host, use the following command (adjust as needed for router forwarding):
```shell
sudo iptables -A INPUT -p udp -m udp --sport 27015:27099 --dport 1025:65355 -j ACCEPT;
```

# Running (Manually)
1. Clone this repository.
2. Set the following environment variables:
- PORT (required)
- MAP (optional)
- MAX_PLAYERS (optional)
- SERVER_FILES_URL (required) – Use steamcmd to generate your own server files, as proprietary files aren’t shared here. Place files in the server directory (flat structure).
3. Start the server:
```shell
./start.sh
```

# Running the app
```shell
make dev
```
