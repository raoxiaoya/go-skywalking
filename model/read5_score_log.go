package model

type Read5ScoreLogModel struct {
	Id         int    `json:"id" primaryKey:"true"`
	ActivityId int    `json:"activity_id"`
	Userid     int    `json:"userid"`
	FromType   int    `json:"from_type"`
	Type       int    `json:"type"`
	AwardType  int    `json:"award_type"`
	AwardValue int    `json:"award_value"`
	Createat   int64  `json:"createat"`
	Createtime string `json:"createtime"`
}

func (Read5ScoreLogModel) TableName() string {
	return "read5_score_log"
}

func (Read5ScoreLogModel) GetId(activityId int, userid int) (id int) {
	var info Read5ScoreLogModel
	Db.Where("activity_id = ? and userid = ? and from_type = 1", activityId, userid).Limit(1).Find(&info)
	if info.Id == 0 {
		return 0
	} else {
		return info.Id
	}
}