# gfx

A wrapper around grep to avoid typing common patterns.

> **Note**: This is a fork version of [gf](https://github.com/tomnomnom/gf) developed by [Tom Hudson](https://github.com/tomnomnom) that includes several improvements.
>
> The changes made to the `gf` add several new features to the `gfx`. Here is a summary of the changes:
>
> - Expands the available flags, including adding a flag for removing a search pattern(s) and making the flags more concise.
> - Support globbing pattern name for conducting searches, dumping, and removing patterns.
> - The `getPatterns` function was refactored to take a pattern name as a parameter, enabling it to only retrieve the designated pattern instead of all the patterns that are available.
> - `gfx` utilizes the pool lib to concurrently perform various dumping, removal, and search ops with a max. of **10** concurrent goroutines.

## What? Why?

I use grep a *lot*. When auditing code bases, looking at the output of [meg](https://github.com/tomnomnom/meg),
or just generally dealing with large amounts of data. I often end up using fairly complex patterns like this one:

```bash
$ grep -HnrE '(\$_(POST|GET|COOKIE|REQUEST|SERVER|FILES)|php://(input|stdin))' *
```

It's really easy to mess up when typing all of that, and it can be hard to know if you haven't got any
results because there are non to find, or because you screwed up writing the pattern or chose the wrong flags.

## Usage

It's fairly easy!

```console
$ gfx --help
gfx [OPTIONS...] pattern-name

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

```

### Example

I want to run grep with the pattern below all at once:

```console
$ gfx -l 2>&1 | grep google
google-account-type_secrets
google-api-key_secrets
google-calendar-uri_secrets
google-captcha_secrets
google-client-email_secrets
google-client-id_secrets
google-client-secret_secrets
google-cloud-platform-api-key_secrets
google-cloud-platform-oauth_secrets
google-gcp-service-account_secrets
google-maps-api-key_secrets
google-oauth-access-token_secrets
google-oauth-id_secrets
google-oauth_secrets
google-patterns_secrets
google-private-key_secrets
google-url_secrets
```

You can supply a pattern name by entering a glob pattern:

```bash
$ gfx google-*
```

> **Note**: Keep in mind that the pattern name serves as a globbing for both the dump and remove _(--rm)_ flags.

```console
$ gfx -d google-*
[google-account-type_secrets] grep -aHnoPr "google[_-]?account[_-]?type(=| =|:| :)" .
[google-calendar-uri_secrets] grep -aHnoPr "https://www\\.google\\.com/calendar/embed\\?src=[A-Za-z0-9%@&;=\\-_\\./]+" .
[google-client-email_secrets] grep -aHnoPr "google[_-]?client[_-]?email(=| =|:| :)" .
[google-client-secret_secrets] grep -aHnoPr "google[_-]?client[_-]?secret(=| =|:| :)" .
[google-cloud-platform-oauth_secrets] grep -aHnoPr "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com" .
[google-oauth-access-token_secrets] grep -aHnoPr "ya29\\.[0-9A-Za-z\\-_]+" .
[google-oauth_secrets] grep -aHnoPr "(ya29.[0-9A-Za-z-_]+)" .
```

### Pattern Files

The pattern definitions are stored in `~/.gf` as little JSON files that can be kept under version control:

```console
$ cat ~/.gf/php-sources.json
{
    "flags": "-HnrE",
    "pattern": "(\\$_(POST|GET|COOKIE|REQUEST|SERVER|FILES)|php://(input|stdin))"
}
```

To help reduce pattern length and complexity a little, you can specify a list of multiple patterns too:

```console
$ cat ~/.gf/php-sources-multiple.json
{
    "flags": "-HnrE",
    "patterns": [
        "\\$_(POST|GET|COOKIE|REQUEST|SERVER|FILES)",
        "php://(input|stdin)"
    ]
}
```

You can use the `-save` flag to create pattern files from the command line:

```bash
$ gfx --save 'php-serialized' '-HnrE' '(a:[0-9]+:{|O:[0-9]+:"|s:[0-9]+:")'
```

Or if you want to remove 'em:

```bash
$ gfx --rm 'php-serialized'
```

### Using custom engines

There are some amazing code searching engines out there that can be a better replacement for grep.
A good example is [the silver searcher](https://github.com/ggreer/the_silver_searcher).
It's faster (like **way faster**) and presents the results in a more visually digestible manner.

In order to utilize a different engine, add `engine: <other tool>` to the relevant pattern file:

```bash
# Using the silver searcher instead of grep for the aws-keys pattern:
# 1. Adding "ag" engine
# 2. Removing the E flag which is irrelevant for ag

{
  "engine": "ag",
  "flags": "-Hanr",
  "pattern": "([^A-Z0-9]|^)(AKIA|A3T|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{12,}"
}
```

> **Warning**: Different engines use different flags, so in the example above, the flag `E` has to be removed from the `aws-keys.json` file in order for ag to successfully run.


## Install

Download a pre-built binary from [release page](https://github.com/dwisiswant0/gfx/releases), unpack and run! Or

#### with Go

If you've got Go *1.19+* installed and configured you can install `gfx` with:

```console
go install github.com/dwisiswant0/gfx@latest
```

#### from Source

Clone this repository, and run:

```console
go build ./...
install gfx $GOBIN
```

To get started with patterns quickly, see also:

- https://github.com/dwisiswant0/secpat2gf
- https://github.com/dwisiswant0/gf-secrets
- https://github.com/1ndianl33t/Gf-Patterns

## License

`gfx` is distributed under MIT. Contributions are welcome! See `LICENSE` file.
