# Monkey

An interpreter for the Monkey programming language written in Go as described
in [Ball, T. (2016). _Writing An Interpreter In
Go._](https://interpreterbook.com/).

## Features

The Monkey programming language has the following features:

- Variable bindings
- Integers and booleans
- Arithmetic operations
- First-class and higher-order functions
- Closures

## Example

Below is an excerpt of the language.

```text
let ispositive = fn(value) {
    return value > -1;
};

let crement = fn(value, side) {
    if (ispositive(side)) {
        return value + 1;
    }
    return value - 1;
};

crement(0, 1);
```

## Installation

Execute the below commands in your shell to install `monkey`.

```shell
git clone https://github.com/vincentlabelle/monkey.git
cd monkey
go install
```

## Usage

The REPL can be run by executing `monkey` in your shell (if `~/go/bin` is in
your `PATH`), and you can exit the REPL by typing `exit()`.
