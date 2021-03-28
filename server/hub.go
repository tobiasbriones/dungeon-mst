/*
 * Copyright (c) 2021 Tobias Briones. All rights reserved.
 */

package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"server/ai"
	"server/model"
	"strconv"
	"time"
)

const matchDuration = 5 * time.Second

type Hub struct {
	clients    map[int]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *ResponseData
	quit       chan struct{}
	match      *model.Match
	startTime  time.Time
}

func (h *Hub) Start() {
	var broadcast = func(message *ResponseData) {
		for _, client := range h.clients {
			client.ch <- message
		}
	}

	var register = func(client *Client) {
		remainingTime := matchDuration - time.Since(h.startTime)

		var players []*PlayerJoin

		for _, client := range h.clients {
			players = append(players, &PlayerJoin{
				Id:        client.id,
				Name:      client.name,
				PointJSON: client.PointJSON,
			})
		}

		client.InitGame(h.match, remainingTime, players)

		h.push(client)

		join := &PlayerJoin{
			Id:        client.id,
			Name:      client.name,
			PointJSON: client.PointJSON,
		}
		enc, _ := json.Marshal(join)
		broadcast(&ResponseData{
			Type: DataTypePlayerJoin,
			Body: string(enc),
		})
		go h.listen(client)
	}

	var unregister = func(client *Client) {
		h.delete(client)
		broadcast(&ResponseData{
			Type: DataTypePlayerLeft,
			Body: strconv.Itoa(client.id),
		})
	}

	h.init()

	go func() {
		for {
			time.Sleep(matchDuration)
			h.init()
			matchJSON := model.NewMatchJSON(h.match)
			matchInit := &MatchInit{
				MatchJSON:     matchJSON,
				RemainingTime: matchDuration,
			}
			enc, err := json.Marshal(matchInit)

			if err != nil {
				log.Println("New match error:", err)
				return
			}
			broadcast(&ResponseData{
				Type: DataTypeGameInitialization,
				Body: string(enc),
			})
		}
	}()

	for {
		select {
		case client := <-h.register:
			register(client)
		case client := <-h.unregister:
			unregister(client)
		case message := <-h.broadcast:
			broadcast(message)
		case <-h.quit:
			log.Println("Hub QUIT")
			return
		}
	}
}

func (h *Hub) Register(c *Client) {
	log.Printf("Client %s (%d) connected.\n", c.name, c.id)
	h.register <- c
}

func (h *Hub) Unregister(c *Client) {
	log.Printf("Client %s (%d) disconnected.\n", c.name, c.id)
	h.unregister <- c
}

func (h *Hub) init() {
	h.match = ai.NewRandomMatch()
	h.startTime = time.Now()
}

func (h *Hub) push(client *Client) {
	h.clients[client.id] = client
}

func (h *Hub) delete(client *Client) {
	delete(h.clients, client.id)
}

func (h *Hub) listen(client *Client) {
	conn := client.conn
	id := client.id

	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err) {
				log.Println("Client disconnected", client.id)
			}
			client.Close()
			h.Unregister(client)
			return
		}
		update := &Update{}

		if err := json.Unmarshal(p, update); err != nil {
			log.Println("Parse update error:", err)
			continue
		}
		update.Id = id
		client.PointJSON = update.PointJSON
		enc, err := json.Marshal(update)

		if err != nil {
			log.Println("Encode update error:", err)
			continue
		}

		h.broadcast <- &ResponseData{
			Type: DataTypeUpdate,
			Body: string(enc),
		}
	}
}

func NewHub(ch chan *ResponseData, quit chan struct{}) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  ch,
		quit:       quit,
	}
}
