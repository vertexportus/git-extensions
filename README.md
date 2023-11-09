# Git Extensions
This repository contains two custom Git extension commands written in Go. These commands are designed to enhance your Git experience by providing additional functionality.

## Build and install

> **_TODO:_** Add build actions to generate downloaded binary releases -> [issue #2](https://github.com/vertexportus/git-extensions/issues/3)

Build using `make` command. Copy binaries to any directory in your `PATH` environment variable.

## Commands

### `git auto-config`

Basic setting of name + email on current repository

```shell
git auto-config --name "Your Name" --email "your@email.com"
```

Supports both long and shorthand formats: `[-n|--name]` and `[-e|--email]`
If you have GPG keys set up on your system, you can instead run:

```shell
git auto-config --gpg
```
supports both `[-g|--gpg]`

without both `--name` and `--email`, it will default to looking for GPG keys in your system and presenting you with a list. Choosing one will pick their name+email as the values.
You can also override the GPG-picked values for email and name passing in either `--name` or `--email` (or both) like so:

```shell
git auto-config --gpg --name "Your Name"
```

To use the selected GPG key to also sign commits, you can pass in the parameter `[-s|--sign]`

```shell
git auto-config --gpg --sign
```

finally, you can skip the confirmation in the end and apply immediately by passing in `[-y|--yes]`

### `git search`

Allows you to search on branches and tags on current repository, and act upon the result

```shell
git search "term"
```
will search local branches that include `term` in their name. By default, it prints the selected branch to Stdout.

There are several parameters that modify the behavior of the command in regard to the selected branch:

`[-c|--checkout]` : does a checkout of selected branch

`[-p|--pull]` : will also do a `git pull` on the repository AFTER checkout of selected branch

`[-m|--merge]` : does a merge of selected branch INTO CURRENT BRANCH

### `git clip`

Simple command that allows you to copy to clipboard certain values

```shell
git clip branch
```

copies current repository current branch to clipboard. Passing in `[-p|--print]` will print out to Stdout instead of clipboard.

```shell
git clip branch --search "term"
```
will use internal git search command to search the branches for including `term`, and copy the result to clipboard (prints with `-p`)

## Contributing

Contributions are welcome! Please feel free to submit a pull request if you have a feature you’d like to add.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.