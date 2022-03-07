package model

type Read5WhiteListModel struct {
	Id         int    `json:"id" primaryKey:"true"`
	ActivityId int    `json:"activity_id"`
	Userid     int    `json:"userid"`
	Createtime int    `json:"createtime"`
	UserName   string `json:"userName" gorm:"column:userName"`
}

func (Read5WhiteListModel) TableName() string {
	return "read5_white_list"
}

func (Read5WhiteListModel) GetId(activityId int, userid int) (id int) {
	var info Read5WhiteListModel
	Db.Where("activity_id = ? and userid = ?", activityId, userid).Limit(1).Find(&info)
	if info.Id == 0 {
		return 0
	} else {
		return info.Id
	}
}
