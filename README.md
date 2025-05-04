# Cache Manager
OBS: THIS ENTIRE REPOSITORY WAS CREATED BY AI, I asked an LLM to build a cache manager for a project of mine that I was making to learn golang. And since I was learning, I was unable of creating that, you can use this repository, but it may not work as expected. Also, it stores data in memory, which in this case are handled by GC, which can give to your project some performance penalties. I would not use this repository now-a-days, but if you want to, feel free to.


This project implement a simple system of caching in Go, with support to item expiration based on it's life time (`life_time`). The cache being managed in memory, using the library `crypto/sha256` to generate unique keys for each stored items.

## Project structure


## Functions

- **Set**: Stores a item in ache with a key.
- **Get**: Get the value which a key points to.
- **Remove**: Using a key, deletes the value it points to
- **Exist**: Verify if the key point to any value.
- **Automatic clean up**: Each item have a lifetime which is decrease each second, when it reaches zero the value is deleted.
