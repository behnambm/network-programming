# P2P Tic-Tac-Toe

## Introduction
This project enables peer-to-peer (P2P) Tic-Tac-Toe gameplay over a TCP network connection. In this game, one side initiates the TCP socket, while the other side connects to it. To ensure smooth communication and game synchronization between the two parties, an "Initial Data" packet is exchanged at the beginning to configure the game settings on both ends. Once initialized, game data is continually exchanged between the players.

## Data Transmission Strategy
Due to the inherent nature of TCP, which operates as a stream, it can be challenging to determine the exact amount of data to read from the socket. To address this issue, we have implemented a simple yet effective solution.

- **Data Size Indicator**: When sending data, we prepend each message with the size of the data in Big-Endian order. This size indicator allows the receiving side to know how much data to expect, ensuring proper parsing of the incoming messages.

- **Synchronized Reading**: By **sending the data size before the actual content**, we streamline the process of reading and writing data over the stream. This synchronization mechanism simplifies data exchange between the two players, eliminating the need for complex data tracking or handling.

![Peek 2023-10-03 13-35](https://github.com/behnambm/network-programming/assets/26994700/4c821cc8-9323-4c5a-a92c-59f056107178)


## Run 

**Install dependencies**
```bash
sudo apt install libxcursor-dev libxinerama-dev libxrandr-dev libffi-dev libxi-dev libgl-dev libxxf86vm-dev
```

**Start Server side**
```bash 
go run *.go -port 9090
```

**Start client side**
```bash 
go run *.go -url localhost:9090
```

## Important functions 

Here are the most important methods/functions that can help you a lot to understand the project.


[gui.go#L199](https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L199)https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L199

[gui.go#L240](https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L240)https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L240

[gui.go#L275](https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L275)https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L275


[gui.go#L313](https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L313)https://github.com/behnambm/network-programming/blob/25f480a45570481b0b963e3602eb8b38261a6150/tic-tac-toe/gui.go#L313



