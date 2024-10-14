# Yocto Lisp

Yocto lisp is a project that will serve as my own personal lisp dialect for solo projects.

I want it to be terse but descriptive.

For example, `first` or `head` insteal of `car`.
`fn` instead of `lambda`. In general I just want it to feel right and make sense to me and not use 'jargon' from old lisp.

I like lisp/scheme, but every dialect has some things I don't like.

I want a small amount of primitives where it's easy to built out your own syntax with.

My end goal is to have an executable that I can use a repl and execute .yoc files.

I need to decide how to implement, using what language etc.

> ## List of built ins
>
> Special forms:
> - quote
> - func
> - if
> - def
> - macro
>
> Functions:
> - append (like cons)
> - first
> - rest
> - same?
>
> Arithmetic
> - +, - , *, /, ^
>
> IO:
> - read
> - print
>
> Type checking:
> - list?
> - number?
> - string?
> - nil?
> - symbol?
