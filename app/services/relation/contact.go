package relation

import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/app/services/dto"
	"webapp_gin/global"
	"webapp_gin/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

//AddNewContact
//state: 0 - 待处理 1 - 已婉拒 2 - 过期自动拒绝 3 - 已同意 4 - 已解除
//status: 0 - 新认识 1 - 已查看
func (r *relationService) AddNewContact(from, to int, message string) error {
	var knowMe models.KnowMe
	var wechatUser models.WechatUser
	err := global.App.DB.Where("id = ?", from).First(&wechatUser).Error
	if err != nil {
		return errors.New("no user found")
	}

	res := global.App.DB.Model(&models.KnowMe{}).Where("know_from = ? and know_to = ?", from, to).Order("created_at desc").First(&knowMe)

	if res.Error == gorm.ErrRecordNotFound {
		knowMe.KnowFrom = strconv.Itoa(from)
		knowMe.KnowTo = strconv.Itoa(to)
		knowMe.Status = 0
		knowMe.State = 0
		knowMe.Message = message
		err = global.App.DB.Create(&knowMe).Error
		if err != nil {
			return errors.New(fmt.Sprintf("create contact record db error: %s", err.Error()))
		} else {
			return nil
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("query contact record db error: %s", err.Error()))
	}

	if knowMe.State == 0 {
		if time.Now().Sub(knowMe.CreatedAt) > 7*24*60*time.Minute {
			err := global.App.DB.Model(&models.KnowMe{}).
				Where("id = (select temp.id from (select id from know_mes where know_from = ? and know_to = ? "+
					"order by created_at desc limit 1) as temp", from, to).
				Update("state", 2).Error
			if err != nil {
				zap.L().Warn("knowMe db error", zap.String("update state error", fmt.Sprintf("from: %v to: %v", from, to)))
			}
			knowMe := models.KnowMe{}
			knowMe.KnowFrom = strconv.Itoa(from)
			knowMe.KnowTo = strconv.Itoa(to)
			knowMe.Status = 0
			knowMe.State = 0
			knowMe.Message = message
			err = global.App.DB.Create(&knowMe).Error
			if err != nil {
				return errors.New(fmt.Sprintf("create contact record db error: %s", err.Error()))
			} else {
				return nil
			}
		} else {
			return errors.New("previous request still valid")
		}
	} else if knowMe.State == 2 || knowMe.State == 3 {
		knowMe := models.KnowMe{}
		knowMe.KnowFrom = strconv.Itoa(from)
		knowMe.KnowTo = strconv.Itoa(to)
		knowMe.Status = 0
		knowMe.State = 0
		knowMe.Message = message
		err = global.App.DB.Create(&knowMe).Error
		if err != nil {
			return errors.New(fmt.Sprintf("create contact record db error: %s", err.Error()))
		} else {
			return nil
		}
	} else {
		return errors.New("previous request still valid")
	}
}

func (r *relationService) ApproveFriendRequest(from, to int) error {
	var knowMe models.KnowMe
	res := global.App.DB.Model(&models.KnowMe{}).
		Where("know_from = ? and know_to = ?", from, to).
		Order("created_at desc").
		First(&knowMe)
	if res.Error == gorm.ErrRecordNotFound {
		return errors.New("no relation record found in db")
	} else if res.Error != nil {
		return errors.New(fmt.Sprintf("query relation db failed woth error: %s", res.Error.Error()))
	}

	if knowMe.State != 0 {
		return errors.New("relation state is not 'ready for approval'")
	}

	return updateStateAndCreateFriend(global.App.DB, from, to)

}

func (r *relationService) ReleaseFriendRelation(from, to int) error {
	var knowMe models.KnowMe
	res := global.App.DB.Model(&models.KnowMe{}).
		Where("know_from = ? and know_to = ?", from, to).
		Order("created_at desc").
		First(&knowMe)
	if res.Error == gorm.ErrRecordNotFound {
		return errors.New("no relation record found in db")
	} else if res.Error != nil {
		return errors.New(fmt.Sprintf("query relation db failed woth error: %s", res.Error.Error()))
	}

	if knowMe.State != 3 {
		return errors.New("relation state is not 'ready for releasing'")
	}

	return updateStateAndDeleteFriend(global.App.DB, from, to)
}

func (r *relationService) GetFriendList(uid int) (*response.Friends, error) {
	var friends []models.Friend
	var friendDto []dto.FriendDto
	res := global.App.DB.Where("friend_to = ?", uid).Find(&friends)

	if res.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if res.Error != nil {
		zap.L().Error("database error", zap.String("looking for user friend error", res.Error.Error()))
		return nil, errors.New(fmt.Sprintf("looking for user's friend failed with db error: %s", res.Error.Error()))
	}

	err := global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join friends on friends.friend_from = wechat_users.id").
		Where("friends.friend_to = ?", uid).
		Where("friends.deleted_at = null").
		Scan(&friendDto).Error

	if err != nil {
		zap.L().Error("database error", zap.String("looking for user friend error", res.Error.Error()))
		return nil, errors.New(fmt.Sprintf("looking for user's friend failed when fetch info from "+
			"other tabls with db error: %s", res.Error.Error()))
	}

	friendsRes := friendDtoToFriendResponse(&friendDto)

	myFriends := response.Friends{
		MyFriends: friendsRes,
	}

	return &myFriends, nil

}

func updateStateAndCreateFriend(db *gorm.DB, from, to int) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(&models.KnowMe{}).
		Where("id = (select temp.id from (select id from know_mes where know_from = ? and know_to = ? "+
			"order by created_at desc limit 1) as temp", from, to).
		Update("state", 3).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	var friends = []models.Friend{
		{FriendFrom: strconv.Itoa(from), FriendTo: strconv.Itoa(to)},
		{FriendFrom: strconv.Itoa(to), FriendTo: strconv.Itoa(from)},
	}

	err = tx.Model(&models.Friend{}).Create(&friends).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func updateStateAndDeleteFriend(db *gorm.DB, from, to int) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(&models.KnowMe{}).
		Where("id = (select temp.id from (select id from know_mes where know_from = ? and know_to = ? "+
			"order by created_at desc limit 1) as temp", from, to).
		Update("state", 5).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	deletedAt := time.Now()

	err = tx.Model(&models.Friend{}).
		Where("friend_from = ? and friend_to = ?", from, to).
		Update("deleted_at", deletedAt).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&models.Friend{}).
		Where("friend_from = ? and friend_to = ?", to, from).
		Update("deleted_at", deletedAt).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func friendDtoToFriendResponse(friendsDtos *[]dto.FriendDto) []response.Friend {
	var friends []response.Friend
	for _, friendDto := range *friendsDtos {
		var username string
		var image string

		if friendDto.UserName == "" {
			username = friendDto.WechatName
		} else {
			username = friendDto.UserName
		}

		if friendDto.Images == "" {
			image = friendDto.AvatarUrl
		} else {
			image = utils.ParseToArray(&friendDto.Images, " ")[0]
		}

		friend := response.Friend{
			Id:         friendDto.Id,
			UserName:   username,
			PetName:    friendDto.PetName,
			Age:        friendDto.Age,
			Location:   friendDto.Location,
			Occupation: friendDto.Occupation,
			Images:     image,
		}

		friends = append(friends, friend)

	}

	return friends
}
