# Yocto Lisp

Yocto lisp is a project that will serve as my own personal lisp dialect for solo projects.

I want it to be terse but descriptive, and just feel 'right' to me.

## Example of current implementation

```
(defn (fib-tail n acc1 acc2)
  (if (<= n 1)
      acc2
          (fib-tail (- n 1) acc2 (+ acc1 acc2)))

  )

(defn (fibonacci n)
  (fib-tail n 0 1))

(print (fibonacci 50))

```
