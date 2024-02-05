package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "io"
  "net/http"
)

type Server struct {
  conns map[*websocket.Conn]bool
}

func NewServer () Server {
  return Server {
    conns: make(map[*websocket.Conn]bool),
  }
}


func (s *Server) handleWs(ws *websocket.Conn) {
  fmt.Println("New connection incoming from: ", ws.RemoteAddr())

  s.conns[ws] = true

  s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
  buf := make([]byte, 1024)
  for {
    n, err := ws.Read(buf)
    if err != nil {
      if err == io.EOF {
        break
      }
      fmt.Println("Error reading from websocket: ", err)
      continue
    }
    msg := buf[:n]
    s.broadcast(msg)
  }
}


func (s *Server) broadcast(b []byte) {
  for ws := range s.conns {
    go func(ws *websocket.Conn) {
      if _, err := ws.Write(b); err != nil {
        fmt.Println("Error writing to websocket: ", err)
      }
    }(ws)
  }
}

func main () {
  server := NewServer()
  http.Handle("/ws", websocket.Handler(server.handleWs))
  http.ListenAndServe(":3000", nil)

}
