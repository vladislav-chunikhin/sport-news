package transformer

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/vladislav-chunikhin/feed-transformer/internal/mapper"
	feedRepoPkg "github.com/vladislav-chunikhin/feed-transformer/internal/repository/feed"
	"github.com/vladislav-chunikhin/lib-go/pkg/redis"
)

const (
	typeKey       = "type"
	cacheMsgIDTTL = 5 * time.Minute
)

func (s *Service) handle(msg amqp.Delivery) error {
	body := msg.Body
	if body == nil || len(body) == 0 {
		return nil
	}

	msgType := ""
	if msgTypeValue, ok := msg.Headers[typeKey]; ok {
		if msgTypeValueStr, okType := msgTypeValue.(string); okType {
			msgType = msgTypeValueStr
		}
	}

	var isProcessed bool
	if err := s.cache.Get(s.ctx, msg.MessageId, &isProcessed); err != nil {
		if !errors.Is(err, redis.ErrCacheMiss) {
			return fmt.Errorf("failed to get msg id from cache: %v", err)
		}
	}

	if isProcessed {
		s.logger.Debugf("there was an attempt to process the message twice, id = %s", msg.MessageId)
		return nil
	}

	var articleItems ArticleItems
	if err := json.Unmarshal(body, &articleItems); err != nil {
		return err
	}

	if len(articleItems.Articles) == 0 {
		return nil
	}

	externalIDs := make([]int, 0, len(articleItems.Articles))
	articleByExtIDMap := make(map[int]*ArticleItem, len(articleItems.Articles))
	for _, articleItem := range articleItems.Articles {
		if articleItem.ID != 0 {
			externalIDs = append(externalIDs, articleItem.ID)
			articleByExtIDMap[articleItem.ID] = articleItem
		}
	}

	if len(externalIDs) != 0 {
		articles, err := s.feedRepository.GetByExternalIDs(s.ctx, externalIDs, "_id", "externalId", "updated")
		if err != nil {
			return err
		}

		if len(articles) != 0 {
			for _, article := range articles {
				if article.Updated == nil {
					continue
				}

				if articleItem, ok := articleByExtIDMap[article.ExternalID]; ok {
					if articleItem.Updated == *article.Updated {
						delete(articleByExtIDMap, article.ExternalID)
					}
				}
			}
		}
	}

	if len(articleByExtIDMap) == 0 {
		return nil
	}

	updatedArticles := make([]*feedRepoPkg.Article, 0, len(articleByExtIDMap))
	for extID := range articleByExtIDMap {
		updatedArticleItem, err := s.feedClient.GetNewsContentByID(s.ctx, extID)
		if err != nil {
			return err
		}

		var article *feedRepoPkg.Article
		article, err = mapper.ToArticle(updatedArticleItem, msgType)
		if err != nil {
			s.logger.Errorf("failed to map an article item to entity")
			continue
		}
		if article != nil {
			updatedArticles = append(updatedArticles, article)
		}
	}

	if len(updatedArticles) == 0 {
		return nil
	}

	if err := s.feedRepository.UpdateOrCreateArticles(s.ctx, updatedArticles); err != nil {
		return err
	}

	isProcessed = true
	if err := s.cache.Set(s.ctx, msg.MessageId, &isProcessed, cacheMsgIDTTL); err != nil {
		s.logger.Errorf("failed to set message ID after successful processing: %v", err)
	}

	return nil
}
