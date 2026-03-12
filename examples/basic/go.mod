module example/basic

go 1.26.1

require git.nix13.pw/scuroneko/laniakea v1.0.0-beta.13

replace (
	git.nix13.pw/scuroneko/laniakea v1.0.0-beta.13 => ../../
)

require (
	git.nix13.pw/scuroneko/extypes v1.2.1 // indirect
	git.nix13.pw/scuroneko/slog v1.0.2 // indirect
	github.com/alitto/pond/v2 v2.7.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/time v0.15.0 // indirect
)
