# apptron
[![Discord](https://img.shields.io/discord/415940907729420288?label=Discord)](https://discord.gg/nQbgRjEBU4) ![GitHub Sponsors](https://img.shields.io/github/sponsors/progrium?label=Sponsors)

Local-first development platform

"The amount of amazing technology in this project is staggering. Seriously, star this, it's amazing." -[ibuildthecloud](https://x.com/ibuildthecloud/status/1996979376106492249)

## User Guide

The "project environment" is the main object of Apptron. It is a full Linux 
environment running in the browser with a VSCode-based editor for you to do 
whatever you want with. For example, you can use it as:

* a development environment and editor
* a sandbox for AI and experiments
* an editor to publish static sites
* an embeddable software playground
* a way to run and share Linux software on the web

However, it is fully extendable, customizable, and self-hosted so you could even
use it as the foundation for your own development platform or software system.

Unlike cloud IDEs, Apptron runs entirely in the browser and does not depend on
the cloud. It also only happens to be an IDE, as it is primarily an IDE for
itself as a general compute environment, similar to Smalltalk.

Since it is written mostly in Go, it has first-class language support for Go. 
However, you are encouraged to get other languages to work on it and add them as
supported languages.

### Linux Environment

Apptron runs Alpine Linux with a custom Linux kernel in [v86](https://github.com/copy/v86)
by way of [Wanix](https://github.com/tractordev/wanix), which gives it extra
capabilities such as native Wasm executable support and access to various DOM
APIs through the filesystem. 

The v86 JIT emulator allows 32-bit x86 software to be run, which can be
installed manually or through the Alpine package manager `apk`. A few packages
are pre-installed including `make`, `git`, and `esbuild`. 

### Persistence

Apptron environments are like Docker images in that changes are not persisted
unless committed or added to the environment build script. However, the project
directory, home directory, and public directory are all persisted via browser
storage and cloud synced. Changes outside these directories will be reset
with every page load. However, you can mount more directories backed by browser
storage.

### Virtual Network

In order to install packages, full internet access is provided through a virtual
network. You are given an IP from a virtual DHCP server on this network with 
every page load. This is known as your session IP. Session IPs are routable to
each other, allowing communication across browser tabs and devices.

If you run software that binds to a TCP port on this IP, it will get a public 
HTTPS endpoint. If the service is HTTP, the endpoint will proxy to it for the 
duration the software is running, similar to Ngrok. Non-HTTP TCP services can be 
used over the endpoint tunneled over WebSocket.

### Using Go

Go can be installed via `apk`, but it is better to use the built-in bundle of
Go 1.25 that includes a pre-compiled standard library. Go runs significantly
slower in the browser, so this saves a lot of time with the first build. 

To mount and set up Go, run `source /etc/goprofile`. 

## Developer Guide

### Prerequisites
* Docker
* Go
* npm
* wrangler

### Start Local Apptron
```sh
make dev
```