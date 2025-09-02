<!--
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-09-02 12:56:50
 * @LastEditTime: 2025-09-02 17:10:16
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /gramferry/README.md
-->

# gramferry

Just UDP to TCP, no more chaos. Written in Go.

Currently, it's an old and stupid pipe.

- IPv6 is not supported (but might be support it in the future).
- Only SINGLE user is allowed when using the UDP Port in client side currently.
- Client side should send the first UDP packet first.
- No auth, no compression, no anti-replay... Just no any advanced features. Your UDP server should take care of herself, and be well.

But it is easy to compile and use.

- No any extra firewall configurations needs to be done.
- No any TUNs or TAPs.
- No any routing table configurations needed in most cases.
- No any config files, only two subcommands, and only two flags.
- No root permissions needed.
- No systemd, no OpenRC.
- Pure Go, no CGO.

It's usefull when you are using WireGuard and your ISP is QoSing UDP... And that's why I developed it ;-)

## Protocol

Currrently: \[len:uint16\]\[data:bytes\]
In the future (may be...): \[linkId:uint32\]\[len:uint16\]\[data:bytes\]

How stupid it is...
