{
  description = "code2svg - A Go server to generate code snippets as SVG images";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        # Use Go 1.24 as requested
        go = pkgs.go_1_24;
        buildGoModule = pkgs.buildGoModule.override { inherit go; };
      in
      {
        packages.code2svg = buildGoModule {
          pname = "code2svg";
          version = "0.1.0";
          src = ./.;

          vendorHash = null;

          subPackages = [ "cmd/code2svg" ];
        };

        packages.default = self.packages.${system}.code2svg;

        apps.code2svg = utils.lib.mkApp {
          drv = self.packages.${system}.code2svg;
        };

        apps.default = self.apps.${system}.code2svg;

        devShells.default = pkgs.mkShell {
          buildInputs = [
            go
            pkgs.just
            pkgs.gopls
            pkgs.go-tools
            pkgs.nixfmt-rfc-style
          ];
        };
      }
    )
    // {
      nixosModules.code2svg =
        {
          config,
          lib,
          pkgs,
          ...
        }:
        {
          imports = [ ./nixos/module.nix ];
          services.code2svg.package = lib.mkDefault self.packages.${pkgs.system}.code2svg;
        };

      nixosModules.default = self.nixosModules.code2svg;
    };
}
