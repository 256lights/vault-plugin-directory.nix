# Nix helper for creating a Vault plugin directory

This repository contains a Nix flake
with a single function: `lib.mkPluginDirectory`.
This function takes in the following parameter attribute set:

```plain
{
  pkgs :: nixpkgs;
  plugins :: [ plugin set (see below) ];
}
```

A plugin set is defined as:

```plain
{
  binary :: derivation;
  [ type :: "secret" | "auth" | "database"; ]
  [ pname :: string; ]
  [ version :: string; ]
}
```

`mkPluginDirectory` will return a derivation
with a `libexec/vault-plugins` directory containing all the listed plugins
as well as a `bin/register-vault-plugins` script
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
