# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Unlicense

{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "systems";

    # Code QL
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    pre-commit-hooks.url = "github:cachix/git-hooks.nix";
  };

  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
      systems,
      ...
    }@inputs:
    let
      pkgs = import nixpkgs {
        config = {
          allowUnfree = true;
        };
      };

      commonPackages = with pkgs; [
        # Go related
        go # Need that obviously
        gofumpt # Go formatter
        golangci-lint # Local/CI linter
        gopls
        gotestsum # Pretty tester
        gotools

        # Library related
        stripe-cli # Stripe integration
        upx # Binary shrinker

        # Dev tools
        openapi-generator-cli
        jq # JSON manipulation
        yq # YAML manipulation
        tokei # CLOC
        reuse # LICENSE compliance

        # Formatting
        nixfmt-rfc-style

        # System tools
        lazygit # TUI Git interface
        mprocs # Process runner
        neovim # Better vim
        helix # Quick text editor
        go-task # Makefile alternative
        vegeta # HTTP Load Testing Tool
        figlet # Terminal text ASCII
      ];

      eachSystem =
        f:
        nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed (system: f nixpkgs.legacyPackages.${system});
      treefmtEval = eachSystem (pkgs: treefmt-nix.lib.evalModule pkgs ./treefmt.nix);
    in
    {
      formatter = eachSystem (pkgs: treefmtEval.${pkgs.system}.config.build.wrapper);

      devShells = eachSystem (pkgs: {
        default = pkgs.mkShellNoCC {
          # packages = commonPackages;
          buildInputs = [ ] ++ commonPackages;

          REUSE_COPYRIGHT = "BROKE DA EAR LLC <https://brokedaear.com>";
          REUSE_LICENSE = "Apache-2.0";

          shellHook = ''
            # eval "$(starship init bash)"
            export PS1='$(printf "\033[01;34m(nix) \033[00m\033[01;32m[%s] \033[01;33m\033[00m$\033[00m " "\W")'
          '';
        };
      });

      checks = eachSystem (pkgs: {
        # Throws an error if any of the source files are not correctly formatted
        # when you run `nix flake check --print-build-logs`. Useful for CI
        treefmt = treefmtEval.${pkgs.system}.config.build.check self;
      });
    };
}
