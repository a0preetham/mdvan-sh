package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/core/base64"
	"github.com/u-root/u-root/pkg/core/cat"
	"github.com/u-root/u-root/pkg/core/chmod"
	"github.com/u-root/u-root/pkg/core/cp"
	"github.com/u-root/u-root/pkg/core/gzip"
	"github.com/u-root/u-root/pkg/core/mkdir"
	"github.com/u-root/u-root/pkg/core/mktemp"
	"github.com/u-root/u-root/pkg/core/mv"
	"github.com/u-root/u-root/pkg/core/rm"
	"github.com/u-root/u-root/pkg/core/shasum"
	"github.com/u-root/u-root/pkg/core/tar"
	"github.com/u-root/u-root/pkg/core/touch"
	"github.com/u-root/u-root/pkg/core/xargs"

	"mvdan.cc/sh/v3/interp"
)

var commandBuilders = map[string]func() core.Command{
	"base64": func() core.Command { return base64.New() },
	"cat":    func() core.Command { return cat.New() },
	"chmod":  func() core.Command { return chmod.New() },
	"cp":     func() core.Command { return cp.New() },
	"gzcat":  func() core.Command { return gzip.New("gzcat") },
	"gzip":   func() core.Command { return gzip.New("gzip") },
	"gunzip": func() core.Command { return gzip.New("gunzip") },
	"mkdir":  func() core.Command { return mkdir.New() },
	"mktemp": func() core.Command { return mktemp.New() },
	"mv":     func() core.Command { return mv.New() },
	"rm":     func() core.Command { return rm.New() },
	"shasum": func() core.Command { return shasum.New() },
	"tar":    func() core.Command { return tar.New() },
	"touch":  func() core.Command { return touch.New() },
	"xargs":  func() core.Command { return xargs.New() },
}

func wasiCoreutils(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		program, programArgs := args[0], args[1:]
		hc := interp.HandlerCtx(ctx)

		switch program {
		case "ls":
			return builtinLs(hc, programArgs)
		case "find":
			return builtinFind(hc, programArgs)
		}

		newCmd, ok := commandBuilders[program]
		if !ok {
			return next(ctx, args)
		}

		cmd := newCmd()
		cmd.SetIO(hc.Stdin, hc.Stdout, hc.Stderr)
		cmd.SetWorkingDir(hc.Dir)
		cmd.SetLookupEnv(func(key string) (string, bool) {
			v := hc.Env.Get(key)
			return v.Str, v.Set
		})
		return cmd.RunContext(ctx, programArgs...)
	}
}

// builtinLs is a WASI-compatible ls using only stdlib.
func builtinLs(hc interp.HandlerContext, args []string) error {
	long := false
	paths := []string{}
	for _, a := range args {
		if a == "-l" {
			long = true
		} else if !strings.HasPrefix(a, "-") {
			paths = append(paths, a)
		}
	}
	if len(paths) == 0 {
		paths = []string{hc.Dir}
	}

	for _, p := range paths {
		if !filepath.IsAbs(p) {
			p = filepath.Join(hc.Dir, p)
		}
		info, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(hc.Stderr, "ls: %v\n", err)
			continue
		}
		if !info.IsDir() {
			printEntry(hc, info, p, long)
			continue
		}
		entries, err := os.ReadDir(p)
		if err != nil {
			fmt.Fprintf(hc.Stderr, "ls: %v\n", err)
			continue
		}
		if len(paths) > 1 {
			fmt.Fprintf(hc.Stdout, "%s:\n", p)
		}
		for _, e := range entries {
			fi, _ := e.Info()
			printEntry(hc, fi, e.Name(), long)
		}
	}
	return nil
}

func printEntry(hc interp.HandlerContext, fi fs.FileInfo, name string, long bool) {
	if long {
		fmt.Fprintf(hc.Stdout, "%s %8d %s %s\n",
			fi.Mode(), fi.Size(), fi.ModTime().Format("Jan _2 15:04"), name)
	} else {
		fmt.Fprintln(hc.Stdout, name)
	}
}

// builtinFind is a WASI-compatible find using only stdlib.
func builtinFind(hc interp.HandlerContext, args []string) error {
	root := hc.Dir
	name := ""
	typeFilter := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-name":
			if i+1 < len(args) {
				i++
				name = args[i]
			}
		case "-type":
			if i+1 < len(args) {
				i++
				typeFilter = args[i]
			}
		default:
			if !strings.HasPrefix(args[i], "-") {
				p := args[i]
				if !filepath.IsAbs(p) {
					p = filepath.Join(hc.Dir, p)
				}
				root = p
			}
		}
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(hc.Stderr, "find: %v\n", err)
			return nil
		}
		if typeFilter == "f" && d.IsDir() {
			return nil
		}
		if typeFilter == "d" && !d.IsDir() {
			return nil
		}
		if name != "" {
			matched, err := filepath.Match(name, d.Name())
			if err != nil || !matched {
				return nil
			}
		}
		fmt.Fprintln(hc.Stdout, path)
		return nil
	})
}

