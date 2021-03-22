# Neoconf - Neovim Config-Tool

Neoconf is a CLI-Tool to help configuring [Neovim](https://github.com/neovim/neovim).

## WIP

It is highly WIP so use at own risk.

## Goals

1. Lua-Free configuration: no need to open the actual configuration files
2. Flexible integration: run from terminal or as a Neovim plugin
3. Abstraction for config: no obscure values you have to google/:help to find out what they do
4. sensible defaults
5. helpers, f.e. install language servers
6. database of plugins incl. installations / stars, searchable

## Basics

Simply make sure neoconf is in your `$PATH`. Download the binary from here or build from source:

```
git clone https://github.com/abenz1267/neoconf
cd neoconf
go install
```

1. `init`: looks for all required folders/files and creates them, if they don't already exist. Also installs all plugins listed in the "plugins.json" file.

## Plugin Management

Barebone management of plugins is in place! Not defining a branch explicitly -> neoconf will first try to clone `master`, if that fails it will try to clone `main`. Right now it is not possible to define `opt` packages. Installation & updates are processed concurrently.

### Commands:

1. `install <plugin1> <plugin2> ...`: installs all plugins provided. Also install missing plugins from "plugins.json". Creates plugin configuration-file under `lua/plugins`
   1. Branch: `glepnir/galaxyline.nvim@SOMEBRANCH`
   2. ....what about post-install hooks, bro? This is a bit sub-optimal at the moment: neoconf will look for a `cd DIR && yarn install` or `yarn install` command in the repository's README file. I'd definitely prefer a `config.json` file with meta-data to be parsed, but this would require work from maintainers.
2. `update`: updates all plugins
3. `clean`: removes plugins not found in "plugins.json", removes `<plugin>.lua` from `lua/plugins`
4. `remove`: this will list all plugins installed and prompt you to enter the index of the plugins you want to remove. Does not remove plugin config.

## Todo

- Configuration management (editor configuration, ...plugin configuration?)
- create a database in order to be able to search for plugins
- neovim plugin!
