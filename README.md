# 🧬 WebSocket middleware for [Fiber](https://github.com/gofiber/fiber)

Based on [Fasthttp Fastws](https://github.com/fasthttp/fastws) for [Fiber](https://github.com/gofiber/fiber)

### Install

```
go get -u github.com/gofiber/fiber
go get -u github.com/fibergo/websocket
```

### Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/fibergo/websocket"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	app.Get("/", websocket.Upgrade(handler))

	app.Listen(3000)
}

// Websocket handler
func handler(conn *websocket.Conn) {
	fmt.Printf("Opened connection\n")
	conn.WriteString("Hello")
	var msg []byte
	var err error
	for {
		_, msg, err = conn.ReadMessage(msg[:0])
		if err != nil {
			if err != websocket.EOF {
				fmt.Fprintf(os.Stderr, "error reading message: %s\n", err)
			}
			break
		}
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing message: %s\n", err)
			break
		}
	}
	fmt.Printf("Closed connection\n")
}


  app.Listen(3000)
  // Access the websocket server: ws://localhost:3000/ws
}
```
