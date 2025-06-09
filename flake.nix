# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Unlicense

{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    # Code QL
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    pre-commit-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
      flake-utils,
      pre-commit-hooks,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
        ci-script-name = "run-ci";
        ci-script =
          (pkgs.writeScriptBin ci-script-name (builtins.readFile ./scripts/ci.sh)).overrideAttrs
            (old: {
              buildCommand = "${old.buildCommand}\n patchShebangs $out";
            });

        ciPackages = with pkgs; [
          # Go related
          go # Need that obviously
          gofumpt # Go formatter
          golangci-lint # Local/CI linter
          gotestsum # Pretty tester
          upx # Binary shrinker
          reuse # LICENSE compliance
          nixfmt-rfc-style
          figlet # Terminal text ASCII
          tokei # CLOC
        ];

        devPackages = with pkgs; [
          gopls
          gotools
          stripe-cli # Stripe integration
          lazygit # TUI Git interface
          mprocs # Process runner
          neovim # Better vim
          helix # Quick text editor
          go-task # Makefile alternative
          vegeta # HTTP Load Testing Tool
          openapi-generator-cli
          jq # JSON manipulation
          yq # YAML manipulation
        ];
      in
      {
        formatter = treefmtEval.config.build.wrapper;

        checks = {
          # Throws an error if any of the source files are not correctly formatted
          # when you run `nix flake check --print-build-logs`. Useful for CI
          treefmt = treefmtEval.config.build.check self;
          pre-commit-check = pre-commit-hooks.lib.${system}.run {
            src = ./.;
            hooks = {
              format = {
                enable = true;
                name = "Format with treefmt";
                entry = "${treefmtEval.config.build.wrapper}/bin/treefmt";
                stages = [ "pre-commit" ];
              };
            };
          };

          # ci =
          #   pkgs.runCommand "ci-runner"
          #     {
          #       nativeBuildInputs = [ ci-script ] ++ ciPackages; # Include ci-script and ciPackages
          #       inherit (self.checks.${system}.pre-commit-check) shellHook; # Reuse pre-commit shellHook if needed
          #     }
          #     ''
          #       echo "Running CI script with essential packages..."
          #       # Ensure the ci-script is executable and in PATH for this check
          #       export PATH=$out/bin:$PATH
          #       run-ci
          #       touch $out # Ensure the output path is created
          #     '';
        };

        devShells = {
          default = pkgs.mkShellNoCC {
            buildInputs =
              [ ci-script ]
              ++ ciPackages
              ++ devPackages
              ++ self.checks.${system}.pre-commit-check.enabledPackages;

            # Environment variables
            REUSE_COPYRIGHT = "BROKE DA EAR LLC <https://brokedaear.com>";
            REUSE_LICENSE = "Apache-2.0";

            shellHook = ''
              ${self.checks.${system}.pre-commit-check.shellHook}
              # eval "$(starship init bash)"
              export PS1='$(printf "\033[01;34m(nix) \033[00m\033[01;32m[%s] \033[01;33m\033[00m$\033[00m " "\W")'
            '';
          };

          ci = pkgs.mkShellNoCC {
            buildInputs = [ ci-script ] ++ ciPackages;
            CI = true;
            shellHook = ''
              echo "Entering CI shell. Only essential CI tools available."
              run-ci
            '';
          };
        };
      }
    );
}
