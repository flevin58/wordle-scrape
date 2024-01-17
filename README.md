# wordle-scrape

This project program creates a words.txt file with a dictionary of 5-letter words compatible with the game of Wordle.
I used it to copy the generated file into my [wordle game repo](https://github.com/flevin58/wordle)
It could be useful if you need to generate a file with words other than english.

## Installation

No installation is needed, since this program is run only once to generate the words.txt file.
Download the project files to your computer and within the wordle folder type:

```
go run main.go
```

The wordle repository can be found here: https://github.com/flevin58/wordle
There is already a words.txt file

## Words dictionary

The words were scraped from the following site:
https://www.wordunscrambler.net/word-list/wordle-word-list

## Project Tree Structure

```
.
├── LICENSE
├── README.md
├── go.mod
├── go.sum
└── main.go
```
