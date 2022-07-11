package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type POST struct {
	gorm.Model
	Publisher  string
	Title      string `gorm:"size:100"`
	Content    string `gorm:"size:255"`
	LikeNum    int
	CommentNum int
	Pic        []PostImage
	Remarks    []Comment
}

type PostsResponse struct {
	User  Info
	Posts POST
}
type Info struct {
	Name  string
	Email string
	IMG   string
}
type Comment struct {
	gorm.Model
	PostID    uint `gorm:"column:post_id"`
	Publisher string
	Content   string `gorm:"size:255"`
}

type LikeList struct {
	gorm.Model
	PostID uint
	Email  string
}

func GetPosts(offset int, limit int) ([]*PostsResponse, int, error) {
	var item []POST
	d := DB.Table("posts").
		Offset(offset).Limit(limit).
		Order("created_at desc").Find(&item)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	num := 0
	var posts []POST
	d = DB.Table("posts").
		Select("posts.*").
		Order("created_at desc").Find(&posts)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	for i, _ := range posts {
		num = i + 1
	}
	resp := make([]*PostsResponse, len(item))
	for i, m := range item {
		u, _ := GetInfo(m.Publisher)
		resp[i] = &PostsResponse{
			User: Info{
				Name:  u.NickName,
				Email: u.Email,
				IMG:   u.Image.Avatar,
			},
			Posts: item[i],
		}
	}
	return resp, num, d.Error
}

func GetPostDetails(offset int, limit int, id int) (POST, int, error) {
	var item POST
	var comment []Comment
	Num := 0
	DB.Where("id = ?", id).First(&item)
	d := DB.Table("comments").
		Where(" post_id = ?", id).
		Offset(offset).Limit(limit).
		Order("created_at desc").Scan(&comment)
	if d.Error != nil {
		return item, 0, d.Error
	}
	item.Remarks = comment
	for i, _ := range comment {
		Num = i + 1
	}
	return item, Num, d.Error
}

func PostMessage(email, content, title string) (*POST, error) {
	var post POST = POST{
		Publisher:  email,
		Title:      title,
		Content:    content,
		LikeNum:    0,
		CommentNum: 0,
	}
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	err := tx.Commit().Error
	return &post, err
}

func GetMyPosts(email string, offset int, limit int) ([]*POST, int, error) {
	item := make([]*POST, 0)

	d := DB.Table("posts").
		Where(" publisher = ?", email).
		Offset(offset).Limit(limit).
		Order("created_at desc").Find(&item)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	var num int
	for i, m := range item {
		var pic []PostImage
		DB.Where("post_id = ?", m.ID).Find(&pic)
		item[i].Pic = pic
		var remark []Comment
		DB.Where("post_id = ?", m.ID).Find(&remark)
		item[i].Remarks = remark
	}
	posts := make([]*POST, 0)
	d = DB.Table("posts").
		Where(" publisher = ?", email).
		Order("created_at desc").Find(&posts)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	for i, _ := range posts {
		num = i + 1
	}
	return item, num, d.Error

}

func GetMyComments(email string, offset int, limit int) ([]*Comment, int, error) {
	item := make([]*Comment, 0)

	d := DB.Table("comments").
		Where(" publisher = ?", email).
		Offset(offset).Limit(limit).
		Order("created_at desc").Find(&item)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	var num int
	for i, _ := range item {
		num = i + 1
	}
	return item, num, d.Error
}
func (post *POST) UploadPostPictures(ID, picUrl, picSha, picPath string) error {
	id, _ := strconv.Atoi(ID)
	Id := uint(id)
	var image = PostImage{
		PostID: Id,
		Url:    picUrl,
		Sha:    picSha,
		Path:   picPath,
	}
	if err := image.UploadPostPicture(); err != nil {
		return err
	}
	if err := DB.Where("id = ?", id).First(&post).Error; err != nil {
		return err
	}
	var IMG []PostImage
	if err := DB.Where("post_id=?", id).Find(&IMG).Error; err != nil {
		return err
	}
	post.Pic = IMG
	return nil
}

func CreateComment(ID uint, content string, email string) (Comment, error) {
	comment := Comment{
		PostID:    ID,
		Publisher: email,
		Content:   content,
	}
	var post POST
	DB.Where("id = ?", ID).First(&post)
	post.CommentNum += 1
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return Comment{}, err
	}
	if err := tx.Model(&POST{}).Where("id = ?", ID).Updates(POST{CommentNum: post.CommentNum}).Error; err != nil {
		tx.Rollback()
		return Comment{}, err
	}
	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		return Comment{}, err
	}
	return comment, tx.Commit().Error
}

// 获取点评
func GetComments(id string, offset int, limit int) ([]Comment, int, error) {
	var item []Comment
	d := DB.Table("comments").
		Where("post_id = ?", id).
		Offset(offset).Limit(limit).
		Order("created_at desc").Find(&item)
	if err := d.Error; err != nil {
		return item, 0, err
	}
	var num int
	var comments []Comment
	d = DB.Table("comments").
		Where("post_id = ?", id).
		Offset(offset).Limit(limit).
		Order("created_at desc").Find(&comments)
	if err := d.Error; err != nil {
		return item, 0, err
	}
	for i, _ := range comments {
		num = i + 1
	}
	return item, num, nil
}

// 判断是否点赞
func CheckLikeList(id string, email string) error {
	var list LikeList
	if err := DB.Where("post_id =? AND email = ?", id, email).First(&list).Error; err == nil {
		var post POST
		DB.Where("id = ?", id).First(&post)
		post.LikeNum -= 1
		DB.Model(&POST{}).Where("id = ?", id).Updates(POST{LikeNum: post.LikeNum})
		return DB.Where("post_id =? AND email = ?", id, email).Delete(&list).Error
	} else if err != nil {
		return err
	}
	return nil
}

// 点赞
func CreateLike(id string, email string) error {
	var post POST
	var like LikeList
	ID, _ := strconv.Atoi(id)
	Id := uint(ID)
	DB.Where("id = ?", id).First(&post)
	post.LikeNum += 1
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Model(&POST{}).Where("id = ?", id).Updates(POST{LikeNum: post.LikeNum}).Error; err != nil {
		tx.Rollback()
		return err
	}
	like = LikeList{
		PostID: Id,
		Email:  email,
	}
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

}

// 删除帖子
func DeletePost(email string, id string) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Where("id = ? AND publisher = ?", id, email).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id = ?", id).Delete(&POST{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
	return nil
}

//删除评论
func DeleteComment(email string, id string) error {
	var comment Comment
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id = ? and publisher = ?", id, email).First(&comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	var post POST
	DB.Where("id = ?", comment.PostID).First(&post)
	post.CommentNum -= 1
	if err := tx.Model(&POST{}).Where("id = ?", comment.PostID).Updates(POST{CommentNum: post.CommentNum}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
