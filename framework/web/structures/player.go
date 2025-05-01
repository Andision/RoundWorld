package structures

import (
	"bytes"
	"github.com/Andision/RoundWorld/framework/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a Message to the peer.
	playerWriteWait = 10 * time.Second

	// Time allowed to read the next pong Message from the peer.
	playerPongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than playerPongWait.
	playerPingPeriod = (playerPongWait * 9) / 10

	// Maximum Message size allowed from peer.
	playerMaxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Player is a middleman between the websocket connection and the table.
type Player interface {
	readPump()
	writePump()

	Send(message []byte) bool
	Close()
}

// NewPlayer handles websocket requests from the peer.
func NewPlayer(user models.User, table Table, w http.ResponseWriter, r *http.Request) Player {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	player := &PlayerImpl{
		user:  user,
		table: table,
		conn:  conn,
		send:  make(chan []byte, 256),
	}

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go player.writePump()
	go player.readPump()

	return player
}

type PlayerImpl struct {
	user  models.User
	table Table

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *PlayerImpl) readPump() {
	defer func() {
		c.table.SendUnregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(playerMaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(playerPongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(playerPongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.table.SendMessage(TableMessage{
			Sender:  c,
			Message: bytes.TrimSpace(bytes.Replace(message, newline, space, -1)),
		})
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *PlayerImpl) writePump() {
	ticker := time.NewTicker(playerPingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			log.Printf("Received Message from send channel: %v", message)
			c.conn.SetWriteDeadline(time.Now().Add(playerWriteWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("writePump: close")
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket Message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {

				log.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(playerWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		default:
			//log.Println("Sender.send channel is full")
		}
	}
}

func (c *PlayerImpl) Send(message []byte) bool {
	select {
	case c.send <- message:
		return true
	default:
		return false
	}
}

func (c *PlayerImpl) Close() {
	close(c.send)
}
