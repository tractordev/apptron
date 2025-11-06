package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("  Go-to-WASM Compiler Demo: Go 1.25 + Wazero 1.5.0")
	fmt.Println("════════════════════════════════════════════════════════════════")

	demoDir, _ := filepath.Abs("..")
	goroot := filepath.Join(demoDir, "gosrc")
	workdir := filepath.Join(demoDir, "work")
	gocache := os.Getenv("HOME") + "/Library/Caches/go-build"

	os.MkdirAll(workdir, 0755)

	fmt.Println("\n[1/4] Generate import configuration from stdlib cache")
	generateImportcfg(goroot, workdir)

	fmt.Println("\n[2/4] Compile hello.go using Go compiler (WASM) in wazero")
	compileInWazero(goroot, workdir, gocache, demoDir)

	fmt.Println("\n[3/4] Link hello.o using Go linker (WASM) in wazero")
	linkInWazero(goroot, workdir, gocache)

	fmt.Println("\n[4/4] Run the compiled hello.wasm in wazero")
	runInWazero(workdir)

	fmt.Println("\n════════════════════════════════════════════════════════════════")
	fmt.Println("  ✓ Complete: Compiled & ran Go program entirely in WASM")
	fmt.Println("════════════════════════════════════════════════════════════════")
}

func generateImportcfg(goroot, workdir string) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(`
		cd %s
		echo "# auto-generated importcfg" > %s/importcfg.compile
		echo "# auto-generated importcfg" > %s/importcfg.link
		for pkg in $(GO111MODULE=off GOROOT=$(pwd) GOOS=wasip1 GOARCH=wasm ./bin/go list -f '{{ join .Deps "\n" }}' fmt | sort -u); do
			export_path=$(GO111MODULE=off GOROOT=$(pwd) GOOS=wasip1 GOARCH=wasm ./bin/go list -export -f '{{.Export}}' $pkg 2>/dev/null)
			if [ -n "$export_path" ] && [ -f "$export_path" ]; then
				echo "packagefile $pkg=$export_path" >> %s/importcfg.compile
				echo "packagefile $pkg=$export_path" >> %s/importcfg.link
			fi
		done
		export_path=$(GO111MODULE=off GOROOT=$(pwd) GOOS=wasip1 GOARCH=wasm ./bin/go list -export -f '{{.Export}}' fmt)
		echo "packagefile fmt=$export_path" >> %s/importcfg.compile
		echo "packagefile fmt=$export_path" >> %s/importcfg.link
	`, goroot, workdir, workdir, workdir, workdir, workdir, workdir))
	out, _ := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	fmt.Println("  ✓ Generated importcfg files")
}

func compileInWazero(goroot, workdir, gocache, demoDir string) {
	wasm, _ := os.ReadFile(filepath.Join(goroot, "pkg/tool/wasip1_wasm/compile"))
	fmt.Printf("  Loaded compile tool: %.1f MB\n", float64(len(wasm))/1024/1024)

	start := time.Now()
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	helloPath := filepath.Join(demoDir, "hello.go")
	objPath := filepath.Join(workdir, "hello.o")
	importcfgPath := filepath.Join(workdir, "importcfg.compile")

	cfg := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithArgs("compile", "-o", objPath, "-p", "main", "-complete", "-importcfg", importcfgPath, helloPath).
		WithFSConfig(wazero.NewFSConfig().
			WithReadOnlyDirMount(goroot, goroot).
			WithDirMount(workdir, workdir).
			WithReadOnlyDirMount(gocache, gocache).
			WithReadOnlyDirMount(demoDir, demoDir))

	m, _ := r.CompileModule(ctx, wasm)
	defer m.Close(ctx)

	r.InstantiateModule(ctx, m, cfg)
	fmt.Printf("  Compiled in %v\n", time.Since(start))

	if stat, err := os.Stat(objPath); err == nil {
		fmt.Printf("  ✓ Created hello.o: %.1f KB\n", float64(stat.Size())/1024)
	}
}

func linkInWazero(goroot, workdir, gocache string) {
	wasm, _ := os.ReadFile(filepath.Join(goroot, "pkg/tool/wasip1_wasm/link"))
	fmt.Printf("  Loaded link tool: %.1f MB\n", float64(len(wasm))/1024/1024)

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	objPath := filepath.Join(workdir, "hello.o")
	outPath := filepath.Join(workdir, "hello.wasm")
	importcfgPath := filepath.Join(workdir, "importcfg.link")

	cfg := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithArgs("link", "-o", outPath, "-buildmode=exe", "-importcfg", importcfgPath, objPath).
		WithEnv("GOROOT", goroot).
		WithFSConfig(wazero.NewFSConfig().
			WithReadOnlyDirMount(goroot, goroot).
			WithDirMount(workdir, workdir).
			WithReadOnlyDirMount(gocache, gocache))

	m, _ := r.CompileModule(ctx, wasm)
	defer m.Close(ctx)

	r.InstantiateModule(ctx, m, cfg)

	if stat, err := os.Stat(outPath); err == nil {
		fmt.Printf("  ✓ Created hello.wasm: %.1f MB\n", float64(stat.Size())/1024/1024)
	}
}

func runInWazero(workdir string) {
	wasm, _ := os.ReadFile(filepath.Join(workdir, "hello.wasm"))
	fmt.Printf("  Loaded binary: %.1f MB\n", float64(len(wasm))/1024/1024)

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	cfg := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithArgs("hello")

	m, _ := r.CompileModule(ctx, wasm)
	defer m.Close(ctx)

	fmt.Println("  Output:")
	fmt.Print("    ")
	r.InstantiateModule(ctx, m, cfg)
	fmt.Println("  ✓ Execution complete")
}
