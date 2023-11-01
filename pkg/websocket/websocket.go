package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	gorillaws "github.com/gorilla/websocket"

	"github.com/Marv963/CryptoTracker/app/pkg/workerpool"
)

type Connection struct {
	ws   *gorillaws.Conn
	send chan []byte // send-Queue
	mu   sync.Mutex  // Secure the write operation
}

func (c *Connection) writer() {
	for message := range c.send {
		c.mu.Lock()
		if err := c.ws.WriteMessage(gorillaws.TextMessage, message); err != nil {
			c.mu.Unlock()
			fmt.Println("write err:", err)
			break
		}
		c.mu.Unlock()
	}
	c.ws.Close()
}

type Server struct {
	conns         map[*Connection]bool
	subscriptions map[string]map[*Connection]bool
	connPairs     map[*gorillaws.Conn][]string
	pool          *workerpool.WorkerPool
	mu            sync.RWMutex // Mutex to protect conns and subscriptions
	activeConns   int64        // active connections
	activeSubs    int64        // active subscriptions
}

func NewServer() *Server {
	return &Server{
		conns:         make(map[*Connection]bool),
		subscriptions: make(map[string]map[*Connection]bool),
		connPairs:     make(map[*gorillaws.Conn][]string),
		pool:          workerpool.NewWorkerPool(4),
	}
}

func (s *Server) Start() {
	fmt.Println("Websocket started")
	go s.HandleWorkerErrors()
}

func (s *Server) addConn(ws *gorillaws.Conn) {
	atomic.AddInt64(&s.activeConns, 1)
	fmt.Printf("Connection added. Active connections: %d\n", atomic.LoadInt64(&s.activeConns))
	s.getOrCreateConnection(ws)
}

func (s *Server) getOrCreateConnection(ws *gorillaws.Conn) *Connection {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.unlockedGetOrCreateConnection(ws)
}

func (s *Server) unlockedGetOrCreateConnection(ws *gorillaws.Conn) *Connection {
	c := s.unlockedfindConn(ws)
	if c == nil {
		c = &Connection{ws: ws, send: make(chan []byte, 256)} // 256 is the size of the send queue
		s.conns[c] = true
		go c.writer() // Start the writer routine for this connection
	}
	return c
}

func (s *Server) findConn(ws *gorillaws.Conn) *Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.unlockedfindConn(ws)
}

func (s *Server) unlockedfindConn(ws *gorillaws.Conn) *Connection {
	for conn := range s.conns {
		if conn.ws == ws {
			return conn
		}
	}
	return nil
}

func (s *Server) removeConn(ws *gorillaws.Conn) {
	atomic.AddInt64(&s.activeConns, -1)
	fmt.Printf("Connection removed. Active connections: %d\n", atomic.LoadInt64(&s.activeConns))
	// fmt.Println("Locked mu in removeConn")
	s.mu.Lock()

	defer s.mu.Unlock()
	conn := s.unlockedfindConn(ws)
	delete(s.conns, conn)
}

func (s *Server) addSubscription(pair string, ws *gorillaws.Conn) {
	atomic.AddInt64(&s.activeSubs, 1)
	fmt.Printf("Subscription added for %s. Active subscriptions: %d\n", pair, atomic.LoadInt64(&s.activeSubs))
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subscriptions[pair] == nil {
		s.subscriptions[pair] = make(map[*Connection]bool)
	}
	c := s.unlockedGetOrCreateConnection(ws)
	s.subscriptions[pair][c] = true
	s.connPairs[ws] = append(s.connPairs[ws], pair)
	fmt.Printf("connPair added for %s. Active connpairs: %d\n", pair, len(s.connPairs[ws]))
}

