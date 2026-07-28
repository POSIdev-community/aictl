{
  pkgs,
  ...
}:
let
  golangci-lint-2-8 = pkgs.golangci-lint.overrideAttrs (_old: {
    version = "2.8.0";
    src = pkgs.fetchFromGitHub {
      owner = "golangci";
      repo = "golangci-lint";
      rev = "v2.8.0";
      hash = "sha256-w6MAOirj8rPHYbKrW4gJeemXCS64fNtteV6IioqIQTQ=";
    };
    vendorHash = "sha256-/Vqo/yrmGh6XipELQ9NDtlMEO2a654XykmvnMs0BdrI=";
  });

  goimportsReviserArgs = "-company-prefixes github.com/POSIdev-community";
  goimportsReviserBin = "${pkgs.goimports-reviser}/bin/goimports-reviser";
in
{
  env = {
    GOPRIVATE = "github.com/POSIdev-community";
  };

  packages = with pkgs; [
    git

    golangci-lint-2-8
    gotools
    commitizen
    enumer
    go-task
    go-mockery
    go-arch-lint
    goimports-reviser

    # Use nixpkgs gopls as-is: rebuilding against Go 1.25.8 fails
    # (current gopls requires Go ≥ 1.26).
    gopls
  ];

  languages.go = {
    enable = true;
    version = "1.25.8";
    lsp.enable = false;
  };
  delta.enable = true;

  tasks = {
    "aictl:build".exec = "task build";
    "aictl:test".exec = "task test";
    "aictl:lint".exec = "task lint";
    "aictl:generate".exec = "task generate";
    "aictl:doc".exec = "task doc";
    "aictl:check".exec = "task check";
    "aictl:quick".exec = "task quick";
  };

  git-hooks.hooks = {
    gofmt.enable = true;
    # Disabled: nixpkgs govet vs pinned Go 1.25.8 mismatch. Covered by golangci-lint.
    govet.enable = false;
    aictl-check = {
      enable = true;
      name = "task check";
      # Repo script sets PATH from .devenv/profile (works for IDE commits).
      entry = "${pkgs.bash}/bin/bash scripts/pre-commit-check.sh";
      pass_filenames = false;
      stages = [ "pre-commit" ];
    };
    commitizen.enable = true;
    goimports-reviser = {
      enable = true;
      name = "goimports-reviser";
      entry = "${goimportsReviserBin} ${goimportsReviserArgs}";
      files = "\\.go$";
    };
  };
}
