# Monkey

An interpreter for the Monkey programming language written in Go as described
in [Ball, T. (2016). _Writing An Interpreter In
Go._](https://interpreterbook.com/).

## Features

The Monkey programming language has the following features:

- Variable bindings
- Integers, booleans, strings, arrays, and hash maps
- Arithmetic operations
- First-class and higher-order functions
- Built-in functions
- Closures

## Example

Below is an excerpt of the language.

```text
// Integers & Booleans
let ispositive = fn(value) {
    return value > -1;
};

let crement = fn(value, side) {
    if (ispositive(side)) {
        return value + 1;
    }
    return value - 1;
};

puts(crement(0, 1));  // 1

// Strings
let greeting = fn(first, last) {
    let begin = "Hello, ";
    let end = "!";
    if (len(last) > 0) {
        return begin + first + " " + last + end;
    }
    return begin + first + end;
};

puts(greeting("Monkey", ""));  // Hello, Monkey!

// Arrays
let accumulate = fn(array, stop) {
    if (len(array) == 0) {
        return array;
    }
    if (last(array) == stop) {
        return array;
    }
    let new = push(rest(array), last(array) + 1);
    accumulate(new, stop);
};

puts(accumulate([1, 2, 3], 5));  // [3, 4, 5]

// Hash maps
let h = {"foo": 1, true: 2, 3: 3};
puts(h["foo"]); // 1
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