// removeSubscription removes a specific websocket ws from the
// subscriptions for a specific currency (pair). If no more websockets are
// registered for this currency, the entry for the currency itself is also
// removed from the subscriptions map to save resources and keep the data
// structure clean.
func (s *Server) removeSubscription(c *Connection, pair string) {
	atomic.AddInt64(&s.activeSubs, -1)
	fmt.Printf("Subscription removed for %s. Active subscriptions: %d\n", pair, atomic.LoadInt64(&s.activeSubs))
	fmt.Println("Locked mu in removeSubscription")
	s.mu.Lock()
	defer s.mu.Unlock()
	defer fmt.Println("Unlocked mu in removeSubscription")
	if _, ok := s.subscriptions[pair]; ok {
		delete(s.subscriptions[pair], c)
		if len(s.subscriptions[pair]) == 0 {
			delete(s.subscriptions, pair)
		}
	}
	for ws, pairs := range s.connPairs {
		if c.ws == ws {
			// Create a new slice containing all pairs except the one to be removed.
			newPairs := []string{}
			for _, p := range pairs {
				if p != pair {
					newPairs = append(newPairs, p)
				}
			}
			// EN: If no pairs are left, delete the entry from the map.
			// Otherwise, assign the updated slice to the map.
			if len(newPairs) == 0 {
				delete(s.connPairs, ws)
				fmt.Printf("No more pairs for connection %v with %s. Removed from connPairs map.\n", ws.RemoteAddr().String(), pair)

			} else {
				s.connPairs[ws] = newPairs
				fmt.Printf("Updated pairs for connection %v in connPairs with %s map.\n", ws.RemoteAddr().String(), pair)

			}
			break
		}
	}
}

// SafeRemoveConn safely closes a websocket ws by first closing the
// websocket and then calling the removeSubscription function to remove
// the websocket from the subscriptions for a currency pair (pair). This
// function ensures that all resources associated with the websocket are
// properly released and that it is safely removed from the relevant data
// structures.
func (s *Server) SafeRemoveConn(c *Connection, pair string) {
	fmt.Printf("SafeRemoveConn called for pair: %s and connection: %v\n", pair, c.ws.RemoteAddr())
	c.ws.Close()
	s.removeSubscription(c, pair)
}

func (s *Server) Shutdown() {
	// Close all websocket connections
	s.mu.Lock() // Set write lock to ensure that no further changes are made to the connections
	for c := range s.conns {
		c.ws.Close()
	}
	s.mu.Unlock() // Release write lock

	// Stop the worker pool
	s.pool.Stop()
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	upgrader := gorillaws.Upgrader{
		// Allow all origins
		// TODO: Origin begrenzen
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade err:", err)
		return
	}
	fmt.Println("new incomming connection from client", ws.RemoteAddr())
	s.addConn(ws)
	defer s.removeConn(ws)
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *gorillaws.Conn) {
	for {
		// Read message from browser
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if gorillaws.IsCloseError(err, gorillaws.CloseNormalClosure, gorillaws.CloseAbnormalClosure, gorillaws.CloseGoingAway) {
				_, ok := s.connPairs[ws]
				if ok {
					conn := s.findConn(ws) // Find the *Connection associated with ws
					if conn != nil {
						// Iterate over all pairs and remove the connection.
						pairs, ok := s.connPairs[ws]
						if ok {
							for _, pair := range pairs {
								s.removeSubscription(conn, pair)
							}
						} else {
							s.removeConn(ws)
						}
					}
				}
				fmt.Println("Websocket closed normally")
				return
			}
			if gorillaws.IsCloseError(err, gorillaws.CloseGoingAway) {
				// TODO: Should somehting happen here?
			}
			fmt.Println("read err:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			fmt.Println("unmarshal err:", err)
			continue
		}

		method, ok := data["method"].(string)
		if !ok || method != "subscribe" {
			fmt.Printf("unexpected message: %s \n", msg)
			continue
		}

		params, ok := data["symbols"].([]interface{})
		if !ok {
			fmt.Println("unexpected params", data["symbols"])
			continue
		}

		for _, p := range params {
			pair, ok := p.(string)
			if !ok {
				fmt.Println("unexpected param type:", p)
				continue
			}
			s.addSubscription(pair, ws)
		}
	}
}

func (s *Server) Broadcast(pair string, data []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for c := range s.subscriptions[pair] {
		s.pool.Tasks <- workerpool.Task{
			Execute: func() error {
				select {
				case c.send <- data:
				default:
					// Report error if necessary
					// Possibly removing the connection and other clean-up
					return fmt.Errorf("could not send message to client: %v", c.ws.RemoteAddr())
				}
				return nil
			},
			Ws:   c.ws,
			Data: data,
			Pair: pair,
		}
	}
}

// New error-handling method/goroutine
func (s *Server) HandleWorkerErrors() {
	for task := range s.pool.Errors {
		fmt.Printf("Handling error for pair: %s err: %v\n", task.Pair, task.Error)
		// Here call SafeRemoveConn or other cleanup logic as needed
		conn := s.findConn(task.Ws) // Find the *Connection associated with task.Ws
		if conn != nil {
			s.SafeRemoveConn(conn, task.Pair) // Pass *Connection, not *gorillaws.Conn
		}
	}
}
