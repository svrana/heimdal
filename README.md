# Heimdal

## Movitations

In a repository with only one build target, setting up [entr](https://github.com/eradman/entr) to trigger the build
on file changes is how I typically manage my golang builds as I write code. (Pairing this with auto-deployments using
[tilt](https://tilt.dev) is sweet, btw). However I recently found myself in a monorepo with many build targets
and wished that I could only build the targets that required a rebuild.

## Setup

Place a .heimdal.toml in your golang project directory.

```
[[target]]
    name = "heimdal"
    pkg = "cmd/heimdal/main.go"
    command = "make"
```

The above example is the heimdal configuration for heimdal itself. In this case, heimdal starts parsing the golang code in cmd/heimdal/main.go.
If main.go or any file included by main.go changes, it runs the specified command, `make` in this case. The 'name' in the config is for
logging only.

Run `heimdal` from the directorry containing the .heimdal.toml file. Change main.go or one of the file that includes.


## Warning

This approach works fine for small projects (like this one, where it's not
really needed) but in the monorepo in which I was working it was slower than
just running `go build` on all the targets. Building all the targets was never
that painful (go build does a great job of not doing unnecessary work), so I
have yet to do the profiling required to figure out where to spend my efforts
optimizing it. I suspect that day will come, but until then..

# Shoulders

1. [dep-notify](github.com/mattmoor/dep-notify) does the heavy-lifting here
