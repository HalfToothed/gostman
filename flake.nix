{
  description = "Gostman flake development environment";

  # Flake inputs
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  # Flake outputs
  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in with pkgs; {
        # Development environment output
        devShells = {
          default = mkShell {
            # The Nix packages provided in the environment
            packages = [
              go_1_23
              golangci-lint
              golangci-lint-langserver
              gopls
              gotools
              direnv
              watchexec
              gnumake
            ];
          };
        };

        packages.default = buildGo123Module {
          pname = "gostman";
          version = "1.1.0";
          src = ./.;
          vendorHash = "sha256-tr4t4zvxPxmFqflnrpTSs9cwnO7dh5CK3hflupFgR0I=";
          nativeBuildInputs = [ installShellFiles ];
          ldflags = [ "-s" "-w" ];
        };
      });
}
