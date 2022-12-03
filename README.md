# Bloomzinho
A very simple bloom filter

* Small and simple code
* Concurrency-safe lookups (insertions needs to be synchronized)
* No allocations on lookups or insertions (unless you use more than 8 hashes)
* No external dependencies, only standard library
* Kinda fast
* Made in Brazil

## TODO

* Better documentation
* "Contains" method