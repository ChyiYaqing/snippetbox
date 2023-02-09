package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/chyiyaqing/snippetbox/internal"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
}

// Define an application struct to hold the application wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	// Define a new command-line flag for the MySQL DSN string.
	// username:password@protocol(ipaddress:port)/dbname?param=value
	dsn := flag.String("dsn", "web:macintosh@tcp(HuaWei:13306)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator .
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the releveant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also define a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new instance of our application struct, containing the
	// dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  internal.PrometheusMiddleware(app.routes()),
	}

	infoLog.Printf("Staring server on %s", cfg.addr)
	err = srv.ListenAndServe()
	// log.Fatal() function will also call os.Exit(1) after writing the message.
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Actual connections to the database are estabilished lazily. as and when
	// needed for the first time. So to verify that everything is set up correctly
	// we need to use the db.Ping() method to create a connection and check for
	// any errors.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
