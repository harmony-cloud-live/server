package server

import (
	"context"
	"net/http"
	"net/url"

	"nhooyr.io/websocket"
)

type eventHandler func(ctx context.Context, c *websocket.Conn, userId string) error
type closeFunc func(ctx context.Context, userId string) error

func (h *HarmonyCloudServer) upgradeMidiSocket(w http.ResponseWriter, r *http.Request) {
	h.upgradeSocket(w, r, h.handleMidiEvent, nil)
}

func (h *HarmonyCloudServer) upgradeControlSocket(w http.ResponseWriter, r *http.Request) {
	h.upgradeSocket(w, r, h.handleControlEvent, h.closeControlSocket)
}

func (h *HarmonyCloudServer) closeControlSocket(ctx context.Context, userId string) error {
	h.logf("closeControlSocket: %v", h.clients[userId].Username)
	delete(h.clients, userId)
	if h.leader != nil && h.leader.UserId == userId {
		for _, c := range h.clients {
			if c != nil {
				err := h.setLeader(ctx, c)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return h.marshalAndBroadcast(ctx, "", GetClients, h.getClients())
}

func (h *HarmonyCloudServer) upgradeSocket(w http.ResponseWriter, r *http.Request, handler eventHandler, close closeFunc) {
	userId, err := getUserId(r)
	if err != nil {
		h.logf("error getting userId %v", err)
		return	

	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		h.logf("upgradeSocket error: %v", err)
		return	
	}
	defer c.CloseNow()

	for {
		err = handler(r.Context(), c, userId)
		if err != nil {
			if close != nil {
				close(r.Context(), userId)
			}
			return
		}
	}
}

func getUserId(req *http.Request) (string, error) {
    parsedUrl, err := url.Parse(req.URL.String())
    if err != nil {
        return "", err
    }
    userId := parsedUrl.Query().Get("userId")
    return userId, nil
}

func (h *HarmonyCloudServer) broadcast(ctx context.Context, senderId string, msg []byte) error {
	for _, c := range h.clients {
		if c.UserId != senderId {
			err := c.Conn.Write(ctx, websocket.MessageBinary, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
