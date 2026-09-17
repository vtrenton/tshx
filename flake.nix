{
  description = "kubectx-style context switcher for Teleport (tsh)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "tshx";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-FSqPLnehNEVfyjex4s7kQ4kezLz9zd49QZAvKmeOu9s=";

          ldflags = [ "-s" "-w" "-X main.version=0.1.0" ];

          meta = with pkgs.lib; {
            description = "kubectx-style context switcher for Teleport (tsh)";
            homepage = "https://github.com/vtrenton/tshx";
            license = licenses.mit;
            mainProgram = "tshx";
          };
        };

        apps.default = flake-utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [ pkgs.go pkgs.gopls ];
        };
      });
}
