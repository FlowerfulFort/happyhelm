package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowerfulfort/happyhelm/internal/helm"
	"github.com/flowerfulfort/happyhelm/internal/tui"
	pickvalues "github.com/flowerfulfort/happyhelm/internal/values"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type pickOptions struct {
	output    string
	namespace string
	noTUI     bool
	debug     bool
}

var pickOpts pickOptions

var pickCmd = &cobra.Command{
	Use:   "pick <chart> [keyword...]",
	Short: "Pick Helm values and generate minimal override YAML",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPick(args[0], args[1:], pickOpts)
	},
}

var releaseCmd = &cobra.Command{
	Use:   "release <release> [keyword...]",
	Short: "Pick values from a deployed Helm release",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPickRelease(args[0], args[1:], pickOpts)
	},
}

func init() {
	addPickFlags(pickCmd)
	addPickFlags(releaseCmd)
	releaseCmd.Flags().StringVarP(&pickOpts.namespace, "namespace", "n", "", "release namespace")
	rootCmd.AddCommand(pickCmd)
	rootCmd.AddCommand(releaseCmd)
}

func addPickFlags(command *cobra.Command) {
	command.Flags().StringVarP(&pickOpts.output, "output", "o", "", "write selected YAML to a file")
	command.Flags().BoolVar(&pickOpts.noTUI, "no-tui", false, "skip TUI and output all matched paths")
	command.Flags().BoolVar(&pickOpts.debug, "debug", false, "print debug information to stderr")
}

func runPick(chart string, keywords []string, opts pickOptions) error {
	raw, err := helm.ShowValues(chart)
	if err != nil {
		return err
	}
	debugf(opts.debug, "fetched %d bytes from helm show values\n", len(raw))
	return runPickFromYAML(raw, keywords, opts)
}

func runPickRelease(release string, keywords []string, opts pickOptions) error {
	raw, err := helm.GetReleaseValues(release, opts.namespace)
	if err != nil {
		return err
	}
	debugf(opts.debug, "fetched %d bytes from helm get values --all\n", len(raw))
	return runPickFromYAML(raw, keywords, opts)
}

func runPickFromYAML(raw []byte, keywords []string, opts pickOptions) error {
	var root any
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return fmt.Errorf("invalid values YAML: %w", err)
	}

	entries := pickvalues.Flatten(root)
	matches := pickvalues.Search(entries, keywords)
	debugf(opts.debug, "flattened %d entries, matched %d entries\n", len(entries), len(matches))

	if len(matches) == 0 {
		return fmt.Errorf("no values matched %q", keywords)
	}

	selected := matches
	if !opts.noTUI {
		var err error
		selected, err = tui.Pick(matches, keywords)
		if err != nil {
			if errors.Is(err, tui.ErrCanceled) {
				return fmt.Errorf("selection canceled")
			}
			return err
		}
		if len(selected) == 0 {
			return fmt.Errorf("no values selected")
		}
	}

	out, err := pickvalues.BuildSkeletonYAML(selected)
	if err != nil {
		return err
	}

	if opts.output == "" {
		fmt.Print(string(out))
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(opts.output), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(opts.output, out, 0o644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", opts.output)
	return nil
}

func debugf(enabled bool, format string, args ...any) {
	if enabled {
		fmt.Fprintf(os.Stderr, "debug: "+format, args...)
	}
}
