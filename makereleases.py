#!/bin/python3

'''
Author: FunctionSir
License: AGPLv3
Date: 2025-09-03 16:13:30
LastEditTime: 2025-09-03 16:54:17
LastEditors: FunctionSir
Description: -
FilePath: /gramferry/makereleases.py
'''

import os


targets = [
    # AIX #
    "aix/ppc64",
    # Darwin #
    "darwin/amd64", "darwin/arm64",
    # Dragonfly #
    "dragonfly/amd64",
    # FreeBSD #
    "freebsd/386", "freebsd/amd64", "freebsd/arm", "freebsd/arm64", "freebsd/riscv64",
    # Illumos #
    "illumos/amd64",
    # Linux #
    "linux/386", "linux/amd64", "linux/arm", "linux/arm64", "linux/loong64", "linux/mips", "linux/mips64",
    "linux/mips64le", "linux/mipsle", "linux/ppc64", "linux/ppc64le", "linux/riscv64", "linux/s390x",
    # NetBSD #
    "netbsd/386", "netbsd/amd64", "netbsd/arm", "netbsd/arm64",
    # OpenBSD #
    "openbsd/386", "openbsd/amd64", "openbsd/arm", "openbsd/arm64", "openbsd/ppc64", "openbsd/riscv64",
    # Plan9 #
    "plan9/386", "plan9/amd64", "plan9/arm",
    # Solaris #
    "solaris/amd64",
    # Windows #
    "windows/386", "windows/amd64", "windows/arm64"
]

for T in targets:
    GOOS, GOARCH = T.split("/")
    EXT = ""
    if GOOS == "windows":
        EXT = ".exe"
    print(f"Building for {T}...")
    os.system(
        f"GOOS={GOOS} GOARCH={GOARCH} go build -o builds/gf-{GOOS}-{GOARCH}{EXT} -ldflags '-s -w'"
    )

print("All done!")
