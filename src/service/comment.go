package service

import (
	"fmt"
	"forum/src/model"
	"forum/src/model/request"
	"time"
)

func CommentTopic(userID, topicID int, req *request.CommentTopic) error {
	_, err := getTopic(topicID)
	if err != nil {
		return err
	}

	now := int(time.Now().Unix())
	tx, _ := model.DB.Beginx()
	_, err = tx.NamedExec("INSERT INTO comment(topic_id, user_id, content, comment_time, status) VALUES(:topic_id, :user_id, :content, :comment_time, :status)",
		model.Comment{
			TopicID:     topicID,
			UserID:      userID,
			Content:     req.Content,
			CommentTime: now,
			Status:      model.CommentStatusNormal,
		},
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("回帖失败")
	}

	_, err = tx.Exec("UPDATE topic SET comment_time = ? WHERE topic_id = ?", now, topicID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("回帖相关操作失败")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("发布回帖失败")
	}

	return nil
}
