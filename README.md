<!--
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 12:56:50
 * @LastEditTime: 2025-09-06 11:22:58
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/README.md
-->

# gramferry

Just UDP to TCP, no more chaos. Written in Go.

Currently, it's an stupid pipe.

- Client side should send the first UDP packet first. (Otherwise it doesn't know who is the receiver)
- No auth, no compression, no encryption, no anti-replay... Just no any advanced features. Your UDP server should take care of herself (yes, I call my servers "she" - it's a style choice, not a grammar mistake), and be well. (In most cases, a UDP service without any protections is not safe at all, and you probably don't want to use it)

But it is easy to compile and use.

- No extra firewall configuration is needed.
- No any TUNs or TAPs.
- No any routing table configurations needed in most cases.
- No config files, only two subcommands, and only two flags.
- No root permissions needed.
- No systemd, no OpenRC.
- Pure Go, no CGO.

It's useful when you are using WireGuard and your ISP is QoSing UDP... And that's why I developed it ;-)

## What's new

- Supported IPv6.
- Added pprof flag for performance analyzing and debugging.
- Now a scanner gorutine will scan the connections pool every 60s, and close connections not used over 300s.

## Protocol

\[len:uint16, big-endian\]\[data:bytes\]

How stupid it is... (But it works) And it's even more stupid some ISPs are QoSing UDP in some stupid ways...

## How it works

Before:

UDP Client <-----\[QoS\]-----> UDP Server

UDP Port X <-----\[QoS\]-----> UDP Port X

After:

UDP Client <->  GramFerry Client  <-> GramFerry Server <-> UDP Server

UDP Port X <->  (X/UDP <-> Y/TCP) <-> (Y/TCP <-> X/UDP) <-> UDP Port X

Client:

Listening UDP and transform UDP to TCP streams. Every different UDP remote will be mapped to a different TCP stream.

If a stream was idled for a long time (300s), it will be closed.

Server:

Listening TCP and transform TCP streams to UDP.

Since there are only 2 more bytes per packet, the cost is really small.
