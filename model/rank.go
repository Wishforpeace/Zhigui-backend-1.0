package model

func GetRank(offset, limit int) ([]*User, int, error) {
	item := make([]*User, 0)
	d := DB.Table("users").
		Select("users.*").
		Offset(offset).Limit(limit).
		Order("earning").Scan(&item)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	var num int
	users := make([]*User, 0)
	d = DB.Table("users").
		Select("users.*").
		Order("earning").Scan(&users)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	for i, _ := range users {
		num = i + 1
	}
	return item, num, nil
}
