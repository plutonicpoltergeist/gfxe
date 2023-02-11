package main

const usage = `gfx [OPTIONS...] pattern-name

Options:
  -s, --save        Save a pattern
  -l, --list        List available patterns
  -d, --dump        Print the grep command of patterns instead
      --rm          Remove patterns
  -h, --help        Print this helps

Examples:
  gfx aws*
  gfx -d aws*
  gfx --rm aws*
  gfx --save pattern-name '-Hnri' 'search-pattern'

`
