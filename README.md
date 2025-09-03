<!--
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 12:56:50
 * @LastEditTime: 2025-09-03 16:01:03
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/README.md
-->

# gramferry

Just UDP to TCP, no more chaos. Written in Go.

Currently, it's an old and stupid pipe.

- IPv6 is not supported (but might be supported it in the future). (Since the parser is currently really simple)
- Client side should send the first UDP packet first. (Otherwise it doesn't know who is the receiver)
- No auth, no compression, no encryption, no anti-replay... Just no any advanced features. Your UDP server should take care of herself (yes, I call my servers "she" - it's a style choice, not a grammar mistake), and be well. (In most cases, a UDP service without any protections is not safe at all, and you probably don't want to use it)
- Currently, it she will not clean the TCP connections automatically. (Might add this feature in the future)

But it is easy to compile and use.

- No extra firewall configuration is needed.
- No any TUNs or TAPs.
- No any routing table configurations needed in most cases.
- No config files, only two subcommands, and only two flags.
- No root permissions needed.
- No systemd, no OpenRC.
- Pure Go, no CGO.

It's useful when you are using WireGuard and your ISP is QoSing UDP... And that's why I developed it ;-)

## Protocol

\[len:uint16, big-endian\]\[data:bytes\]

How stupid it is... (But it works) And it's even more stupid some ISPs are QoSing UDP in some stupid ways...
