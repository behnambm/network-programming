# P2P Tic-Tac-Toe

## Introduction
This project enables peer-to-peer (P2P) Tic-Tac-Toe gameplay over a TCP network connection. In this game, one side initiates the TCP socket, while the other side connects to it. To ensure smooth communication and game synchronization between the two parties, an "Initial Data" packet is exchanged at the beginning to configure the game settings on both ends. Once initialized, game data is continually exchanged between the players.

## Data Transmission Strategy
Due to the inherent nature of TCP, which operates as a stream, it can be challenging to determine the exact amount of data to read from the socket. To address this issue, we have implemented a simple yet effective solution.

- **Data Size Indicator**: When sending data, we prepend each message with the size of the data in Big-Endian order. This size indicator allows the receiving side to know how much data to expect, ensuring proper parsing of the incoming messages.

- **Synchronized Reading**: By sending the data size before the actual content, we streamline the process of reading and writing data over the stream. This synchronization mechanism simplifies data exchange between the two players, eliminating the need for complex data tracking or handling.


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
