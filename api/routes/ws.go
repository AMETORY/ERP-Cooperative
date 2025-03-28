package routes

import (
	"ametory-cooperative/config"
	"ametory-cooperative/services"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gopkg.in/olahol/melody.v1"
)

func SetupWSRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	r.GET("/ws/:channelId", func(c *gin.Context) {
		erpContext.AppService.(*services.AppService).Websocket.HandleRequest(c.Writer, c.Request)
	})

	erpContext.AppService.(*services.AppService).Websocket.HandleConnect(func(s *melody.Session) {
		userID, err := parseToken(s.Request.URL.Query().Get("token"))
		if err != nil {
			s.Close()
			return
		}
		msg := gin.H{
			"message":   "Connected",
			"sender_id": *userID,
		}
		b, _ := json.Marshal(msg)
		erpContext.AppService.(*services.AppService).Websocket.BroadcastFilter(b, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
		fmt.Println("Connected", s.Request.URL.Path)
	})
	erpContext.AppService.(*services.AppService).Websocket.HandleDisconnect(func(s *melody.Session) {
		fmt.Println("Disconnected", s.Request.URL.Path)
	})

	erpContext.AppService.(*services.AppService).Websocket.HandleMessage(func(s *melody.Session, msg []byte) {
		erpContext.AppService.(*services.AppService).Websocket.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})
}

func parseToken(authToken string) (*string, error) {

	token, err := jwt.ParseWithClaims(authToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App.Server.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	userID := token.Claims.(*jwt.StandardClaims).Id
	return &userID, nil
}
