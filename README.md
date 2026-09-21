# tshx

A `kubectx`-style context switcher for [Teleport](https://goteleport.com/) (`tsh`).

## Purpose

Teleport stores each cluster you've logged into as a profile file in
`~/.tsh/<proxy-domain>.yaml`, and tracks which one is active in
`~/.tsh/current-profile` (a plain text file containing the domain name).

`tshx` lets you switch between these cluster profiles the same way
`kubectx` lets you switch between Kubernetes contexts: pick one from a
list with the arrow keys and hit Enter, and `tshx` updates
`current-profile` for you.

## Install

*Homebrew:*

```
brew install vtrenton/tap/tshx
```

*Nix:*

run only (no install):

```
nix run github:vtrenton/tshx
```

Install into your profile:

```
nix profile add github:vtrenton/tshx
```

*Go:*
```
go install github.com/vtrenton/tshx@latest
```

## Build from source

Requires Go.

```
make build      # builds ./tshx
make install    # go install . onto your $PATH
```

`tshx --version` prints the running version.

## Usage

Run with no arguments to open an interactive picker. Use the up/down
arrow keys to choose a cluster (the current one is pre-selected and
marked `(current)`), then press Enter to switch:

```
tshx
```

Switch directly to a named cluster without the picker:

```
tshx <proxy-domain>
```

Switch back to the previously active cluster:

```
tshx -
```

### Environment variables

- `TSH_HOME` — override the Teleport config directory (defaults to `~/.tsh`).

## License

[MIT](LICENSE)
