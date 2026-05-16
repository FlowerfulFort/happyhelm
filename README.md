# happyhelm

`happyhelm` provides the Helm plugin command `helm pick`, an interactive values picker for creating small override `values.yaml` files from chart defaults or deployed release values.

Instead of copying an entire chart's default values, pick only the paths you want to override and print a minimal YAML skeleton.

## Install

```sh
git clone <repo>
cd happyhelm
go build -o happyhelm .
helm plugin install .
```

## Usage

```sh
helm pick <chart> [keyword...]
helm pick release <release> [keyword...]
```

Examples:

```sh
helm pick traefik/traefik nodePort
helm pick traefik/traefik service.type
helm pick traefik/traefik nodePort externalTrafficPolicy
helm pick traefik/traefik service.type -o values/traefik.yaml
helm pick release traefik nodePort -n kube-system
```

Use `--no-tui` to output all matched paths without opening the picker:

```sh
helm pick traefik/traefik nodePort --no-tui
helm pick release traefik nodePort -n kube-system --no-tui
```

## Picker Keys

- Up/Down or `j`/`k`: move
- Space: toggle current item
- `a`: select all / deselect all
- Enter: confirm
- `q` or Esc: cancel

## Output

The command writes YAML to stdout by default, so it can be redirected or piped:

```sh
helm pick traefik/traefik nodePort > values/traefik.yaml
```

With `-o, --output`, parent directories are created automatically:

```sh
helm pick traefik/traefik service.type -o values/traefik.yaml
```

## Development

```sh
go test ./...
go build -o happyhelm .
```
