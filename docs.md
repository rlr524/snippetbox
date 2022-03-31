#main.go

##infoLog
> Use log.New() to create a logger for writing information messages. This takes three parameters:
The destination to write the logs to (os.Stdout) 2. A string prefix for the message,
followed by a tab and 3. Flags to indicate what additional information to include (local date and time).
Note that the date and time are joined using a bitwise OR operator (|).

##errorLog
> Create a logger for writing error messages in the same way, but use stderr as the destination and use the 
log.Llongfile flag to include the relevant path, file name and line number.

##srv
> Initialize a new http.Server struct. Set the Addr and Handler fields so that 
the server uses the same network address as in the Config struct and the routes and set the ErrorLog field so that the server uses the custom errorLog logger in the event of any problems.