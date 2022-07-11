package forum

import (
	"zhigui/model"
)

type CommentResponse struct {
	Publisher string
	Content   string
}

func PostMessage(email, content, title string) (*model.POST, error) {

	post, err := model.PostMessage(email, content, title)
	if err != nil {
		return nil, err
	}
	return post, nil

}

func GetMyPosts(email string, offset int, limit int) ([]*model.POST, int, error) {
	msg, num, err := model.GetMyPosts(email, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return msg, num, nil
}

func PostComment(ID uint, content string, email string) (CommentResponse, error) {
	comment, err := model.CreateComment(ID, content, email)
	var resp CommentResponse
	if err != nil {
		return resp, err
	} else {
		resp = CommentResponse{
			Publisher: comment.Publisher,
			Content:   comment.Content,
		}
		return resp, nil
	}
}

func GetComments(id string, offset int, limit int) ([]model.Comment, int, error) {
	var comments []model.Comment
	comments, num, err := model.GetComments(id, offset, limit)
	return comments, num, err
}

// 获取我的评论
func GetMyComments(email string, offset int, limit int) ([]*model.Comment, int, error) {
	msg, num, err := model.GetMyComments(email, offset, limit)
	if err != nil {
		return nil, num, err
	}
	return msg, num, nil
}

// 删除我的帖子
func DeletePost(email string, id string) error {
	err := model.DeletePost(email, id)
	return err
}

// 删除我的评论
func DeleteComment(email string, id string) error {
	err := model.DeleteComment(email, id)
	return err
}
