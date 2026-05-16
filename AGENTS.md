You are implementing a new Helm plugin project.

Project name: happyhelm
Helm plugin command name: pick
User-facing command: helm pick

Goal:
Build an interactive Helm values picker that helps users generate a minimal override values.yaml file without copying the entire chart default values.

Background:
Helm often prints “Happy Helming!” after install/upgrade. The project name happyhelm is a light reference to that, but the actual command should be serious and practical: `helm pick`.

Core UX:
The user should be able to run:

  helm pick traefik/traefik nodePort

The tool should:
1. Run `helm show values <chart>` internally.
2. Parse the returned YAML.
3. Flatten the values tree into dot-path entries.
4. Filter entries matching the query terms, such as `nodePort`.
5. Show all matching entries in a terminal UI.
6. Allow the user to move with arrow keys or hjkl.
7. Allow the user to toggle selections with Space.
8. Allow Enter to confirm.
9. Print a minimal YAML override skeleton to stdout containing only selected paths.

Example:

Input command:

  helm pick traefik/traefik nodePort

TUI should show something like:

  Search: nodePort

  [ ] ports.traefik.nodePort        null
  [x] ports.web.nodePort            null
  [x] ports.websecure.nodePort      null

  ↑/↓ or j/k: move    space: toggle    enter: confirm    q/esc: quit

After selecting `ports.web.nodePort` and `ports.websecure.nodePort`, stdout should be:

  ports:
    web:
      nodePort: null
    websecure:
      nodePort: null

Implementation language:
Use Go.

Recommended libraries:
- cobra for CLI command structure
- bubbletea for TUI
- bubbles/list if useful for the list UI
- gopkg.in/yaml.v3 for YAML parsing and output

Repository structure:

  happyhelm/
    go.mod
    main.go
    plugin.yaml
    cmd/
      root.go
      pick.go
    internal/
      helm/
        values.go
      values/
        flatten.go
        search.go
        skeleton.go
      tui/
        picker.go

Helm plugin metadata:
Create `plugin.yaml` like:

  name: "pick"
  version: "0.1.0"
  usage: "Pick Helm values and generate minimal override YAML"
  description: "Interactive Helm values picker for creating small override files"
  command: "$HELM_PLUGIN_DIR/happyhelm"

Commands for MVP:

  helm pick <chart> [keyword...]

Examples:

  helm pick traefik/traefik nodePort
  helm pick traefik/traefik service.type
  helm pick traefik/traefik nodePort externalTrafficPolicy
  helm pick traefik/traefik nodePort -o values/traefik.yaml

Flags:
- `-o, --output <file>`: write selected minimal YAML to a file instead of stdout
- `--no-tui`: skip TUI and output all matched paths as YAML
- `--debug`: print debug information to stderr

Behavior details:

1. Fetch chart values
   Implement a function that runs:

     helm show values <chart>

   Capture stdout.
   Return an error if Helm is not installed, the chart is invalid, or the command fails.

2. Parse YAML
   Parse the YAML into a tree.
   Preserve null values as null.
   It is acceptable for MVP to lose comments.

3. Flatten YAML
   Convert nested maps into entries:

     type ValueEntry struct {
       Path  string
       Value any
     }

   Example:

     service:
       type: LoadBalancer
     ports:
       web:
         nodePort: null

   Should become:

     service.type = LoadBalancer
     ports.web.nodePort = null

4. Lists
   For MVP, do not flatten inside list items.
   If a value is a list, treat the list field as a selectable leaf.
   Example:

     additionalArguments:
       - "--api.insecure=true"

   Should become:

     additionalArguments = ["--api.insecure=true"]

   Later we may add better list handling, but avoid unsafe partial list overrides for now.

5. Search matching
   Support:
   - case-insensitive keyword match against the dot path
   - exact-ish path matching, e.g. `service.type`
   - multiple keywords as OR matching for MVP

   Example:

     helm pick traefik/traefik nodePort externalTrafficPolicy

   Should show paths containing either `nodePort` or `externalTrafficPolicy`.

6. TUI selection
   Implement a simple multi-select TUI:
   - up/down arrows move cursor
   - j/k also move cursor
   - Space toggles current item
   - Enter confirms
   - q or Esc cancels
   - a toggles select all / deselect all
   - / enters filter mode if easy; if not, leave it for later

   On cancel, exit with non-zero code and print no YAML.

7. Build minimal override skeleton
   Given selected paths, rebuild a nested map.

   Selected paths:

     service.type
     ports.web.nodePort
     ports.websecure.nodePort

   Output:

     service:
       type: LoadBalancer
     ports:
       web:
         nodePort: null
       websecure:
         nodePort: null

   The selected values should be copied from the default chart values.

8. Output
   If `--output` is not set, print YAML to stdout.
   If `--output` is set, write the YAML file.
   Create parent directories if necessary.
   Do not print noisy logs to stdout because stdout is intended to be pipeable YAML.
   Write status/debug logs to stderr.

9. Error handling
   Handle:
   - no matches found
   - helm command failure
   - invalid YAML
   - TUI cancel
   - output file write failure

10. Tests
   Add unit tests for:
   - flattening nested YAML
   - keyword search
   - rebuilding skeleton from selected paths
   - list handling as a leaf
   - null value preservation

Initial MVP scope:
Do not implement these yet:
- values.schema.json support
- comment preservation
- natural language interpretation
- helm upgrade/diff integration
- existing values.yaml merge
- release-based mode using `helm get values`

Future roadmap, but not required now:
- `helm pick-release <release> -n <namespace> <keyword...>`
- schema description display
- comment extraction from original values.yaml
- editor mode using `$EDITOR`
- merge mode into an existing values.yaml
- helm-diff integration

Expected deliverable:
Produce a complete working Go project with:
- `go.mod`
- source files
- `plugin.yaml`
- README.md
- unit tests
- clear installation instructions

README should include:

  git clone <repo>
  cd happyhelm
  go build -o happyhelm .
  helm plugin install .

And usage examples:

  helm pick traefik/traefik nodePort
  helm pick traefik/traefik service.type -o values/traefik.yaml

Important design constraint:
The plugin name should be `pick`, not `happy`.
The repository/project name is `happyhelm`, but the actual Helm command should be:

  helm pick ...

Keep implementation simple and robust. Prefer a working MVP over advanced features.
