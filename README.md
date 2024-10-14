# Yocto Lisp

Yocto lisp is a project that will serve as my own personal lisp dialect for solo projects.

I want it to be terse but descriptive.

For example, `first` or `head` insteal of `car`.
`fn` instead of `lambda`. In general I just want it to feel right and make sense to me and not use 'jargon' from old lisp.

I like lisp/scheme, but every dialect has some things I don't like.

I want a small amount of primitives where it's easy to built out your own syntax with.

My end goal is to have an executable that I can use a repl and execute .yoc files.

## Specific goals

1. Easily rename some things.
  car -> first
  cdr -> rest
  define -> def
  lambda -> func/fn or something like that

2. easily add syntax
  - [items] -> (list items)
  - {} -> dictionary
  - some kind of splat operator like ...
  - make all definitions alike. So no special keyword for defining functions, rather just use `(def name (fn [x] (body)))` (same with macro)

3. Consider infix?
