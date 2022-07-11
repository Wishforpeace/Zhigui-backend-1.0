package rank

import (
	"errors"
	"zhigui/model"
)

type RankResponse struct {
	ID       uint   `json:"ID"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
	Gender   string `json:"gender"`
	Earning  string `json:"earning"`
}

func GetRank(offset int, limit int) ([]*RankResponse, int, error) {

	item, num, err := model.GetRank(offset, limit)
	if err != nil {
		return nil, 0, errors.New("获取失败")
	}

	resp := make([]*RankResponse, len(item))
	for i, m := range item {
		resp[i] = &RankResponse{
			ID:       m.ID,
			Email:    m.Email,
			NickName: m.NickName,
			Gender:   m.Gender,
			Earning:  m.Earning,
		}
	}
	return resp, num, nil
}
