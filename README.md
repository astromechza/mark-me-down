# mark-me-down

```
$ ./mark-me-down --help
mark-me-down is a simple binary for rendering Github Flavoured Markdown content.

Run it with a single argument (a path to a file) that will be formatted into html
and served whenever a request hits the local server.

The --listen-port field is provided in order to specify a particular port.

  -listen-port int
        Server the markdown on this port (default 80)
```

## Example

```
$ ./mark-me-down --listen-port 8080 README.md
Listening on localhost:8080...
Attempting to open a browser window to the address..
```

## How to build

```
$ ./make
Building official darwin amd64 binary
Output Folder build/darwin_amd64
github.com/AstromechZA/mark-me-down
Done
-rwxr-xr-x  1 username  staff  8603356 May 24 21:36 build/darwin_amd64/mark-me-down
build/darwin_amd64/mark-me-down: Mach-O 64-bit executable x86_64

Building official linux amd64 binary
Output Folder build/linux_amd64
github.com/AstromechZA/mark-me-down
Done
-rwxr-xr-x  1 username  staff  8664056 May 24 21:36 build/linux_amd64/mark-me-down
build/linux_amd64/mark-me-down: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped
```

This will build Linux and OSX binaries in 64bit. If you need to build for other
operating systems, modify the `./make` file.
