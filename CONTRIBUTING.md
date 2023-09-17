# Contributing Guide

Hello! We'd love to see your contribution on this repository soon, even if it's just a typo fix!

Contributing means anything from reporting bugs, ideas, suggestion, code fix, even new feature.

Bear in mind to keep your contributions under the [Code of Conduct](./.github/CODE_OF_CONDUCT.md) for the community.

## Bug report, ideas, and suggestion

The [issues](https://github.com/teknologi-umum/polarite/issues) page is a great way to communicate to us.
Other than that, we have a [Telegram group](https://t.me/teknologi_umum_v2) that you can discuss your ideas into.
If you're not an Indonesian speaker, it's 100% fine to talk in English there.

Please make sure that the issue you're creating is in as much detail as possible. Poor communication might lead to a big
mistake, we're trying to avoid that.

## Pull request

**A big heads up before you're writing a breaking change code or a new feature: Please open up an
[issue](https://github.com/teknologi-umum/polarite/issues) regarding what you're working on, or just talk in the
[Telegram group](https://t.me/teknologi_umum_v2).**

### Prerequisites

You will need a few things to get things working:

1. Latest stable version of [Rust](https://www.rust-lang.org/tools/install).
2. An IDE or text editor with Rust plugin installed.

### Getting Started

1. [Fork](https://help.github.com/articles/fork-a-repo/) this repository to your own Github account
   and [clone](https://help.github.com/articles/cloning-a-repository/) it to your local machine.
2. Run `cargo update` to install the dependencies needed.
3. Run `cargo run` to start the development application.
4. Have fun!

You are encouraged to use [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0-beta.2/)
for your commit message. But it's not really compulsory.

### Testing your change

Creating tests are not necessary, but it's always great if you can provide tests!

```sh
$ cargo test
```

### Before creating a PR

Please test (command above) and format your code accordingly to pass the CI.

```sh
$ cargo fmt
```

And you're set!