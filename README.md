# Forge

<br>
Forge is an automation tool written in Go. It's much similar to GNU Make.

The instructions that Forge executes are written in JSON format in a ``forgeMe.json`` or ``forgeMe`` file.

<br>
### Getting Forge

Forge can be cloned from GitHub.
```
$ git clone https://www.github.com/KILLinefficiency/Forge.git
```

<br>
### Installing Forge

Forge can be compiled easily with:
```
$ go build forge.go
```

However, you can also install Forge using the ``install.sh`` shell script. This script compiles Forge, copies it to ``~/.forge`` and adds it to the PATH environment variable using the ``~/.bashrc`` file.

```
$ ./install.sh
```

Restart your terminal after running this script.

Users running shells other than Bash can make changes to the ``install.sh`` script suitable to their shell config file.

### Getting Started

Forge instructions are written in JSON format. These instructions are to be specified in a ``forgeMe.json`` or ``forgeMe`` file. Forge searches for the ``forgeMe.json`` file if files with both names are present in the directory.

The ``forgeMe.json`` file consists of three keys:

* heads
* conditions
* variables
* settings

Each of these four keys start with ``!`` and have their own JSON objects as values

```json
{
	"!settings": {},
	"!variables": {},
	"!heads": {}
}
```
<br>
#### !heads

A head is a collection of shell commands which are run one by one. The ``!heads`` JSON object can contain multiple heads with an array as their values.

```json
{
	"!heads": {
		"build": ["gcc main.c -o main", "echo Compiled Successfully!"],
		"remove": ["rm main"],
		"alert": ["echo main removed."]
	}
}
```

Even if a head contains a single shell command, it has to be specified in am array.

These heads can be run by calling Forge and specifying the heads as the command-line arguments. Like,

```
$ forge build
```
This will run all the commands inside the **build** head. You can also pass multiple existing heads as the command-line arguments. Like,

```
$ forge build remove
```
This will run the **build** head first and then the **remove** head.

You can also chain multiple heads inside a head. Other heads can be referenced inside a head by writing ``^`` in front of the head(s) that needs to be referenced in the current head.

```json
{
	"!heads": {
		"build": ["gcc main.c -o main", "echo Compiled Successfully!"],
		"remove": ["rm main"],
		"alert": ["echo main removed."],
		"clean": ["^remove", "^alert", "echo Done."]
	}
}
```

Here, the head called "**clean**" calls two heads, **remove** and **alert** in the specified order and then runs a shell command.
A head which contains one or more referenced heads can also contain a shell command.

You can see a list of all available heads with:

```
$ forge --heads
```

<br>
### !conditions

The **!conditions** key has a JSON object as it's value.

**!conditions** allows a head to be executed only if specified files are present on the system. If all of the specified files are not present then the head will just be deleted.

The **!conditions** JSON object has the name of the heads a it's keys and an array as it's value. You can have multiple conditional heads in **!conditions**. This array contains the addresses of the strictly existing files required by the head as it's values.

```json
{
	"!conditions": {
		"build": ["main.c", "my_lib.h"]
	},
	
	"!heads": {
		"build": ["gcc main.c -o main", "echo Compiled Successfully!"]
	}
}
```

In this example, the Forge will only be able to execute the **build** head if the files **main.c** and **my_lib.h** are present. You can pass relative or absoulte address of the files in the array.

If the files do not exist, then the head(s) will disappear from the list of heads that you get by running ``forge --heads``.