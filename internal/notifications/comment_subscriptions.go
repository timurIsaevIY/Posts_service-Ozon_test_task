package notifications

import (
	"context"
	"sync"

	"Ozon_Post_comment_system/internal/models"
)

type Observers interface {
	Subscribe(ctx context.Context, postID int) (<-chan *models.Comment, error)
	Unsubscribe(postID int, targetChan chan *models.Comment)
	Notify(postID int, comment *models.Comment)
}

type Observer struct {
	mu          sync.RWMutex
	subscribers map[int][]chan *models.Comment // postID -> список подписчиков
}

func NewObserver() *Observer {
	return &Observer{
		subscribers: make(map[int][]chan *models.Comment),
	}
}

func (o *Observer) Subscribe(ctx context.Context, postID int) (<-chan *models.Comment, error) {
	commentChan := make(chan *models.Comment, 1) // Буферизованный канал
	o.mu.Lock()
	o.subscribers[postID] = append(o.subscribers[postID], commentChan)
	o.mu.Unlock()

	// 🔥 Авто-отписка при разрыве соединения
	go func() {
		<-ctx.Done()
		o.Unsubscribe(postID, commentChan)
	}()

	return commentChan, nil
}

func (o *Observer) Unsubscribe(postID int, targetChan chan *models.Comment) {
	o.mu.Lock()
	defer o.mu.Unlock()

	subs, exists := o.subscribers[postID]
	if !exists {
		return
	}

	// 🔥 Удаляем `targetChan` из списка подписчиков
	newSubs := subs[:0] // Создаём новый список подписчиков
	for _, ch := range subs {
		if ch != targetChan {
			newSubs = append(newSubs, ch) // Оставляем только активные подписчики
		} else {
			close(ch) // ✅ Закрываем канал, чтобы избежать утечек памяти
		}
	}

	// Если подписчиков больше нет, удаляем `postID`
	if len(newSubs) == 0 {
		delete(o.subscribers, postID)
	} else {
		o.subscribers[postID] = newSubs
	}
}

func (o *Observer) Notify(postID int, comment *models.Comment) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	subs, exists := o.subscribers[postID]
	if !exists {
		return
	}

	for _, ch := range subs {
		select {
		case ch <- comment:
		default:
		}
	}
}
