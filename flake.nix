# modified from https://github.com/akirak/flake-templates/blob/master/node-typescript/flake.nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    # Code QL
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    pre-commit-hooks.url = "github:cachix/git-hooks.nix";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      treefmt-nix,
      pre-commit-hooks,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
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
          htop
          mprocs
        ];

        treefmtEval = pkgs: treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      in
      {
        formatter = pkgs: treefmtEval.${pkgs.system}.config.build.wrapper;
        devShell = pkgs.mkShell {
          buildInputs = [ commonPackages ];
          shellHook = ''
            # eval "$(starship init bash)"
            export PS1='$(printf "\033[01;34m(nix) \033[00m\033[01;32m[%s] \033[01;33m\033[00m$\033[00m " "\W")'
          '';
        };
      }
    );
}
