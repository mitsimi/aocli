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

## Documentation

The program looks for a configuration file in the following places in order:

1. file provided via the flag
2. project folder
3. home folder
4. home/.config

### Configuration

The configuration file is either a TOML, YAML or JSON file with the following keys:
| Key         | Description                                                                  | Default                    | Possible Values         |
| ----------- | ---------------------------------------------------------------------------- | -------------------------- | ----------------------- |
| `session`   | Your Advent of Code session cookie.                                          |                            |                         |
| `year`      | The year of the Advent of Code event. Defaults to the current or last event. | current or last event year | 2015, 15, 2020, 20      |
| `structure` | The folder structure for saving puzzles and inputs.                          | single-year                | multi-year, single-year |

## Example


## Installation


## Documentation
