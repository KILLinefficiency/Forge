# Forge


Forge is an automation tool written in Go. It's much similar to GNU Make.

The instructions that Forge executes are written in JSON format in a ``forgeMe.json`` or ``forgeMe`` file.

### Getting Forge

Forge can be cloned from GitHub.
```
$ git clone https://www.github.com/KILLinefficiency/Forge.git
```
### Installing Forge

Forge can be compiled easily with:
```
$ go build forge.go
```

However, you can also install Forge using the ``install.sh`` shell script. This script compiles Forges, copies it to ``~/.forge`` and adds it to the PATH environment variable.

```
$ ./install.sh
```