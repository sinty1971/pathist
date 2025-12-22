#!/usr/bin/env python3
import argparse
import platform
import urllib.request
from pathlib import Path

def detect_arch():
    machine = platform.machine().lower()
    if "arm64" in machine or "aarch64" in machine:
        return "arm64"
    if "x86_64" in machine or "amd64" in machine:
        return "x86_64"
    raise RuntimeError(f"Unsupported architecture: {machine}")

def build_url(version: str, arch: str) -> tuple[str, str]:
    if version == "latest":
        return ("https://github.com/bufbuild/buf/releases/latest/download/"
                f"buf-Windows-{arch}.exe", "latest")
    tag = version if version.startswith("v") else f"v{version}"
    return (f"https://github.com/bufbuild/buf/releases/download/{tag}/"
            f"buf-Windows-{arch}.exe", tag)

def main():
    parser = argparse.ArgumentParser(description="Download buf for Windows.")
    parser.add_argument("--version", default="latest")
    parser.add_argument("--destination", default=str(Path.home() / "prj/bin/buf.exe"))
    args = parser.parse_args()

    arch = detect_arch()
    url, tag = build_url(args.version, arch)
    dest = Path(args.destination)
    dest.parent.mkdir(parents=True, exist_ok=True)

    print(f"Downloading buf ({arch}) from {url} ...")
    req = urllib.request.Request(url, headers={"User-Agent": "pathist-buf-downloader"})
    with urllib.request.urlopen(req) as resp, open(dest, "wb") as f:
        f.write(resp.read())
    size = dest.stat().st_size
    print(f"Downloaded version tag: {tag}")
    print(f"Saved to {dest} ({size} bytes)")

if __name__ == "__main__":
    main()
