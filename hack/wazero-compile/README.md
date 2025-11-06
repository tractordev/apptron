# wazero-compile

This is a PoC of the Go compiler compiled to WASI, run in Wazero to compile
a hello world program to WASI that can also run in Wazero.

It also gives you an idea of how to use `compile` and `link` directly to
build a Go program.


## Run Demo

Even though the setup builds Go, you also need Go installed already.

```sh
make run-demo
```