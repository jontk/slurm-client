// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docsOutputDir string
	docsFormat    string
)

func init() {
	docsCmd.Flags().StringVarP(&docsOutputDir, "output", "o", "../../docs/cli", "Output directory for documentation")
	docsCmd.Flags().StringVarP(&docsFormat, "format", "f", "markdown", "Documentation format: markdown, man, rest")
}

var docsCmd = &cobra.Command{
	Use:   "generate-docs",
	Short: "Generate documentation for the CLI",
	Long: `Generate documentation for all CLI commands in various formats.

This command auto-generates comprehensive documentation from the CLI
structure including all commands, subcommands, flags, and examples.

Supported formats:
  - markdown: Markdown files for MkDocs/GitHub
  - man: Manual pages for Unix systems
  - rest: ReStructuredText for Sphinx

Examples:
  # Generate markdown docs in default location
  slurm-cli generate-docs

  # Generate docs in custom directory
  slurm-cli generate-docs --output ./custom-docs

  # Generate man pages
  slurm-cli generate-docs --format man --output ./man
`,
	Hidden: true, // Hide from regular help to avoid confusion
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create output directory if it doesn't exist
		if err := os.MkdirAll(docsOutputDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		absPath, err := filepath.Abs(docsOutputDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		log.Printf("Generating %s documentation in: %s", docsFormat, absPath)

		// Generate documentation based on format
		switch docsFormat {
		case "markdown", "md":
			if err := doc.GenMarkdownTree(rootCmd, absPath); err != nil {
				return fmt.Errorf("failed to generate markdown docs: %w", err)
			}
			log.Println("✓ Markdown documentation generated successfully")

		case "man":
			header := &doc.GenManHeader{
				Title:   "SLURM-CLI",
				Section: "1",
				Source:  "SLURM REST API Client",
			}
			if err := doc.GenManTree(rootCmd, header, absPath); err != nil {
				return fmt.Errorf("failed to generate man pages: %w", err)
			}
			log.Println("✓ Man pages generated successfully")

		case "rest", "rst":
			if err := doc.GenReSTTree(rootCmd, absPath); err != nil {
				return fmt.Errorf("failed to generate ReST docs: %w", err)
			}
			log.Println("✓ ReStructuredText documentation generated successfully")

		default:
			return fmt.Errorf("unsupported format: %s (use: markdown, man, or rest)", docsFormat)
		}

		return nil
	},
}
