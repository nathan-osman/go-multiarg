## go-multiarg

[![Build Status](https://travis-ci.org/nathan-osman/go-multiarg.svg?branch=master)](https://travis-ci.org/nathan-osman/go-multiarg)
[![GoDoc](https://godoc.org/github.com/nathan-osman/go-multiarg?status.svg)](https://godoc.org/github.com/nathan-osman/go-multiarg)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

Many Go applications that I write use command-line arguments to control application behavior. The problem with command-line arguments is that they are visible to other users, eliminating the possibility of using them for sensitive data.

Environment variables and configuration files present two solutions to this problem. Unfortunately, this means writing a lot of boilerplate code to check for arguments in different places. go-multiarg exists to simplify the process.

### Usage

To use go-multiarg, pass a `struct` to the `multiarg.Load()` function with the default values set:

    import (
        "os"
        "github.com/nathan-osman/go-multiarg"
    )

    const configFilename = "/etc/myapp/config.json"

    type Config struct {
        NumTries int    `multiarg:"number of tries"`
        LoginURL string `multiarg:"URL for login"`
    }

    func main() {
        config := Config{
            NumTries: 3,
            LoginURL: "https://example.com/login",
        }
        if ok := multiarg.Load(&config, &multiarg.Config{}); !ok {
            os.Exit(1)
        }

        // NumTries can be specified in three ways:
        // - via the configuration file using "num_tries"
        // - via CLI using "--num-tries"
        // - via env. variable NUM_TRIES
    }

