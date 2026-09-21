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

See [Nix](#nix) below for declarative NixOS/home-manager setup or
non-flake installs.

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

## Nix

One-off run, no install:

```
nix run github:vtrenton/tshx
```

Install into your profile:

```
nix profile install github:vtrenton/tshx
```

No flakes:

```
nix-env -if https://github.com/vtrenton/tshx/archive/refs/heads/master.tar.gz
```

Declarative (NixOS, home-manager, or your own flake) — add as an input,
then either reference `tshx.packages.${system}.default` directly or pull
in `tshx.overlays.default` to get `pkgs.tshx`:

```nix
{
  inputs.tshx.url = "github:vtrenton/tshx";

  outputs = { nixpkgs, tshx, ... }: {
    nixosConfigurations.myhost = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        { nixpkgs.overlays = [ tshx.overlays.default ]; }
        { environment.systemPackages = [ pkgs.tshx ]; }
      ];
    };
  };
}
```

Same pattern in a home-manager flake (`home.packages` instead of
`environment.systemPackages`).

## License

[MIT](LICENSE)
