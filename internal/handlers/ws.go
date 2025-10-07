package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"LinkUp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type wsIncoming struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type Presence struct {
	mu     sync.Mutex
	online map[uint]time.Time
}
// Auto-generated swagger comments for NewPresence
// @Summary Auto-generated summary for NewPresence
// @Description Auto-generated description for NewPresence — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func NewPresence() *Presence { return &Presence{online: map[uint]time.Time{}} }
func (p *Presence) Online(uid uint) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.online[uid] = time.Now()
}
func (p *Presence) Offline(uid uint) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.online, uid)
}
// Auto-generated swagger comments for LastSeen
// @Summary Auto-generated summary for LastSeen
// @Description Auto-generated description for LastSeen — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)
func (p *Presence) LastSeen(uid uint) (time.Time, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.online[uid]
	return t, ok
}
func (p *Presence) IsOnline(uid uint) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.online[uid]
	return ok
}

// Возвращает короче все короче статусы короче  для короче  ws короче клиентов
// Auto-generated swagger comments for AllStatuses
// @Summary Auto-generated summary for AllStatuses
// @Description Auto-generated description for AllStatuses — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)
func (p *Presence) AllStatuses() map[uint]map[string]interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	status := map[uint]map[string]interface{}{}
	for uid, t := range p.online {
		status[uid] = map[string]interface{}{
			"online":   true,
			"lastSeen": t,
		}
	}
	return status
}

// ------------------- RoomHubs -------------------

type RoomHubs struct {
	hubs map[uint]*Hub
}
// Auto-generated swagger comments for NewRoomHubs
// @Summary Auto-generated summary for NewRoomHubs
// @Description Auto-generated description for NewRoomHubs — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func NewRoomHubs() *RoomHubs { return &RoomHubs{hubs: map[uint]*Hub{}} }
func (r *RoomHubs) hub(roomID uint) *Hub {
	if h, ok := r.hubs[roomID]; ok {
		return h
	}
	h := NewHub(roomID)
	r.hubs[roomID] = h
	go h.Run()
	return h
}
// Auto-generated swagger comments for Emit
// @Summary Auto-generated summary for Emit
// @Description Auto-generated description for Emit — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)
func (r *RoomHubs) Emit(roomID uint, ev Event) { r.hub(roomID).Broadcast(ev) }

// ------------------- Hub -------------------

type Hub struct {
	roomID     uint
	register   chan *Client
	unregister chan *Client
	broadcast  chan Event
	clients    map[*Client]bool
}
// Auto-generated swagger comments for NewHub
// @Summary Auto-generated summary for NewHub
// @Description Auto-generated description for NewHub — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func NewHub(roomID uint) *Hub {
	return &Hub{
		roomID:     roomID,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Event, 64),
		clients:    map[*Client]bool{},
	}
}
// Auto-generated swagger comments for Run
// @Summary Auto-generated summary for Run
// @Description Auto-generated description for Run — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			h.Broadcast(Event{Type: "presence_join", Payload: gin.H{"userId": c.userID}})
			h.broadcastPresenceUpdate()
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.Broadcast(Event{Type: "presence_leave", Payload: gin.H{"userId": c.userID}})
			h.broadcastPresenceUpdate()
		case ev := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- ev:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}
}

