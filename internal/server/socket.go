package server

import (
	"context"
	"net/http"
	"net/url"

	"nhooyr.io/websocket"
)


type eventHandler func(ctx context.Context, c *websocket.Conn, userId string) error

func (h *HarmonyCloudServer) upgradeMidiSocket(w http.ResponseWriter, r *http.Request) {
	h.upgradeSocket(w, r, h.handleMidiEvent)
}

func (h *HarmonyCloudServer) upgradeControlSocket(w http.ResponseWriter, r *http.Request) {
	h.upgradeSocket(w, r, h.handleControlEvent)
}

func (h *HarmonyCloudServer) upgradeSocket(w http.ResponseWriter, r *http.Request, handler eventHandler) {
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
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure || 
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			return
		}
		if err != nil {
			h.logf("error handling event: %v, userId: %v", err, userId)
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