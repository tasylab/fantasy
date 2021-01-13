package fantasy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fasthttp/router"
	"github.com/tasylab/fantasy/tcpserver"
	"github.com/valyala/fasthttp"
)

// Resolver is the interface that a user defined resolver need to implement.
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

// Fantasy holds a fantasy instance with the given backend as a resolver.
// It uses bigcache with default config as the default backend.
type Fantasy struct {
	resolver Resolver
}

// New returns a new fantasy instance.
func New(r Resolver) *Fantasy {
	return &Fantasy{
		resolver: newHook(r),
	}
}

// HTTPServer starts a fasthttp server to provide HTTP interface to interact with fantasy API.
func (f *Fantasy) HTTPServer(listenAddr string) error {
	r := router.New()
	r.GET("/{key}", func(ctx *fasthttp.RequestCtx) {
		key := ctx.UserValue("key").(string)
		val, err := f.resolver.Get(key)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		_, _ = ctx.Write(val)
	})
	r.POST("/{key}", func(ctx *fasthttp.RequestCtx) {
		key := ctx.UserValue("key").(string)
		if err := f.resolver.Set(key, ctx.Request.Body()); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		_, _ = ctx.WriteString(key)
	})
	r.DELETE("/{key}", func(ctx *fasthttp.RequestCtx) {
		key := ctx.UserValue("key").(string)
		if err := f.resolver.Del(key); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		_, _ = ctx.WriteString(key)
	})
	r.OPTIONS("/purge", func(ctx *fasthttp.RequestCtx) {
		if err := f.resolver.Purge(); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		_, _ = ctx.WriteString("OK")
	})
	return fasthttp.ListenAndServe(listenAddr, r.Handler)
}

// TCPServer creates a TCP server to interact with Fantasy API over plain TCP connection.
func (f *Fantasy) TCPServer(listenAddr string) error {
	server := tcpserver.New(listenAddr)
	server.OnNewClient(func(c *tcpserver.Client) {
		// new client connected
		// lets send some message
		_ = c.Send("__connected__\n\r")
	})
	server.OnNewMessage(func(c *tcpserver.Client, message string) {
		defer func() {
			_ = c.Send("\n\r")
		}()
		args, err := parseCommand(message)
		if err != nil {
			_ = c.Send(err.Error())
			return
		}
		resp, err := decideWhatToDo(f, args)
		if err != nil {
			_ = c.Send(err.Error())
			return
		}
		_ = c.SendBytes(resp)
	})
	server.OnClientConnectionClosed(func(c *tcpserver.Client, err error) {
		c.Close()
	})
	return server.Listen()
}

func decideWhatToDo(f *Fantasy, args []string) (response []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			response = make([]byte, 0)
			err = fmt.Errorf("%w", r)
		}
	}()
	requireMinimum(1, len(args))
	command := strings.ToUpper(args[0])
	switch command {
	case CommandGET:
		requireMinimum(2, len(args))
		key := args[1]
		val, err := f.resolver.Get(key)
		return val, err
	case CommandSET:
		requireMinimum(3, len(args))
		key := args[1]
		value := args[2]
		err := f.resolver.Set(key, []byte(value))
		return []byte(key), err
	case CommandDELETE:
		requireMinimum(2, len(args))
		key := args[1]
		err := f.resolver.Del(key)
		return []byte(key), err
	case CommandPURGE:
		err := f.resolver.Purge()
		return []byte("OK"), err
	case CommandLEN:
		val, err := f.resolver.Len()
		return []byte(fmt.Sprint(val)), err
	}
	return nil, errors.New("fantasy: unknown command passed")
}
