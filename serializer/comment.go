package serializer

import (
	"giligili/model"
)

// Comment 评论序列化器
type CommentResponse struct {
	ID        uint              `json:"id"`
	User      User              `json:"user"`
	VideoId   uint              `json:"video_id"`
	Content   string            `json:"content"`
	ParentID  *uint             `json:"parent_id,omitempty"` // omitempty 表示如果 ParentID 为空则不显示在 JSON 中
	Replies   []CommentResponse `json:"replies,omitempty"`
	CreatedAt int64             `json:"created_at"`
}

// BuildComment 序列化评论
func BuildComment(comment model.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		User:      BuildUser(comment.User),
		VideoId:   comment.VideoId,
		Content:   comment.Content,
		ParentID:  comment.ParentID,
		Replies:   []CommentResponse{}, // 在下面处理
		CreatedAt: comment.CreatedAt.Unix(),
	}
}

// BuildComments 序列化评论切片
func BuildComments(items []model.Comment) []CommentResponse {
	var comments []CommentResponse

	commentMap := make(map[uint]*CommentResponse)
	for _, item := range items {
		// 序列化当前评论
		serializedComment := BuildComment(item)

		if item.ParentID == nil {
			// 如果 ParentID 为空说明是顶级评论，直接添加到结果切片里
			comments = append(comments, serializedComment)

			// 在映射（commentMap）中存储对这个顶级评论序列化对象的引用，以便后续添加回复
			commentMap[item.ID] = &comments[len(comments)-1]
		} else {
			// 如果 ParentID 不为空，则说明是子评论
			// 我们将其添加到其父评论的 Replies 切片中
			if parentComment, ok := commentMap[*item.ParentID]; ok {
				parentComment.Replies = append(parentComment.Replies, serializedComment)
			} else {
				// 如果父评论还没有被添加到映射中（可能是因为父评论被软删除或其他原因）
				// 可以在这里处理这种情况，例如记录日志或者添加一个孤立的评论
				// comments = append(comments, serializedComment)
			}
		}
	}

	return comments
}
