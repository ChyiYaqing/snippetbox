# [Let's Go!](https://lets-go.alexedwards.net/)

> Learn to Build Professional Web Applications with Go

Let's Go is a clear, concise and easy-to-follow guide to web development with Go.

It packs in everything you need to know about best practices, project structure
and practical code patterns -- without skimping on important details and explanations.


## Chapter 6. Middleware

* Panic recovery

In a simple Go application, when your code panics it will result in the application being terminated straight away.

But web application is a bit more sophisticated. Go's HTTP server assumes that the effect of any panic is isolated to the goroutine serving the active HTTP request (remember, every request is handled in it's own goroutine).


## Chapter 10. Security improvemenets

Run the generate_cert.go tool.
```
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
2023/02/18 09:47:51 wrote cert.pem
2023/02/18 09:47:51 wrote key.pem
```