package server

import (
	"encoding/json"
	"net/http"
	"psv-generator/internal/generator"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WsHandlerWithChannel(auctionResultChan chan generator.AuctionResult) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

		if err != nil {
			return nil
		}

		defer ws.Close()

		for auctionRes := range auctionResultChan {
			data, _ := json.Marshal(auctionRes)

			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				return err
			}
		}

		return nil
	}
}

func WsHandlerWithHub(hub *BroadcastHub) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

		if err != nil {
			return nil
		}

		done := make(chan struct{})
		clientCh := make(chan generator.AuctionResult, 10)

		hub.Subscribe(clientCh)

		// handle cleanup
		defer func() {
			hub.Unsubscribe(clientCh)
			ws.Close()
		}()

		// read pump
		go func() {
			defer close(done)
			for {
				_, _, err := ws.ReadMessage()
				if err != nil {
					return
				}
			}
		}()

		// write pump
		for {
			select {
			case <-done:
				// handle client disconnect
				return nil
			case auctionRes, ok := <-clientCh:
				if !ok {
					return nil
				}
				data, _ := json.Marshal(auctionRes)
				if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
					return err
				}
			}
		}
	}
}
