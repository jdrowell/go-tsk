# go-tsk
Go bindings for The Sleuth Kit forensic library

# What is Working

* disk images (open, properties)
* volume systems (open, properties)
* partitions (open, properties)
* filesystems (open, properties, walk)
* files (copy)
* error handling

# How to Use

Check the busytsk folder for examples.

# Dependencies

These bindings are up to date with version 4.12.0 of [sleuthkit](https://github.com/sleuthkit/sleuthkit).
You must have sleuthkit-dev installed to be able to link your binary. Also, make sure **NOT** to have something
like <code>CGO_ENABLED=0</code> in your environment.

# Building

You can build the demo utility by entering the busytsk folder and running <code>go build</code>.

# Documentation

I'm slowly adding documentation in <code>go doc</code> format. Just use that command to explore it.
