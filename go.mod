module mvdan.cc/sh/v3

go 1.25.0

require (
	github.com/creack/pty v1.1.24
	github.com/go-quicktest/qt v1.101.0
	github.com/google/go-cmp v0.7.0
	github.com/google/renameio/v2 v2.0.2
	github.com/rogpeppe/go-internal v1.14.1
	golang.org/x/sys v0.42.0
	golang.org/x/term v0.41.0
	mvdan.cc/editorconfig v0.3.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/hugelgupf/go-shlex v0.0.0-20200702092117-c80c9d0918fa // indirect
	github.com/hugelgupf/vmtest v0.0.0-20240307030256-5d9f3d34a58d // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/u-root/gobusybox/src v0.0.0-20250101170133-2e884e4509c7 // indirect
	github.com/u-root/mkuimage v0.0.0-20250905073043-9a40452f5d3b // indirect
	github.com/u-root/u-root v0.15.1-0.20251208185023-2f8c7e763cf8 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	mvdan.cc/sh/moreinterp v0.0.0-local
)

tool golang.org/x/tools/cmd/stringer

replace mvdan.cc/sh/moreinterp => ./moreinterp
