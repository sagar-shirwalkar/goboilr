# GoBoilr

![Go Version](https://img.shields.io/github/go-mod/go-version/sagar-shirwalkar/goboilr?style=flat-square&logo=go&color=%2300ADD8)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagar-shirwalkar/goboilr?style=flat-square)](https://goreportcard.com/report/github.com/sagar-shirwalkar/goboilr)
[![codecov](https://img.shields.io/codecov/c/github/sagar-shirwalkar/goboilr/main.svg?style=flat-square&label=Coverage)](https://codecov.io/gh/sagar-shirwalkar/goboilr)
[![GitHub tag](https://img.shields.io/github/tag/sagar-shirwalkar/goboilr?include_prereleases=&sort=semver&color=fbd12b&style=flat-square&label=Tag)](https://github.com/sagar-shirwalkar/goboilr/releases/)
![GitHub code size](https://img.shields.io/github/languages/code-size/sagar-shirwalkar/goboilr?style=flat-square&label=Code%20Size)
[![Issues](https://img.shields.io/github/issues/sagar-shirwalkar/goboilr?style=flat-square&color=21ca26&label=Issues)](https://github.com/sagar-shirwalkar/goboilr/issues)
[![License](https://img.shields.io/github/license/sagar-shirwalkar/goboilr?style=flat-square&color=bd2bfb&label=License)](#license)

**GoBoilr** is a lightweight, zero-dependency code generator for Go. It eliminates boilerplate by automatically generating type-safe Getters, Setters, Constructors and Builders for your structs using standard Go tags and comments.

It is designed to work seamlessly with `go generate` and produces high-performance code with no runtime reflection.

Can also be combined with **[GoLoom](https://github.com/sagar-shirwalkar/goloom)** to quickly generate rich domain models directly from JSON schemas:

[![sagar-shirwalkar - goboilr](https://img.shields.io/static/v1?label=sagar-shirwalkar&message=GoLoom&color=fbd12b&logo=github&style=flat-square)](https://github.com/sagar-shirwalkar/goloom "Go to GitHub repo")

## Features

* **Zero Dependencies:** Built entirely with the standard library (`go/parser`, `go/ast`).
* **Getters & Setters:** Generates standard accessor methods based on struct tags.
* **Validation Hooks:** Auto-wires setters to your custom validation logic.
* **Constructors:** Opt-in generation of "All-Args" constructors.
* **Builders:** Builder generation - super handy if "All-Args" constructors have too many arguments.
* **Safe Integration:** Constructor code is generated in a commented-out block for easy copy-pasting, preventing redeclaration errors.
* **Complex Type Support:** Handles slices, maps, pointers, and external package imports (like `time.Time`, `json.RawMessage`) automatically.
* **Standard Tooling:** Works out-of-the-box with `go generate`.

## Installation

### Option 1: Install as a Tool (Convenient personal use)

Install the binary globally to use it via the command line or `go generate`.

```bash
go install github.com/sagar-shirwalkar/goboilr@latest
```

### Option 2: Add as a dependency (CI-friendly, recommended for teams)

If you want to track the version in your go.mod file:

```bash
go get github.com/sagar-shirwalkar/goboilr
```

## How to use

### 1. Configure your structs

In any Go file (let's say, /models/user.go), add the **go:generate** directive at the top.

Add **//gen:new** right above your structs to request an all-args constructor.

Use the gen struct tag to define accessors (getters, setters with or without validation). Example:

```go
package models

import (
    j "encoding/json"
    "fmt"
    "time"
)

/* Trigger for directing the generator to this file */

//go:generate goboilr

// gen:new
// gen:builder
type User struct {
    // Embedded field (No tag, but affects constructor/builder)
    Base

    Name string `gen:"get,set"`
    Age  int    `gen:"get,set,val"`

    // Complex types
    Tags     []string     `gen:"get"`
    Metadata j.RawMessage `gen:"get"`
    Ptr      *int         `gen:"set"`
}

type Base struct {
    ID string `gen:"get"`
}
```

The generator works with complex types and imports, and if you use the `val` tag, validation will be carried out in the generated setter.

The validation logic must be written by you, but your validation `func` must follow a naming convention for the validation hook on the setter to work.

Let's take the **User's** `age` as an example:

```go
// The generator will create SetAge that calls this.
func (u *User) validateAge(v int) bool {
    return v >= u.age // Age being set must be greater than current age
}
```

### 2. Run the Generator

Run the standard Go command from your project root, targering /models/user.go:

```bash
go generate ./...
```

### 3. The output

GoBoilr will create two files in the same directory as your source file:

1. **user_accessors.go** : Contains the compiled, ready-to-use methods.

    ```go
    func (x *User) Age() int {
        return x.age
    }


    func (x *User) SetAge(v int) bool {
        if !x.validateAge(v) {
            return false
        }
        x.age = v
        return true
    }
    ```

2. **user_constructors.go** : Contains the constructor logic wrapped in comments (to avoid potential conflicts), which you can simply copy into your **user.go** file if you need it. Also contains the generated builders.

    ```go
    /*
    func NewUser(base Base, name string, age int, ...) *User {
        return &User{...}
    }
    */

    // -----------------------------------------------------------------------------
    // User Builder
    // -----------------------------------------------------------------------------

    type UserBuilder struct {
        target *User
    }

    func NewUserBuilder() *UserBuilder {
        return &UserBuilder{
            target: &User{},
        }
    }

    func (b *UserBuilder) Build() *User {
        return b.target
    }

    func (b *UserBuilder) Base(v Base) *UserBuilder {
        b.target.Base = v
        return b
    }

    func (b *UserBuilder) Name(v string) *UserBuilder {
        b.target.Name = v
        return b
    }

    func (b *UserBuilder) Age(v int) *UserBuilder {
        b.target.Age = v
        return b
    }

    /* And other builder funcs */
    ```

## Documentation

[![view - Documentation](https://img.shields.io/badge/view-Documentation-fbd12b?style=for-the-badge)](/docs/ "Go to project documentation")

## License

Released under [MIT](/LICENSE) by [@sagar-shirwalkar](https://github.com/sagar-shirwalkar)
