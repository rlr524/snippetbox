# main.go

## infoLog
> Use log.New() to create a logger for writing information messages. This takes 
> three parameters:
>> 1. The destination to write the logs to (os.Stdout) 
>> 2. A string prefix for the message, followed by a tab and 
>> 3. Flags to indicate what additional information to include (local date 
      > and time).
>>> Note that the date and time are joined using a bitwise OR operator (|).

## errorLog
> Create a logger for writing error messages in the same way, but use stderr 
> as the destination and use the log.Llongfile flag to include the relevant 
> path, file name and line number.

## srv
> Initialize a new http.Server struct. Set the Addr and Handler fields so that 
the server uses the same network address as in the Config struct and the routes 
> and set the ErrorLog field so that the server uses the custom errorLog logger 
> in the event of any problems.

## openDB()
> Wraps sql.Open() and returns a sql.DB connection pool for a given DSN


# snippets.go

## SnippetModel.Insert()
> This is the SQL statement we want to execute. Note, it's split into multiple 
> lines for readability using backticks instead of standard quotes.

### m.DB.Exec()
> Use the Exec() method on the embedded connection pool to execute the
> statement. The first parameter is the SQL statement, followed by the
> title, content, and expiry values for the placeholder (?) parameters. This
> method returns a sql.Result object, which contains some basic information
> about what happened when the statement was executed.

### result.LastInsertId()
> Use the LastInsertID() method on the result object to get the ID of
> the newly inserted record in the snippets table.


# handlers.go
## createSnippet()
### w.Header().Set()
> If the method is not POST use the w.WriteHeader() method to send a 405 
> status code and the w.Write() method to write a "Method Not Allowed" 
> response body. We then return from the function so that the subsequent code 
> is not executed. Use the Header().Set() method to provide additional 
> information to the user by setting a header as Allow:POST. Need to ensure 
> headers are set BEFORE WriteHeader() or Write() are called. Remember 
> the difference between WriteHeader() and Header().Set() is that WriteHeader() 
> only sets a status code and that can't be changed once set or set again where 
> Header().Set() sends the headers in the standard key:value format.

### app.clientError()
> Can combine the WriteHeader() and Write() into using the http.Error() method 
> which takes the ResponseWriter, an error message string, and the http code 
> to be returned. We can avoid having to do error handling on the Write() 
> method using this method.


# routes.go
## Route descriptions
| Method | Pattern         | Handler           | Action                       |
|--------|-----------------|-------------------|------------------------------|
| GET    | /               | home              | Display the home page        |
| GET    | /snippet/:id    | showSnippet       | Display a specific snippet   |
| GET    | /snippet/create | createSnippetForm | Display the new snippet form |
| POST   | /snippet/create | createSnippet     | Create a new snippet         |
| GET    | /static/        | http.Fileserver   | Serve a specific static fil  |

## Middleware TODO
> Need to break the middleware into dynamic and standard types. The dynamic will be
> app.session.Enable and app.requireAuthentication and all others will be standard.
> The standard applies to all routes and dynamic will will apply to the /, 
> /snippet/create, /snippet/:id, /user/signup, and /user/login routes.
