* handler:
> responsible for executing application logic and for writing HTTP response headers and bodies.

* router: (servemux)
> This stores a mapping between the URL patterns for your application and the corresponding handlers.

servemux supports two different types of URL patterns: `fixed paths` and `subtree paths`.
    * fixed paths: don't end with a trailing slash.
    * subtree paths: do end with a trailing slash.

