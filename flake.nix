# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Unlicense

{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    da-flake = {
      url = "github:brokedaear/da-flake";
      inputs.nixpkgs.follows = "nixpkgs";
    };

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
      da-flake,
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
        ci-script = da-flake.lib.${system}.mkScript {
          name = ci-script-name;
          scriptPath = ./scripts/ci.sh;
        };

        ciPackages =
          with pkgs;
          [
            go # Need that obviously
            gofumpt # Go formatter
            golangci-lint # Local/CI linter
            gotestsum # Pretty tester
            buf # protobuf linter/formatter
            protobuf
          ]
          ++ da-flake.lib.${system}.ciPackages;

        devPackages =
          with pkgs;
          [
            gopls
            gotools
            stripe-cli # Stripe integration
            protoc-gen-go
          ]
          ++ da-flake.lib.${system}.devPackages;
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
              lint-go = {
                enable = true;
                name = "Lint Go files";
                entry = "golangci-lint run";
                pass_filenames = false;
                types = [ "go" ];
                stages = [ "pre-push" ];
              };
              unit-tests = {
                enable = true;
                name = "Run unit tests";
                entry = "gotestsum --format testdox ./...";
                pass_filenames = false;
                stages = [ "pre-push" ];
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
          default = pkgs.mkShell {
            buildInputs =
              [ ci-script ]
              ++ ciPackages
              ++ devPackages
              ++ self.checks.${system}.pre-commit-check.enabledPackages;

            inherit (da-flake.lib.${system}.envVars) REUSE_COPYRIGHT REUSE_LICENSE;

            shellHook = ''
              ${self.checks.${system}.pre-commit-check.shellHook}
              # eval "$(starship init bash)"
              export PS1='$(printf "\033[01;34m(nix) \033[00m\033[01;32m[%s] \033[01;33m\033[00m$\033[00m " "\W")'
            '';
          };

          ci = pkgs.mkShell {
            buildInputs = [ ci-script ] ++ ciPackages;
            CI = true;
            shellHook = ''
              echo "Entering CI shell. Only essential CI tools available."
            '';
          };
        };
      }
    );
}
