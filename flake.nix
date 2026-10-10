{
  description = "A game framework, or something";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:UnstoppableMango/nix-systems";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.inputs.systems.follows = "systems";
    };

    unmango-pkgs = {
      url = "github:unmango/pkgs";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.systems.follows = "systems";
      inputs.flake-parts.follows = "flake-parts";
      inputs.gomod2nix.follows = "gomod2nix";
      inputs.treefmt-nix.follows = "treefmt-nix";
    };

    a2b = {
      url = "github:UnstoppableMango/a2b";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.systems.follows = "systems";
      inputs.flake-parts.follows = "flake-parts";
      inputs.treefmt-nix.follows = "treefmt-nix";
      inputs.mangopkgs.follows = "unmango-pkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports = with inputs; [
        systems.flakeModule or { }
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          pkgs,
          self',
          system,
          inputs',
          ...
        }:
        let
          version = "0.0.1";
          go = pkgs.go_1_27;
          inherit (inputs'.unmango-pkgs.packages) likec4;
          arch = inputs'.a2b.legacyPackages.lib.likec4;
          archSrc = ./docs/architecture;
        in
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = with inputs; [
              gomod2nix.overlays.default
            ];
          };

          packages.default = pkgs.callPackage ./nix { inherit go version; };

          # Relative base and hash history let one build serve from any URL prefix.
          packages.architecture = arch.build {
            name = "architecture";
            src = archSrc;
            base = "./";
            useHashHistory = true;
            title = "unmango/game architecture";
          };

          devShells.default = pkgs.mkShellNoCC {
            packages = [
              go
              (pkgs.gomod2nix.override { inherit go; })
            ]
            ++ (with pkgs; [
              buf
              direnv
              gopls
              ginkgo
              protoc-gen-go
              protoc-gen-connect-go
              gnumake
              nixfmt
            ])
            ++ [ likec4 ];
          };

          checks.architecture = arch.validate {
            name = "architecture";
            src = archSrc;
          };
          checks.architecture-site = self'.packages.architecture;

          treefmt.programs = {
            actionlint.enable = true;
            nixfmt.enable = true;
            gofmt.enable = true;
          };
        };
    };
}
