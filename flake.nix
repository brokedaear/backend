# modified from https://github.com/akirak/flake-templates/blob/master/node-typescript/flake.nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

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
      ...
    }@inputs:
    let
      pkgs = import nixpkgs {
        config = {
          allowUnfree = true;
        };
      };

      commonPackages = with pkgs; [
        # Development related
        go
        gofumpt
        gopls
        stripe-cli
        upx

        # System tools
        lazygit
        mprocs
        neovim
        helix
        go-task
        vegeta # HTTP Load Testing Tool
      ];

      eachSystem =
        f:
        nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed (system: f nixpkgs.legacyPackages.${system});
      treefmtEval = eachSystem (pkgs: treefmt-nix.lib.evalModule pkgs ./treefmt.nix);
    in
    {
      formatter = eachSystem (pkgs: treefmtEval.${pkgs.system}.config.build.wrapper);

      devShells = eachSystem (pkgs: {
        default = pkgs.mkShell {
          packages = commonPackages;

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