// Отправка статусов всех пользователей
// Auto-generated swagger comments for broadcastPresenceUpdate
// @Summary Auto-generated summary for broadcastPresenceUpdate
// @Description Auto-generated description for broadcastPresenceUpdate — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)
func (h *Hub) broadcastPresenceUpdate() {
	status := h.clientsStatus()
	ev := Event{Type: "presence_update", Payload: status}
	for c := range h.clients {
		select {
		case c.send <- ev:
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}

// Получаем статусы пользователей из Presence
// Auto-generated swagger comments for clientsStatus
// @Summary Auto-generated summary for clientsStatus
// @Description Auto-generated description for clientsStatus — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)
func (h *Hub) clientsStatus() map[uint]map[string]interface{} {
	status := map[uint]map[string]interface{}{}
	for c := range h.clients {
		status[c.userID] = map[string]interface{}{
			"online":   true,
			"lastSeen": time.Now(), // Здесь можно заменить на c.handler.presence.LastSeen(c.userID)
		}
	}
	return status
}
// Auto-generated swagger comments for Broadcast
// @Summary Auto-generated summary for Broadcast
// @Description Auto-generated description for Broadcast — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (h *Hub) Broadcast(ev Event) { h.broadcast <- ev }

// ------------------- Client -------------------

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan Event
	userID  uint
	handler *Handler
}
// Auto-generated swagger comments for addClient
// @Summary Auto-generated summary for addClient
// @Description Auto-generated description for addClient — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (h *Handler) addClient(roomID, userID uint, conn *websocket.Conn) *Client {
	cl := &Client{
		hub:     h.rooms.hub(roomID),
		conn:    conn,
		send:    make(chan Event, 32),
		userID:  userID,
		handler: h,
	}
	cl.hub.register <- cl
	return cl
}
// Auto-generated swagger comments for readPump
// @Summary Auto-generated summary for readPump
// @Description Auto-generated description for readPump — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.handler.presence.Offline(c.userID)
	}()
	for {
		var incoming wsIncoming
		if err := c.conn.ReadJSON(&incoming); err != nil {
			log.Println("read:", err)
			break
		}
		switch incoming.Type {
		case "typing":
			c.hub.Broadcast(Event{Type: "typing", Payload: gin.H{"userId": c.userID}})
		case "message":
			typ, _ := incoming.Payload["type"].(string)
			text, _ := incoming.Payload["text"].(string)
			imageURL, _ := incoming.Payload["imageUrl"].(string)

			msg := models.Message{
				RoomID:   c.hub.roomID,
				UserID:   c.userID,
				Type:     typ,
				Text:     text,
				ImageURL: imageURL,
			}
			if msg.Type == "" {
				msg.Type = "text"
			}
			if err := c.handler.db.Create(&msg).Error; err == nil {
				c.hub.Broadcast(Event{Type: "message", Payload: gin.H{
					"id":        msg.ID,
					"roomId":    msg.RoomID,
					"userId":    msg.UserID,
					"type":      msg.Type,
					"text":      msg.Text,
					"imageUrl":  msg.ImageURL,
					"createdAt": msg.CreatedAt,
				}})
			}
		}
	}
}
// Auto-generated swagger comments for writePump
// @Summary Auto-generated summary for writePump
// @Description Auto-generated description for writePump — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (c *Client) writePump() {
	defer c.conn.Close()
	for ev := range c.send {
		if err := c.conn.WriteJSON(ev); err != nil {
			break
		}
	}
}

// ------------------- Handler WS -------------------
// Auto-generated swagger comments for initWS
// @Summary Auto-generated summary for initWS
// @Description Auto-generated description for initWS — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (h *Handler) initWS(c *gin.Context, roomID uint) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		respondErr(c, 400, "upgrade failed")
		return
	}
	userID := uid(c)
	h.presence.Online(userID)
	cl := h.addClient(roomID, userID, conn)
	go cl.writePump()
	go cl.readPump()
}
// Auto-generated swagger comments for RoomWebSocket
// @Summary Auto-generated summary for RoomWebSocket
// @Description Auto-generated description for RoomWebSocket — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (h *Handler) RoomWebSocket(c *gin.Context) {
	roomIDstr := c.Param("id")
	rid64, _ := strconv.ParseUint(roomIDstr, 10, 64)
	roomID := uint(rid64)

	// Проверка приватной комнаты
	var r models.Room
	if err := h.db.First(&r, roomID).Error; err != nil {
		respondErr(c, 404, "room not found")
		return
	}
	if r.IsPrivate {
		var m models.RoomMember
		if err := h.db.Where("room_id = ? AND user_id = ?", roomID, uid(c)).First(&m).Error; err != nil {
			respondErr(c, 403, "not a member of private room")
			return
		}
	}
	h.initWS(c, roomID)
}
