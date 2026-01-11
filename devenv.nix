{
  pkgs,
  lib,
  config,
  ...
}:
{
  languages.go.enable = true;

  packages = [
    pkgs.git
    pkgs.pulumi-bin
    pkgs.yamllint
    pkgs.golangci-lint
  ];

  scripts = {

    lint = {
      description = "Lints the go project";
      exec = ''
        cd $(git rev-parse --show-toplevel)/infra
        golangci-lint run
      '';
    };

    test-infra = {
      description = "Test go code";
      exec = ''
        cd $(git rev-parse --show-toplevel)/infra
        go test -v ./...
      '';
    };

    build-project = {
      description = "builds go project";
      exec = ''
        cd $(git rev-parse --show-toplevel)/infra
        go build ./...
      '';
    };

    infra-up = {
      description = "runs pulumi up to standup the infrastructure";
      exec = "pulumi up";
    };

    devhelp = {
      description = "Prints this message";
      exec = ''
        echo 
        echo Helper scripts
        echo
        ${pkgs.gnused}/bin/sed -e 's| |••|g' -e 's|=| |' <<EOF | ${pkgs.util-linuxMinimal}/bin/column -t | ${pkgs.gnused}/bin/sed -e 's|^|🦾 |' -e 's|••| |g'
        ${lib.generators.toKeyValue { } (lib.mapAttrs (name: value: value.description) config.scripts)}
        EOF
        echo
      '';
    };
  };

  # https://devenv.sh/git-hooks/
  git-hooks.hooks = {
    # Commit convention (Conventional Commits)
    convco.enable = true;

    # Go tools - configure to run only in infra directory
    infra-gotest = {
      enable = true;
      name = "Go tests (infra)";
      entry = "cd infra && go test -v ./...";
      files = "^infra/.*\\.go$";
      language = "system";
      pass_filenames = false;
    };
    infra-golangci-lint = {
      enable = true;
      name = "golangci-lint (infra)";
      entry = "cd infra && golangci-lint run";
      files = "^infra/.*\\.go$";
      language = "system";
      pass_filenames = false;
    };
    infra-gofmt = {
      enable = true;
      name = "gofmt (infra)";
      entry = "cd infra && gofmt -l -w .";
      files = "^infra/.*\\.go$";
      language = "system";
      pass_filenames = false;
    };

    # YAML tools
    yamllint.enable = true;
  };

  enterShell = ''
    devhelp
  '';

  # See full reference at https://devenv.sh/reference/options/
}
