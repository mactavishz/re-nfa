# NFA-based Regex Engine

This project implements a regular expression engine using Non-deterministic Finite Automata (NFA). It's written in Go and provides a simple, yet powerful way to perform regex matching.

The implementation is grounded in the ideas of automata theory and draws inspiration from Ken Thompson's regular expression search algorithm. This solution differs from Ken Thompson's approach of use reverse Polish Notation and a stack to generate the NFA. Instead, it employs a compact LL(1) parser to construct the NFA.

It is used as a learning exercise to understand the inner workings of FSM-based regex engines and the theory behind them not as a production-ready library.

## Features

- NFA-based regex matching
- Support for basic regex operations:
  - Concatenation
  - Alternation (`|`)
  - Kleene star (`*`)
  - Plus (`+`)
  - Optional (`?`)
- Grouping with parentheses

## Installation

To use this regex engine, you need to have Go installed on your system. If you don't have Go installed, you can download it from [the official Go website](https://golang.org/dl/).

1. Clone the repository:

   ```
   git clone https://github.com/mactavishz/re-nfa.git
   ```

2. Navigate to the project directory:

   ```
   cd re-nfa
   ```

3. Build the project:

   ```
   make
   ```

## Usage

After building the project, you can use the regex engine from the command line:

```
./re <regex_pattern> <input_string>
```

If you don't provide an input string, the program will read from stdin.

### Examples

1. Match a simple pattern:

   ```
   ./re "a(b|c)*" "abcbc"
   ```

2. Use stdin for input:

   ```
   echo "abcbc" | ./re "a(b|c)*"
   ```

## Running Tests

To run the unit tests for this project:

```
make test
```

## Project Structure

- `main.go`: Entry point of the application
- `pkg/fsm/nfa.go`: NFA structure and matching logic
- `pkg/utils/parser.go`: Regex parsing logic
- `pkg/utils/tokenizer.go`: Tokenization of regex patterns
- `nfa_test.go`: Unit tests for the NFA matcher

## Limitations

- This implementation does not support some advanced regex features like lookahead/lookbehind, backreferences, or non-greedy matching.
- Performance may vary compared to highly optimized regex engines in standard libraries.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- This project was inspired by the study of formal languages and automata theory.
- Especially, the book "Introduction to the Theory of Computation" by Michael Sipser.
- Russ Cox's article on [Regular Expression Matching Can Be Simple And Fast](https://swtch.com/~rsc/regexp/regexp1.html) was also a great resource.
- Ken Thompson's paper on [Regular Expression Search Algorithm](https://dl.acm.org/doi/10.1145/363347.363387) was a pioneering work in this field.
