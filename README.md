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

```
go install github.com/vtrenton/tshx@latest
```

Homebrew (macOS or Linux):

```
brew install vtrenton/tap/tshx
```

Nix:

```
nix run github:vtrenton/tshx
```

### macOS security warning

macOS will flag the Homebrew-installed binary as being from an
"unidentified developer" (or outright call it malware). This is
because the binary isn't code-signed or notarized by Apple — doing so
requires an Apple Developer Program membership ($99/year), which this
project isn't paying for. The binary is not malicious; this is just
macOS Gatekeeper's default behavior for any unsigned binary downloaded
from the internet.

To run it anyway, clear the quarantine flag after installing:

```
xattr -d com.apple.quarantine "$(brew --prefix)/bin/tshx"
```

**Note:** macOS re-quarantines the binary every time Homebrew downloads a
new one, so you'll need to re-run this command after every
`brew upgrade tshx`, not just the first install.

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
