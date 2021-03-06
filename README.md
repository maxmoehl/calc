# Overview

`calc` is a simple cli to calculate mathematical expressions written in pure go. It supports
basic operations  like `+`, `-`, `*`, `/` and parentheses `(` and `)` out of the box.
Further functionality can be added by macros. The `macros` directory contains a macro for
calculating the square root and power.

## Constraints

Since the [plugin](https://pkg.go.dev/plugin) package only supports macOS, FreeBSD and Linux
with cgo enabled, the plugin loading mechanism is only enabled for those platforms. This is
achieved by setting the build directive in `plugins.go`. Any code in that file will only be
included if the target platform is listed in this directive.

# Installation

Currently, no pre-built executables are available but building it is easy since there are zero
dependencies (although that wouldn't be a problem anyways). The application was built for go
1.16, and it is recommended to have at least this version. After installing go form the
official website, using the application only takes one more command:

```
go install github.com/maxmoehl/calc/cli@latest
```

Depending on your setup you might have to add the `$GOBIN` directory to your path variable.
By default, it is `$HOME/go/bin`. If you want to use macros see section [Macros](#macros)

# Configuration

The plugin in home can be overridden by setting the environment variable `CALC_PLUGIN_DIR`.
If the variable is not present, the default directory `$HOME/.calc` is used, if the
environment variable is empty no plugins will be loaded.

To get debug information set `DEBUG=1` as an environment variable or directly pass it to
the executable:
```
$ DEBUG=1 calc "1+1"
```

# Usage

It can either be used to calculate a single expression:
```
$ calc "1+1"
2
```
or launched in interactive mode to execute multiple expressions:
```
$ calc -interactive
> 1+1
2
> pow{2,2}
4
> _
```
Type `exit` or press CTRL+C to exit

If executed without any arguments, a little help section gets printed:
```
$ calc
Usage:
  either execute a single calculation:
    calc <mathematical expression>
  or start the interactive mode:
    calc -interactive

Loaded macros:
  sqrt, pow
```

Plans for the **future**: Allow recent results to be reused in calculations. This would allow
for more flexibility in calculation, i.e.:
```
$ calc -interactive
> 1
1
> $0 + 1
2
> $0 + 1
3
> _
```

# Syntax

## Extended Backus???Naur form

_I'm not 100% sure if this is correct_ ??\_(???)_/??

```
digit      = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" | "0" ;
letter     = "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j" | "k" |
             "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" |
             "w" | "x" | "y" | "z" ;

number     = { digit }, [ ".", [ { digit } ] ] ;
identifier = { letter } ;

plus_minus = "+" | "-" ;
mul_div    = "*" | "/" ;

parameter  = expression, { ",", expression };
macro      = identifier, "{", [ parameter, ] "}" ;
operator   = plus_minus | mul_div ;
operand    = number | macro ;


expression = operand | [ plus_minus, ] [ "(", ] expression, [ { operator, expression, } ] [ ")" ] ;
             
```

# Macros

_And make sure your system is [supported](#constraints)_

The application loads all plugins with the file ending `*.so `inside `$HOME/.calc`. For each
plugin the `Index` is checked and any listed macros are registered.
To enable the default macros do the following:
1. Build the macros in plugin mode: `go build -buildmode=plugin ./macros`
2. Create the directory `$HOME/.calc` if it does not exist: `mkdir $HOME/.calc`
3. Copy the built file to the newly created directory: `mv macros.so $HOME/.calc/macros.so`

## Invoking macros

Plugins can be invoked by their identifier and braces containing the parameters delimited by
commas:
```
identifier{parameter, ..., parameter}
```

For example calculating two squared:
```
$ calc "pow{2, 2}"
```

## Building your own macros (plugins)

The Plugin has a single requirement: An exported variable named `Index` of type `types.Index`.
This index maps the identifier of a macro (which can be used in mathematical expressions)
to a function name that can be used to create a representation of the macro in memory.
An example says more than a thousand words, take a look at the files in `macros/` for
a working example.

The Macro itself needs the following:
1. A function to create a new instance of it
2. A struct that holds all necessary information

The function to create a new macro needs to implement the `types.NewMacro` interface which
has the following signature:

```go
type NewMacro func(parameter []Node) (Macro, error)
```

Inside this function basic validation should be done to ensure correct number of arguments.
The returned macro needs to implement the `types.Macro` interface which is defined as follows:

```go
type Macro interface {
	Eval() (float64, error)
}
```

The `Eval` function is used to evaluate the macro. Since the parameters are `types.Node`s
you first have to evaluate all parameters and use the results to do your own calculation. If
any errors occur while evaluating the parameters it is recommended to return `math.NaN()` and
the encountered error without modifying the error or returning your own.

After you've written your plugin ensure that the package name is `main` and try to build it
using `buildmode=plugin`. Copy the resulting `*.so` file to `$HOME/.calc` and run the `calc`
cli to test if it works.
