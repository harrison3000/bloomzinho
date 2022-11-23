# Bloomzinho
A very simple bloom filter

* Small and simple code
* Concurrency-safe (for lookups, insertions needs to be synchronized)
* No allocations on lookups (unless you use more than 8 hashes)
* No external dependencies
* Kinda fast
* Made in Brazil