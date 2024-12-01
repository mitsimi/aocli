# Advent of Code CLI

`aocli` is a CLI client for advent of code.

The purpose of this tool is so you don't have to leave your editor while participating (if you are lazy like me) or if you want to save seconds if you are trying to be competitive.

Please do not use this tool for automated AI assisted submissions. This is actively [discouraged by the AoC team](https://adventofcode.com/about#faq_ai_leaderboard) and against the spirit of AoC.

## Features

> [!WARNING]
> Downloading and saving the puzzles text and your inputs locally is for personal convenience. Uploading them to a public repository is discouraged by the authors of this project and [the AoC team](https://adventofcode.com/about#faq_copying).

- Load Advent of Code session cookie from a file or environment variable.
- Read puzzle description and save it to a file in Markdown format.
- Download puzzle input.
- Submit your puzzle answer and check if it is correct.
- Validate arguments (year, day, puzzle part) and check if puzzle is unlocked.
- If year is not provided, default to the current or last Advent of Code event.
- Infer puzzle day when possible (last unlocked puzzle for current and past events).

## Installation

You can download the latest release from the [releases page](https://github.com/mitsimi/aocli/releases/latest).

Install with go:

```sh
go install https://github.com.com/mitsimi/aocli
```

## Documentation

### Commands

- `new` - Create a new folder for the puzzle and download the puzzle data.
- `download` - Download the puzzle data and save it locally.
- `submit` - Submit your puzzle answer and check if it is correct.

### Configuration

The program looks for a configuration file in the following places:

1. home folder (~)
2. project folder
3. file provided via the flag

The configuration values getting merged in the order above. The last value (flag) wins. This means that you can provide a default configuration in your home folder and override it with a project specific configuration.
This should should help to have the session token somewhere safe and not in the project folder, so it can't be leaked.

The configuration file is either a TOML, YAML or JSON file with the following keys:
| Key | Description | Default | Possible Values |
| ----------- | ---------------------------------------------------------------------------- | -------------------------- | ----------------------- |
| `session` | Your Advent of Code session cookie. | | |
| `year` | The year of the Advent of Code event. Defaults to the current or last event. | current or last event year | 2015, 15, 2020, 20 |
| `structure` | The folder structure for saving puzzles and inputs. | single-year | multi-year, single-year |

## Example

```sh
# You can get the session from the cookies of https://adventofcode.com

mkdir ./template # Feel free to add your code boilerplate in this folder
aocli new -y 2020 -d 1 # This will create the "day_1" folder and downloads the problem into it

# After you solved the problem
cd day_01
aocli submit -l 1 <answer>

# You can also pipe the answer into the cli program
run program | aocli submit -l 1

# Only download the puzzle description
aocli download -D
```
