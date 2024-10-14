# Yocto Lisp Specification


## 1. Overview

Yocto Lisp is a personal Lisp dialect designed for simplicity, readability, and extensibility. It aims to provide a small set of primitives that allow for easy syntax extension and customization.

Key features:
- Intuitive function names (e.g., `first` instead of `car`)
- Simplified syntax for common operations
- Easy-to-use macro system for extending the language
- Consistent function definition syntax
- Optional infix notation support


## 2. Syntax

### 2.1 Basic Structure
- S-expressions: `(function arg1 arg2 ...)`
- Lists: `[item1 item2 ...]` (syntactic sugar for `(list item1 item2 ...)`)
- Dictionaries: `{key1 value1, key2 value2, ...}`

### 2.2 Special Syntax
- Splat operator: `...` (for unpacking lists or dictionaries)

### 2.3 Comments
- Single-line comments: `;`
- Multi-line comments: `;; ... ;;`


## 3. Data Types

- Numbers (integers and floating-point)
- Strings
- Symbols `:symbolname`
- Booleans (`true` and `false`)
- Lists
- Dictionaries
- Functions
- Nil (represented as `nil`)

## 4. Special Forms

- `def`: Define variables and functions
- `func`: Create anonymous functions
- `if`: Conditional execution
- `quote`: Prevent evaluation
- `macro`: Define macros

## 5. Core Functions

- List operations:
  - `first`: Return the first element of a list (replaces `car`)
  - `rest`: Return all but the first element of a list (replaces `cdr`)
  - `append`: Construct a new list by splatting and grouping args into a list
- Arithmetic operations: `+`, `-`, `*`, `/`
- Comparison operations: `=`, `<`, `>`, `<=`, `>=`
- Boolean operations: `and`, `or`, `not`
- I/O operations: `print`, `read`

## 6. Macros

Macros in Yocto Lisp allow for powerful syntactic abstraction and language extension. They are defined using the `macro` special form and operate on the abstract syntax tree before evaluation.

### 6.1 Macro Definition

Macros are defined using the following syntax:

```lisp
(def name (macro [args] body))
```

### 6.2 Macro Expansion

Macro expansion happens at compile-time, transforming the code before it's evaluated.

### 6.3 Hygiene

Yocto Lisp implements hygienic macros to prevent variable capture and ensure that macros behave consistently regardless of the context in which they are used.

### 6.4 Example Macros

```lisp
; Unless macro
(def unless (macro [test then]
  `(if (not ~test) ~then)))

; While loop macro
(def while (macro [test & body]
  `(loop []
     (when ~test
       ~@body
       (recur)))))
```

## 7. Standard Library

## 7. Standard Library

The Yocto Lisp standard library provides a set of commonly used functions and utilities to facilitate programming in the language.

### 7.1 List Operations

- `last`: Return the last element of a list
- `length`: Return the length of a list
- `map`: Apply a function to each element of a list
- `filter`: Select elements from a list that satisfy a predicate
- `reduce`: Reduce a list to a single value using a combining function

### 7.2 String Operations

- `str-concat`: Concatenate multiple strings
- `str-split`: Split a string into a list of substrings
- `str-join`: Join a list of strings with a separator
- `str-upper`: Convert a string to uppercase
- `str-lower`: Convert a string to lowercase

### 7.3 Math Functions

- none yet

### 7.4 Functional Programming Utilities

- none yet

### 7.5 I/O and System Interaction

- `read-file`: Read the contents of a file
- `write-file`: Write content to a file
- `cmd`: Execute a shell command and return its output

### 7.6 Data Structures

- `dict. key` get dict value for key
- `dict. (key)` same, but programmatically
- ``

### 7.7 Control Flow

- `when`: Conditional execution when a test is true
- `cond`: Multi-way conditional branching
- `loop`: General looping construct

## 8. Special comments

Some type of inline, one word comment method to add comments inside of code that gets ignored

e.g.

`(for x y (do something))`
to
`(for x .in y (do something))`
or
`(for x ~in y (do something))`
or
`(for x ;in y (do something))`
or
`(for x *in y (do something))`
