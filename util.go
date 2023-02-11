package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"os/exec"
	"path/filepath"

	"github.com/gobwas/glob"
)

func doSearch(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
}

func listPatterns() ([]string, error) {
	patDir, err := getPatternDir()
	if err != nil {
		return nil, fmt.Errorf(errGetPatternDir, err)
	}

	files, err := filepath.Glob(filepath.Join(patDir, "*.json"))
	if err != nil {
		return nil, err
	}

	var out []string
	for _, file := range files {
		base := filepath.Base(file)
		out = append(out, strings.TrimSuffix(base, ".json"))
	}

	return out, nil
}

func getPatterns(patGlob string) ([]pattern, error) {
	var pats []pattern

	list, err := listPatterns()
	if err != nil {
		return pats, err
	}

	g := glob.MustCompile(patGlob)
	for _, l := range list {
		if g.Match(l) {
			pat, err := getPattern(l)
			if err != nil {
				return pats, err
			}

			pats = append(pats, pat)
		}
	}

	return pats, nil
}

func getPattern(patName string) (pattern, error) {
	var pat pattern

	patDir, err := getPatternDir()
	if err != nil {
		return pat, errors.New(errOpenUserPatternDir)
	}
	fpath := filepath.Join(patDir, fmt.Sprint(patName, ".json"))

	f, err := os.Open(fpath)
	if err != nil {
		return pat, fmt.Errorf(errNoPattern, patName)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&pat); err != nil {
		return pat, fmt.Errorf(errPatternFileMalformed, fpath, err)
	}

	pat.Filename = strings.TrimSuffix(filepath.Base(fpath), ".json")
	pat.Filepath = fpath

	if pat.Pattern == "" {
		if len(pat.Patterns) == 0 {
			return pat, fmt.Errorf(errPatternFileNoPattern, fpath)
		}
		pat.Pattern = fmt.Sprint("(", strings.Join(pat.Patterns, "|"), ")")
	}

	return pat, nil
}

func getPatternDir() (string, error) {
	home := os.Getenv("HOME")

	configPath := filepath.Join(home, ".config/gf")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	legacyPath := filepath.Join(home, ".gf")
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath, nil
	}

	return "", fmt.Errorf(errPatternDirNotFound, configPath, legacyPath)
}

func savePattern(name, flags, pat string) error {
	if name == "" {
		return errors.New(errPatternNameEmpty)
	}

	if flags == "" {
		return errors.New(errPatternFlagsEmpty)
	}

	if pat == "" {
		return errors.New(errPatternValueEmpty)
	}

	p := &pattern{
		Flags:   flags,
		Pattern: pat,
	}

	patDir, err := getPatternDir()
	if err != nil {
		return fmt.Errorf(errGetPatternDir, err)
	}

	path := filepath.Join(patDir, fmt.Sprint(name, ".json"))
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf(errCreatePatternFile, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	if err := enc.Encode(p); err != nil {
		return fmt.Errorf(errWritePatternFile, err)
	}

	return nil
}

func isStdin() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
