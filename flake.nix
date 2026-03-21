# nix flake check && nix build .#default
{
  description = "Smart text replacer (optionally inplace and through a whole dir tree)";
  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";
  };

  outputs =
    { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems =
        f:
        nixpkgs.lib.genAttrs allSystems (
          system:
          f {
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      packages = forAllSystems (
        { pkgs }:
        {
          default = pkgs.buildGoModule {
            name = "sv";
            src = self;
            vendorHash = "sha256-TJ3lSQtW0E2JrznGVm8hOQGVpXjJyXY2guAxku2O9A4%3D";
            goSum = ./go.sum;
            subPackages = [ "cmd/sv" ];
          };
        }
      );
    };
}
