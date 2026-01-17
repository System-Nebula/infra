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

    preview = {
      description = "Displays the expected infra";
      exec = ''
        cd $(git rev-parse --show-toplevel)/infra
        pulumi preview
      '';
    };

    build = {
      description = "builds go project";
      exec = ''
        cd $(git rev-parse --show-toplevel)/infra
        go build
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
        cat << EOF
        ##############################################
            .___          .__           .__           
          __| _/_______  _|  |__   ____ |  | ______   
         / __ |/ __ \  \/ /  |  \_/ __ \|  | \____ \  
        / /_/ \  ___/\   /|   Y  \  ___/|  |_|  |_> > 
        \____ |\___  >\_/ |___|  /\___  >____/   __/  
             \/    \/          \/     \/     |__|
        ##############################################

        EOF
        ${pkgs.gnused}/bin/sed -e 's| |••|g' -e 's|=| |' <<EOF | ${pkgs.util-linuxMinimal}/bin/column -t | ${pkgs.gnused}/bin/sed -e 's|^|* |' -e 's|••| |g'
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
    golangci-lint = {
      enable = true;
      name = "golangci-lint";
      files = "^infra/.*\\.go$";
      entry = "bash -c 'cd infra && golangci-lint run'";
      pass_filenames = false;
    };

    gotest = {
      enable = true;
      name = "gotest";
      files = "^infra/.*\\.go$";
      entry = "bash -c 'cd infra && go test ./...'";
      pass_filenames = false;
    };

    gofmt = {
      enable = true;
      files = "^infra/.*\\.go$";
    };

    # YAML tools
    yamllint = {
      enable = true;
      files = "^(?:(?!\\.github/).)*\\.ya?ml$";
    };
  };

  enterShell = ''
    devhelp
  '';

  # See full reference at https://devenv.sh/reference/options/
}
