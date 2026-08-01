# Chip-8
Chip-8 emulator that can run chip-8 applications through the command line

During my time at University I learnt how to 
# Quick Start
## Prerequisites
- An installation of [Yggdrasil] (https://yggdrasil-network.github.io/installation.html) is needed to run this application.
    - Once downloaded, users need to run the installer and  add the following ip address as a peer to a file in ProgramFiles/Yggdrasil/yggdrasil.conf: 202:960c:a87:2609:9ec0:29f8:e21d:690a
    - The peers list looks like this:
```
Peers: [
  tcp://a.b.c.d:e
  tls://d.c.b.a:e
  tcp://[a:b:c::d]:e
  tls://[d:c:b::a]:e
]
```
## How to Install
1. Clone repository to your local machine
2. run the following command to download all dependencies and create a virtual environment:
`uv sync`
3. if your IDE does not automatically activate the venv, run `source .venv/bin.activate` in the project directory
## Command to run a application
`python3 main.py {path of chip-8 application here} 0`
- NOTE: the 0 represents a configurable opcode for 8xy6 and 8xye. It shouldn't affect most applications but just in case
- It is considered best practice (FOR THIS PROJECT) to keep the chip-8 application in the same relative path as the project for easy access
## Example Command
`python3 main.py Invaders.ch8 0`

Within this repository are multiple chip 8 applications. Feel free to mess around with them!