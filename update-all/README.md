# update-all

Running daily routines

## Preview

![update-all img preview][preview]

[preview]: ./img/preview.png

## Installation

1. build the tool with:

```sh
cargo build --release
```

2. The executable binary will reside in the folder `target/release/update-all`

## How to use

**1. create a config file**

```sh
./update-all --create
```

(Use `./update-all -h` to get help messages)

By default, the config file will be created in `~/.config/update-all/update-all.config.yaml`

**2. override the config file**

Edit the config file, extend the list of content as you wish.

The program would read the config file and execute the commands you specified in it.

```yaml
# ~/.config/update-all/update-all.config.yaml`
- interval_minute: 10
  name: pyenv
  args:
    - update
- interval_minute: 60
  name: brew
  args:
    - update
- interval_minute: 0
  name: ls
  args:
    - "-alh"
```

|       name        |                               meaning                               |
| :---------------: | :-----------------------------------------------------------------: |
| `interval_minute` |          the minimum numebr of minte to rerun this command          |
|      `name`       |                         program name to run                         |
|      `args`       | list of arguments(use`""` if the argument contains `-` sign inside) |
