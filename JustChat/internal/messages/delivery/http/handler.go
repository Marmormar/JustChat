package http

import (
	"JustChat/internal/messages/model"
	"JustChat/internal/messages/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	uc usecase.MessageUseCase
}

func NewMessageHandler(r *gin.RouterGroup, uc usecase.MessageUseCase) {
	h := &MessageHandler{uc: uc}
	chats := r.Group("/message")
	{
		chats.GET("/:id", h.GetMessageByID)
		chats.GET("/chat/:chat_id", h.GetMessageByChatID)
		chats.POST("/", h.SaveMessage)
		chats.DELETE("/:id", h.DeleteMessageByID)

	}
}

func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	myuserIDStr := c.GetHeader("X-User-ID")
	myuserID, err := strconv.ParseInt(myuserIDStr, 10, 64)
	if err != nil || myuserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	message, err := h.uc.GetMessageByID(c.Request.Context(), id, myuserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}
func (h *MessageHandler) DeleteMessageByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	myuserIDStr := c.GetHeader("X-User-ID")
	myuserID, err := strconv.ParseInt(myuserIDStr, 10, 64)
	if err != nil || myuserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	err = h.uc.DeleteMessage(c.Request.Context(), id, myuserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": id})
}
func (h *MessageHandler) GetMessageByChatID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	myuserIDStr := c.GetHeader("X-User-ID")
	myuserID, err := strconv.ParseInt(myuserIDStr, 10, 64)
	if err != nil || myuserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	messages, err := h.uc.GetMessagesByChatID(c.Request.Context(), id, myuserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, messages)
}
func (h *MessageHandler) SaveMessage(c *gin.Context) {
	var message model.Message

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	// Получаем X-User-ID из заголовка
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-User-ID header"})
		return
	}
	message.CreatorID = userID
	myuserIDStr := c.GetHeader("X-User-ID")
	myuserID, err := strconv.ParseInt(myuserIDStr, 10, 64)
	if err != nil || myuserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	tosend, err := h.uc.SaveMessage(c.Request.Context(), &message, myuserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tosend)
}
