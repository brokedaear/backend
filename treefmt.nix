# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Unlicense

{
  # See https://github.com/numtide/treefmt-nix#supported-programs

  projectRootFile = "flake.nix";

  settings.global.includes = [
    "*.go"
    "*.yaml"
    "*.yml"
    "*.md"
    "*.nix"
    "*.proto"
    "*.sql"
  ];

  settings.global.excludes = [
    "*"
  ];

  settings.global.fail-on-change = false;

  programs.gofumpt.enable = true;
  programs.goimports.enable = true;
  programs.protolint.enable = true;
  programs.sqlfluff = {
    enable = true;
    dialect = "postgres";
  };

  # GitHub Actions
  programs.yamlfmt.enable = true;
  programs.actionlint.enable = true;

  programs.taplo.enable = true;

  # Markdown
  programs.mdformat.enable = true;

  # Nix
  programs.nixfmt.enable = true;
}
