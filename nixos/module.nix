{
  config,
  lib,
  ...
}:

let
  cfg = config.services.code2svg;
in
{
  options.services.code2svg = {
    enable = lib.mkEnableOption "code2svg server";

    package = lib.mkOption {
      type = lib.types.package;
      description = "The code2svg package to use.";
    };

    port = lib.mkOption {
      type = lib.types.port;
      default = 8080;
      description = "The port on which the code2svg server listens.";
    };

    openFirewall = lib.mkOption {
      type = lib.types.bool;
      default = false;
      description = "Whether to open the port in the firewall.";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.code2svg = {
      description = "Code SVG Server";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];

      environment = {
        PORT = toString cfg.port;
      };

      serviceConfig = {
        ExecStart = "${cfg.package}/bin/code2svg";
        Restart = "always";
        DynamicUser = true;
        # Basic hardening
        ProtectSystem = "full";
        NoNewPrivileges = true;
        PrivateTmp = true;
      };
    };

    networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.port ];
  };
}
