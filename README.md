# Chat Server
A Go Chat Server that can receive and distribute network packets concurrently to broadcast messages and handle authentication.

During my time at University I learnt about networking protocols and security terminology such as TCP and Password Salting at a high, abstract level. While I thought this was useful info, I was 
unsure how to implement this practically in my own systems. This project was the outcome of me putting my theoretical knowledge to the test and applying it in a meaningful way.

## Cool features
- Concurrency through goroutines allow multiple connections and network packets to be handled simultaneously, allowing for optimal performance
- Authentication features for this prototype include password salting, sign up, login and database storage in SQLite
- Message broadcasting allows every connected user to receive the exact same msgs in the same order.

## Prerequisites
- An installation of [Yggdrasil](https://yggdrasil-network.github.io/installation.html) is needed to run this application. Without it, the client cannot connect to the backend server.
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
1. Download the client executable file [here](https://github.com/JoseCarlitosGordo/ChatServer/releases/tag/v0.5.0)
2. Run it
3. Enjoy!
