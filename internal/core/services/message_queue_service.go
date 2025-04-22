package services

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type MessageQueueService struct {
	repo ports.MessageQueueRepository
}

func NewMessageQueueService(repo ports.MessageQueueRepository) ports.MessageQueueService {
	return &MessageQueueService{repo: repo}
}

func (c *MessageQueueService) Create(username string, chatID uuid.UUID) error {
	return c.repo.Create(&domain.MessageQueue{Username: username, ChatID: chatID, Count: 0, Timestamp: time.Now()})
}

func (c *MessageQueueService) Get(username string) ([]domain.MessageQueue, error) {
	return c.repo.GetAll(username)
}

func (c *MessageQueueService) ReceiveMessage(username string, chatID uuid.UUID) error {
	queue, err := c.repo.GetQueue(username, chatID)
	if err != nil {
		return err
	}
	queue.Count += 1
	queue.Timestamp = time.Now()
	err = c.repo.Update(queue)
	if err != nil {
		return err
	}
	return nil
}

func (c *MessageQueueService) ReadMessage(username string, chatID uuid.UUID) error {
	queue, err := c.repo.GetQueue(username, chatID)
	if err != nil {
		return err
	}
	queue.Count = 0
	err = c.repo.Update(queue)
	if err != nil {
		return err
	}
	return nil
}
