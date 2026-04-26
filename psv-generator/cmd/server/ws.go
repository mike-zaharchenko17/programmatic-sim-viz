package server

import (
	"encoding/json"
	"psv-generator/internal/generator"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var (
	upgrader = websocket.Upgrader{}
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

		clientCh := make(chan generator.AuctionResult, 10)

		defer ws.Close()

		hub.Subscribe(clientCh)
		defer hub.Unsubscribe(clientCh)

		for auctionRes := range clientCh {
			data, _ := json.Marshal(auctionRes)

			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				return err
			}
		}

		return nil
	}
}
