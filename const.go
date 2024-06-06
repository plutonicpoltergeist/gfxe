package main

const usage = `gfex [OPTIONS...] pattern-name

Options:
  -s, --save        Save a pattern
  -l, --list        List available patterns
  -d, --dump        Print the grep command of patterns instead
      --rm          Remove patterns
  -h, --help        Print this helps
  -i, --hidden      Include hidden files and folders

Examples:
  gfxe aws*
  gfxe -d aws*
  gfxe --rm aws*
  gfxe --save pattern-name '-Hnri' 'search-pattern'

`
