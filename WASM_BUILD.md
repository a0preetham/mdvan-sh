# Building gosh as a WASI binary

Compiles `cmd/gosh` (the shell interpreter) to WebAssembly using the WASI target,
runnable with wasmtime.

## Prerequisites

### Go 1.21+ (WASI support added in 1.21)

```bash
curl -fsSL https://go.dev/dl/go1.24.3.linux-amd64.tar.gz -o /tmp/go.tar.gz
mkdir -p ~/.local/go-sdk
tar -C ~/.local/go-sdk -xzf /tmp/go.tar.gz
export PATH="$HOME/.local/go-sdk/go/bin:$PATH"
go version  # should print go1.24.3
```

### wasmtime

```bash
curl -fsSL https://github.com/bytecodealliance/wasmtime/releases/download/v32.0.0/wasmtime-v32.0.0-x86_64-linux.tar.xz -o /tmp/wasmtime.tar.xz
mkdir -p ~/.local/bin
tar -xJf /tmp/wasmtime.tar.xz -C /tmp
cp /tmp/wasmtime-v32.0.0-x86_64-linux/wasmtime ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
wasmtime --version  # should print wasmtime 32.0.0
```

Add both exports to `~/.bashrc` to persist across sessions.

## Build

```bash
GOOS=wasip1 GOARCH=wasm go build -o gosh.wasm ./cmd/gosh
```

Output: `gosh.wasm` (~5.4 MB, no modifications to the source needed).

## Run

```bash
# Inline command
wasmtime --dir=. gosh.wasm -c 'echo hello; for i in 1 2 3; do echo "item $i"; done'

# Pipe via stdin
echo 'echo hello from wasm' | wasmtime --dir=. gosh.wasm

# Execute a script file (grant the directory containing the script)
wasmtime --dir=/path/to/scripts gosh.wasm /path/to/scripts/myscript.sh
```

`--dir` grants the WASM module read/write access to that host directory.
Use `--dir=/` to grant full filesystem access (less sandboxed).
