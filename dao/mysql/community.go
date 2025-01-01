package mysql

import (
	"Dandelion/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {

	sqlStr := `select community_id,community_name from community`
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据id查询社区详情
func GetCommunityDetailByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select
				community_id,community_name,introduction,create_time 
				from community where community_id = ?
	`
	if err = db.Get(communityDetail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			zap.L().Error("GetCommunityDetailByID/db.Get(community, sqlStr, id) failed: ", zap.Error(err))
		}
	}
	return communityDetail, err
}
