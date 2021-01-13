# fantasy
--
    import "."


## Usage

```go
const (
	CommandGET    = "GET"
	CommandSET    = "SET"
	CommandDELETE = "DELETE"
	CommandPURGE  = "PURGE"
	CommandLEN    = "LEN"
)
```
Commands.

#### type Fantasy

```go
type Fantasy struct {
}
```

Fantasy holds a fantasy instance with the given backend as a resolver. It uses
bigcache with default config as the default backend.

#### func  New

```go
func New(r Resolver) *Fantasy
```
New returns a new fantasy instance.

#### func (*Fantasy) HTTPServer

```go
func (f *Fantasy) HTTPServer(listenAddr string) error
```
HTTPServer starts a fasthttp server to provide HTTP interface to interact with
fantasy API.

#### func (*Fantasy) TCPServer

```go
func (f *Fantasy) TCPServer(listenAddr string) error
```
TCPServer creates a TCP server to interact with Fantasy API over plain TCP
connection.

#### type Resolver

```go
type Resolver interface {
	// Get reads entry for the key.
	Get(key string) ([]byte, error)
	// Contains checks existence of the key.
	Contains(key string) (bool, error)
	// Set sets the value for the key.
	Set(key string, value []byte) error
	// Del deletes the key and its associated value.
	Del(key string) error
	// Purge resets the cached keys.
	Purge() error
	// Len returns the total number of currently stored items.
	Len() (int, error)
}
```

Resolver is the interface that a user defined resolver need to implement.
