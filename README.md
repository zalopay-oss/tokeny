# Tokeny

Tokeny is a minimal CLI **[TOTP](https://tools.ietf.org/html/rfc6238) (Time-Based One-Time Password)** generator. 

## 1. Installation

### Using Homebrew

```
brew install zalopay-oss/tap/tokeny
```

or

```
brew tap zalopay-oss/tap
brew install tokeny
```

### Using pre-built binaries

Download proper binary from [releases section](https://github.com/zalopay-oss/tokeny/releases) and put into your paths.

### Build from source

You can also clone this repo and build it yourself.

## 2. Usage

Please consult `tokeny --help` for all features' usages.

```bash
NAME:
   tokeny - Another TOTP generator

USAGE:
   tokeny [global options] command [command options] [arguments...]

COMMANDS:
   setup    setup master password
   add      add new entry
   get      get OTP
   delete   delete selected entry
   list     list all entries
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### Master Password

**Tokeny** requires a Master Password for authenticating you against the whole application.

Master Password can be set **only once**, on the very first time you run `tokeny setup`. After that, all other commands will ask for Master Password once for every 5 minutes.

In case you lost your Master Password, the only way to reset it is removing all data (including token entries), located at `$HOME/.tokeny`.

## 3. Caution

Please think twice before using **Tokeny**, since having token generator in your machine may make you lose the benefits of **Two-Factor** Authentication.
