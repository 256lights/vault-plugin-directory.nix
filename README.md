# Nix helper for creating a Vault plugin directory

This repository contains a Nix flake
with a single function: `lib.mkPluginDirectory`.
This function takes in the following parameter attribute set:

```plain
{
  plugins :: [ plugin set (see below) ],
  pkgs :: nixpkgs,
}
```

A plugin set is defined as:

```plain
{
  binary :: derivation,
  [ type :: "secret" | "auth" | "database", ]
  [ pname :: string, ]
  [ version :: string ]
}
```

`mkPluginDirectory` will return a derivation
with a `bin` directory containing all the listed plugins
as well as a `register-vault-plugins.sh` script
that will register all the plugins using the Vault CLI.

## Example

```nix
vault-plugin-directory.lib.mkPluginDirectory {
  inherit pkgs;
  plugins = [
    { binary = my-vault-secret-plugin; }
  ];
}
```

## License

[Apache 2.0](LICENSE)
