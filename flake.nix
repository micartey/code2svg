{
  description = "Code SVG - A Go server to generate code snippets as SVG images";

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
        packages.code-svg = buildGoModule {
          pname = "code-svg";
          version = "0.1.0";
          src = ./.;

          vendorHash = null;

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postInstall = ''
            mkdir -p $out/share/code-svg
            cp code_preview.svg $out/share/code-svg/

            wrapProgram $out/bin/code-svg \
              --run "cd $out/share/code-svg"
          '';
        };

        packages.default = self.packages.${system}.code-svg;

        apps.code-svg = utils.lib.mkApp {
          drv = self.packages.${system}.code-svg;
        };

        apps.default = self.apps.${system}.code-svg;

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
      nixosModules.code-svg =
        {
          config,
          lib,
          pkgs,
          ...
        }:
        {
          imports = [ ./nixos/module.nix ];
          services.code-svg.package = lib.mkDefault self.packages.${pkgs.system}.code-svg;
        };

      nixosModules.default = self.nixosModules.code-svg;
    };
}
