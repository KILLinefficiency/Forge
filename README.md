# Forge

<br>
Forge is an automation tool written in Go. It's much similar to GNU Make.
<br>

The instructions that Forge executes are written in JSON format in a ``forgeMe.json`` or ``forgeMe`` file. The ``forgeMe.json`` or ``forgeMe`` file should exist in the directory where ``forge`` is run.

<br>
Forge is intended to be used for compiling projects and work around with files related to it.

<br>

## Getting Forge

Forge can be cloned from GitHub.
```
$ git clone https://www.github.com/KILLinefficiency/Forge.git
```

<br>

## Installing Forge

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

<br>

## Getting Started

Forge instructions are written in JSON format. These instructions are to be specified in a ``forgeMe.json`` or ``forgeMe`` file. Forge searches for the ``forgeMe.json`` first.

The ``forgeMe.json`` file consists of four keys:

* settings
* variables
* conditions
* heads

Each of these four keys start with ``!`` and have their own JSON objects as values.

This is the correct order of the keys:

```json
{
	"!settings": {},
	"!variables": {},
	"!conditions": {},
	"!heads": {}
}
```
<br>

### !heads

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

Even if a head contains a single shell command, it has to be specified in an array.

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
		"clean": ["^remove", "^alert", "echo Done"]
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

The **!conditions** JSON object has the name of the heads as it's keys and an array as it's value. You can have multiple conditional heads in **!conditions**. This array contains the addresses of the strictly existing files required by the head.

Even if there is only one required file, it has to be written in an array.

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

In this example, Forge will only be able to execute the **build** head if the files **main.c** and **my_lib.h** are present. You can pass relative or absolute addresses of the files in the array.

If the files do not exist, then the head(s) will disappear from the list of heads that you get by running ``forge --heads``.

Extended example:

```json
{
	"!conditions": {
		"build": ["main.c", "my_lib.h"],
		"clean": ["main"]
	},
	
	"!heads": {
		"build": ["gcc main.c -o main", "echo Compiled Successfully!"],
		"clean": ["rm main", "echo Done"]
	}
}
```

<br>

### !variables

Forge allows the users to run shell commands and capture the output (**stdout**) of the commands inside variables.

The **!variables** JSON object consists of variable name as the key and a string containing the shell command as the value. You can have multiple variables in **!variables**.

These variables can be used in the **!heads** JSON object by writing ``%`` infront of the variable names.

Consider a file ``data.txt`` containing just a word, ``Batman``.

**data.txt**
```
Batman
```

Using variables in ``forgeMe.json`` file.

```json
{
	"!variables": {
		"dir_name": "cat data.txt"
	},
	
	"!heads": {
		"make_dir": ["mkdir %dir_name"]
	}
}
```

On running ``forge make_dir`` a directory named "Batman" will be created. This ``forgeMe.json`` file will read the file ``data.txt`` and put it's content inside a variable called ``dir_name``. This variable is then accessed inside a head in the **!heads** object using ``%dir_name``.

If you want a variable to just have a string value instead of the **stdout** of a command, you can use ``echo``. Like,

```json
{
	"!variables": {
		"name": "echo Bruce Wayne"
	}
}
```

This will assign the value "Bruce Wayne" to the variable **name** like a normal variable.

<br>

### !settings

As the name suggests, the **!settings** object is used for specifying setting regarding the ``forgeMe.json`` file.

You can have the following settings in the **!settings** object:

* showSTDOUT
* showSTDERR
* delimiter
* default
* every

#### showSTDOUT

Forge, by default, shows the command, it's **stdout** or **stderr** on the terminal. If you don't want the **stdout** to appear, it can be disabled.

```
{
	"!settings": {
		"showSTDOUT": "false"
	},
	
	"!heads": {
		"show_text": ["echo I am Batman!"]
	}
}
```

Here, **showSTDOUT** can have two values, "**true**" or "**false**".

Note that the values are supposed to be passed as strings and not as booleans.

#### showSTDERR

You can supress the **stderr** given out by the heads.

```json
{
	"!settings": {
		"showSTDERR": "false"
	},
	
	"!heads": {
		"list_dir": ["ls -hello", "ls -alh"]
	}
}
```

#### delimiter

The default delimiter for the shell commands that separates the command-line utility and it's arguments is the space character (" ").

However, the user can change the delimiter. Like,

```json
{
	"!settings": {
		"delimiter": "-"
	},
	
	"!heads": {
		"folders": ["mkdir-Batman-Green Arrow"]
	}
}
```

This ``forgeMe.json`` file uses the hyphen character (``-``)  as the delimiter instead of using  a space character and makes two folders, "Batman" and "Green Arrow".

#### default

Forge executes one or more heads only when they are passed as command-line arguments to it.

However, Forge is also capable of executing a default head if no heads are passed as command-line arguments.

This default head can be specified by the user. Like,

```json
{
	"!settings": {
		"default": "build"
	},
	
	"!heads": {
		"build": ["gcc main.c -o main"],
		"clean": ["rm main"]
	}
}
```

On running Forge without any arguments, like,

```
$ forge
```

the **build** head will be executed automatically as no heads are specified in the command-line arguments.

#### every

You can also execute a head after every specific interval.

The **every** key has an array of two strings as it's value.

The first element of the array is the interval (in seconds) (specified as a string) you want a head to be run after. The second element is the name of the head which is needed to be run.

Like,

```json
{
	"!settings": {
		"every": ["10", "build"]
	},
	
	"!heads": {
		"build": ["gcc main.c -o main"]
	}
}
```

Running ``forge build`` will keep executing the **build** head after every 10 seconds.

<br>

## Using Pipes, Redirections and other shell related operations

Forge can't execute shell commands containing symbols like, ``|``, ``>``, ``>>``, ``;``, ``&&`` and ``||``.

This is because everything except for the command name is treated as arguments to that command and these symbols make invalid arguments.

Like,

**Wrong Example**

```json
{
	"!heads": {
		"proc": ["ps -A | grep atom"]
	}
}
```

The **proc** head in this ``forgeMe.json`` file uses a pipe in it's shell command. This is invalid and ``ps`` will throw an error about it.

However, there is a way to make this work.

This involves changing the delimiter and passing the shell command as a string to a shell inside the ``forgeMe.json`` file.

Running:
```
ps -A | grep atom
```

is equivalent to running:
```
sh -c "ps -A | grep atom"
```

This can be used in the ``forgeMe.json`` file.

Like,

**Correct Example**

```json
{
	"!settings": {
		"delimiter": ","
	},
	
	"!heads": {
		"proc": ["sh,-c,ps -A | grep atom"]
	}
}
```

This example uses comma (``,``) as delimiter and works as expected.

Note that there should be no space character before and after the comma (``,``).
