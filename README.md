# GoBoilr

![License](https://img.shields.io/github/license/sagar-shirwalkar/goboilr?style=flat-square&color=dark-green&label=License)
![Go Version](https://img.shields.io/github/go-mod/go-version/sagar-shirwalkar/goboilr?style=flat-square&logo=go&color=%2300ADD8)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagar-shirwalkar/goboilr?style=flat-square)](https://goreportcard.com/report/github.com/sagar-shirwalkar/goboilr)
![Issues](https://img.shields.io/github/issues/sagar-shirwalkar/goboilr?style=flat-square&label=Issues)
![GitHub repo size](https://img.shields.io/github/repo-size/sagar-shirwalkar/goboilr?style=flat-square&label=Size)


**GoBoilr** is a lightweight, zero-dependency code generator for Go. It eliminates boilerplate by automatically generating type-safe Getters, Setters, and Constructors for your structs using standard Go tags and comments.

It is designed to work seamlessly with `go generate` and produces high-performance code with no runtime reflection.

## Features

* **Zero Dependencies:** Built entirely with the standard library (`go/parser`, `go/ast`).
* **Getters & Setters:** Generates standard accessor methods based on struct tags.
* **Validation Hooks:** Auto-wires setters to your custom validation logic.
* **Constructors:** Opt-in generation of "All-Args" constructors.
* **Safe Integration:** Constructor code is generated in a commented-out block for easy copy-pasting, preventing redeclaration errors.
* **Complex Type Support:** Handles slices, maps, pointers, and external package imports (e.g., `time.Time`, `json.RawMessage`) automatically.
* **Standard Tooling:** Works out-of-the-box with `go generate`.

## Installation

### Option 1: Install as a Tool (Recommended)

Install the binary globally to use it via the command line or `go generate`.

```bash
go install github.com/sagar-shirwalkar/goboilr@latest
```

### Option 2: Add as a dependency

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

//gen:new
type User struct {
    name      string       `gen:"get,set"`
    age       int          `gen:"get,set,val"`
    weight    int          `gen:"get,set,val"`
    birthDate time.Time    `gen:"get,set"`
    config    j.RawMessage `gen:"get,set"`
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

2. **user_constructors.go** : Contains the constructor logic wrapped in comments (to avoid potential conflicts). You can simply copy this func into your **user.go** file if you need it.

    ```go
    /*
    func NewUser(id string, name string, ...) *User {
        return &User{...}
    }
    */
    ```

## License

Distributed under the MIT License. See `LICENSE` for more information.
