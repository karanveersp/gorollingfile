# gorollinglog

`gorollinglog` is a Go package that provides a rolling file logger. It allows logging to rolling files, where the active file always has the given name (ex. `myapp.log`), and backup files have an appended integer suffix `myapp_1.log`, `myapp_2.log`, and so on. The higher the integer, the older the logs. 

This package is useful for implementing loggers that automatically rotate log files when they reach a certain size.

There are so many more ways to configure a rolling logger and you can certainly use this implementation as a
starting point for your own custom logger!

## Installation

To use `gorollinglog` in your Go project, you can import it with:

```go
import "github.com/karanveersp/gorollinglog"
```

And then, you can create a new rolling file logger with the `NewRollingFile` function. It takes the following parameters:

- `dir`: The directory where the log files will be stored.
- `name`: The name of the log file.
- `mode`: The mode to open the log file. Use 'w' for overwrite and 'a' for append.
- `maxSizeMb`: The maximum size in megabytes that the log file can reach before it is rotated.
- `backups`: The number of backup files to keep.

## Usage


```go
package main

import (
    "log"

    "github.com/karanveersp/gorollinglog"
)

func main() {
    // Create a new RollingLogFile with the desired settings
    rf, err := gorollinglog.NewRollingFile("/path/to/logs", "logfile.log", 'a', 10.0, 5)
    if err != nil {
        log.Fatal("Error creating RollingLogFile:", err)
    }

    // Create a new logger using log.New and pass the RollingLogFile as the writer
    log := log.New(rf, "", log.LstdFlags)

    // Use the custom logger to log messages
    log.Println("Logging some data")
    log.Println("Logging more data")
}
```

In the above code, we create a new `RollingLogFile` with the desired settings using `NewRollingFile`.

Next, we create a new logger using `log.New` and pass the `RollingLogFile` as the writer. This logger can now be used to log messages, and the log output will be written to the rolling files handled by the `RollingLogFile`.

With this setup, your logger will automatically rotate log files when they reach the specified size and keep the specified number of backup files.


## Contributing

If you find any issues with `gorollinglog` or have suggestions for improvements, feel free to create a pull request or open an issue. Your contributions are welcome!

## License

`gorollinglog` is open-source software licensed under the [MIT License](https://github.com/karanveersp/gorollinglog/blob/master/LICENSE). You can use, modify, and distribute it freely according to the terms of the license.

Thank you for using `gorollinglog`! If you have any questions or need further assistance, please don't hesitate to reach out. Happy logging!
