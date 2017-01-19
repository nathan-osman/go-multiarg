## go-multiarg

[![Build Status](https://travis-ci.org/nathan-osman/go-multiarg.svg?branch=master)](https://travis-ci.org/nathan-osman/go-multiarg)
[![GoDoc](https://godoc.org/github.com/nathan-osman/go-multiarg?status.svg)](https://godoc.org/github.com/nathan-osman/go-multiarg)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

Many Go applications that I write use command-line arguments to control application behavior. The problem with command-line arguments is that they are visible to other users, eliminating the possibility of using them for sensitive data.

Environment variables and configuration files present two solutions to this problem. Unfortunately, this means writing a lot of boilerplate code to check for arguments in different places. go-multiarg exists to simplify the process.

### Example

To use go-multiarg, pass a `struct` to the `multiarg.Load()` function with the default values set:

    import (
        "os"
        "github.com/nathan-osman/go-multiarg"
    )

    type Config struct {
        NumTries int    `multiarg:"number of tries"`
        LoginURL string `multiarg:"URL for login"`
    }

    func main() {
        config := Config{
            NumTries: 3,
            LoginURL: "https://example.com/login",
        }
        if ok, _ := multiarg.Load(&config, &multiarg.Config{
            JSONFilenames: []string{"/etc/myapp/config.json"},
        }); !ok {
            os.Exit(1)
        }

        // Do stuff
    }

In the example above, `NumTries` could be set to `3` in any of these three ways:

- using the JSON configuration file: `{"num_tries": 3}`
- using an environment variable: `NUM_TRIES=3`
- using a CLI argument: `--num-tries 3`
