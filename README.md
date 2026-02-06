# code2svg

![img](./pkg/code2svg/code_preview.svg)

## Introduction

code2svg is a simple Go-based microservice that transforms source code into beautiful, syntax-highlighted SVG images. 
Inspired by the aesthetics of modern code editors like NvChad, it features:

### Setup

To run the server locally, you only need to have Go installed.

```bash
# Clone the repository
git clone https://github.com/micartey/code2svg.git

# Navigate into the directory
cd code2svg

# Run the server using just
just run
```

#### Nix & NixOS

If you are using Nix, you can run the server directly:

```bash
nix run github:micartey/code2svg
```

To host it on a NixOS server, add the flake to your inputs and import the module:

```nix
{
  inputs.code2svg.url = "github:micartey/code2svg";

  outputs = { self, nixpkgs, code2svg }: {
    nixosConfigurations.my-server = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        code2svg.nixosModules.default
        {
          services.code2svg = {
            enable = true;
            port = 8080;
            openFirewall = true;
          };
        }
      ];
    };
  };
}
```

#### Build from source

If you prefer to build the binary yourself:

```bash
go build -o server .
./server
```

The server will start on `:8080` by default.

## Usage

The primary endpoint is `/svg`. You can provide the code either as a query parameter or in the request body.

### Query Parameters

| Parameter | Type   | Description                                      |
|-----------|--------|--------------------------------------------------|
| `code`    | string | Base64 encoded source code to be rendered. |

```bash
curl "http://localhost:8080/svg?code=Zm4gbWFpbigpIHsKICAgIHByaW50bG4hKCJIZWxsbyIpOwp9"
```

### Request Body

Alternatively, you can send the base64 encoded code directly in the request body.

```bash
cat <file> | base64 | curl -d @- http://localhost:8080/svg
```

