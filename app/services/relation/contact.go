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
		if time.Now().Sub(knowMe.CreatedAt) > 1*time.Minute {
			err := global.App.DB.Model(&models.KnowMe{}).
				Where("id = (select temp.id from (select id from know_mes where know_from = ? and know_to = ? "+
					"order by created_at desc limit 1) as temp)", from, to).
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
	} else if knowMe.State == 1 || knowMe.State == 2 || knowMe.State == 4 || knowMe.State == 5 {
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

func (r *relationService) RejectFriendRequest(from, to int) error {
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
		return errors.New("relation state is not 'ready for reject'")
	}

	return updateState(from, to, 1)

}

func (r *relationService) ReleaseFriendRelation(from, to int) error {
	//todo: bi-directional search know_mes table
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
		Select("distinct(wechat_users.id), wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, pets.pet_type,"+
			"wechat_users.gender, wechat_users.age, wechat_users.location,wechat_users.occupation, "+
			"wechat_users.avatar_url, wechat_users.images").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join friends on friends.friend_from = wechat_users.id").
		Where("friends.friend_to = ?", uid).
		Where("friends.deleted_at is null").
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

func (r *relationService) GetRequestedFriendToMe(uid int) (*response.FriendsToMes, error) {
	var friendDto []dto.FriendDto
	var friendsToMes response.FriendsToMes

	d, _ := time.ParseDuration("-24h")
	sevenDays := time.Now().Add(7 * d)

	err := global.App.DB.Model(&models.KnowMe{}).
		Where("know_to", uid).
		Where("created_at < ?", sevenDays).
		Where("state = 0").
		Update("state", 2).Error

	if err != nil {
		zap.L().Error("db know_mes table error", zap.String("update state to 2 failed with error: ", err.Error()))
	}

	err = global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images, "+
			"know_mes.message, know_mes.state").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join know_mes on know_mes.know_from = wechat_users.id").
		Where("know_mes.know_to = ?", uid).
		Where("know_mes.state != 5").
		Where("know_mes.created_at > ?", sevenDays).
		Scan(&friendDto).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			friendsToMes.FriendsInSevenDays = nil
		} else {
			return nil, errors.New(fmt.Sprintf("query requested friend list failed with db error: %s", err.Error()))
		}
	}

	friendsToMes.FriendsInSevenDays = friendDtoToFriendToMeResponse(&friendDto)

	err = global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images, "+
			"know_mes.message, know_mes.state").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join know_mes on know_mes.know_from = wechat_users.id").
		Where("know_mes.know_to = ?", uid).
		Where("know_mes.state != 5").
		Where("know_mes.created_at <= ?", sevenDays).
		Scan(&friendDto).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			friendsToMes.FriendsOutSevenDays = nil
		} else {
			return nil, errors.New(fmt.Sprintf("query requested friend list failed with db error: %s", err.Error()))
		}
	}

	friendsToMes.FriendsOutSevenDays = friendDtoToFriendToMeResponse(&friendDto)

	return &friendsToMes, nil

}

func (r *relationService) GetFriendsInSevenDays(uid int) (*response.FriendsToMes, error) {
	var friendDto []dto.FriendDto
	var friendsToMes response.FriendsToMes

	d, _ := time.ParseDuration("-24h")
	sevenDays := time.Now().Add(7 * d)

	err := global.App.DB.Model(&models.KnowMe{}).
		Where("know_to", uid).
		Where("created_at < ?", sevenDays).
		Where("state = 0").
		Update("state", 2).Error

	if err != nil {
		zap.L().Error("db know_mes table error", zap.String("update state to 2 failed with error: ", err.Error()))
	}

	err = global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images, "+
			"know_mes.message, know_mes.state").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join know_mes on know_mes.know_from = wechat_users.id").
		Where("know_mes.know_to = ?", uid).
		Where("know_mes.state != 5").
		Where("know_mes.created_at > ?", sevenDays).
		Order("know_mes.created_at desc").
		Scan(&friendDto).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			friendsToMes.FriendsInSevenDays = nil
		} else {
			return nil, errors.New(fmt.Sprintf("query requested friend list failed with db error: %s", err.Error()))
		}
	}

	friendsToMes.FriendsInSevenDays = friendDtoToFriendToMeResponse(&friendDto)

	return &friendsToMes, nil
}

func (r *relationService) GetFriendsOutOfSevenDays(uid int) (*response.FriendsToMes, error) {
	var friendDto []dto.FriendDto
	var friendsToMes response.FriendsToMes

	d, _ := time.ParseDuration("-24h")
	sevenDays := time.Now().Add(7 * d)

	err := global.App.DB.Model(&models.KnowMe{}).
		Where("know_to", uid).
		Where("created_at < ?", sevenDays).
		Where("state = 0").
		Update("state", 2).Error

	if err != nil {
		zap.L().Error("db know_mes table error", zap.String("update state to 2 failed with error: ", err.Error()))
	}

	err = global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images, "+
			"know_mes.message, know_mes.state").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join know_mes on know_mes.know_from = wechat_users.id").
		Where("know_mes.know_to = ?", uid).
		Where("know_mes.state != 5").
		Where("know_mes.created_at <= ?", sevenDays).
		Order("know_mes.created_at desc").
		Scan(&friendDto).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			friendsToMes.FriendsOutSevenDays = nil
		} else {
			return nil, errors.New(fmt.Sprintf("query requested friend list failed with db error: %s", err.Error()))
		}
	}

	friendsToMes.FriendsOutSevenDays = friendDtoToFriendToMeResponse(&friendDto)

	return &friendsToMes, nil
}

func (r *relationService) GetMyFriendRequest(uid int) (*response.MesToFriends, error) {
	var friendDto []dto.FriendDto
	var friends response.MesToFriends

	d, _ := time.ParseDuration("-24h")
	sevenDays := time.Now().Add(7 * d)

	err := global.App.DB.Model(&models.KnowMe{}).
		Where("know_from", uid).
		Where("created_at < ?", sevenDays).
		Where("state = 0").
		Update("state", 2).Error

	if err != nil {
		zap.L().Error("db know_mes table error", zap.String("update state to 2 failed with error: ", err.Error()))
		return nil, nil
	}

	err = global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images, "+
			"know_mes.message, know_mes.state, know_mes.created_at").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join know_mes on know_mes.know_to = wechat_users.id").
		Where("know_mes.know_from = ?", uid).
		Where("know_mes.state != 5").
		Order("know_mes.created_at desc").
		Scan(&friendDto).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, errors.New(fmt.Sprintf("query requested friend list failed with db error: %s", err.Error()))
		}
	}

	friends.MyFriendRequests = friendDtoToFriendToMyFriendRequest(&friendDto)

	return &friends, nil
}

//GetFriendStatus
//friendStatus -1-想认识ta 0-已申请 1-发消息 2-去同意
func (r *relationService) GetFriendStatus(from, to int) int {
	var state int
	res := global.App.DB.Model(&models.KnowMe{}).
		Select("state").
		Where("know_from = ? and know_to = ? ", from, to).
		Order("created_at desc").
		Limit(1).
		First(&state)
	if res.Error == gorm.ErrRecordNotFound {
		resRev := global.App.DB.Model(&models.KnowMe{}).
			Select("state").
			Where("know_from = ? and know_to = ? ", to, from).
			Order("created_at desc").
			Limit(1).
			First(&state)
		if resRev.Error == gorm.ErrRecordNotFound {
			return -1
		} else if res.Error != nil {
			return -1
		}

		if state == 0 {
			return 2
		} else if state == 2 {
			return 1
		} else {
			return -1
		}
	} else if res.Error != nil {
		return -1
	}

	if state == 0 {
		return 0
	} else if state == 3 {
		return 1
	} else {
		return -1
	}
}

func (r *relationService) UpdateAllNewFriendRequestStatus(uid int) error {
	res := global.App.DB.Model(&models.KnowMe{}).
		Where("know_to = ? and status = 0", uid).
		Update("status", 1)
	if res.RowsAffected == 0 {
		zap.L().Info("no new friend request to update", zap.String("db info", "no new friend request record found to update"))
		return nil
	}

	if res.Error != nil {
		return errors.New(fmt.Sprintf("update new friend request status failed with error, %s", res.Error.Error()))
	}

	return nil
}

//general method to update user's contact state
func updateState(from, to, state int) error {
	err := global.App.DB.Model(&models.KnowMe{}).
		Where("id = (select temp.id from (select id from know_mes where know_from = ? and know_to = ? "+
			"order by created_at desc limit 1) as temp)", from, to).
		Update("state", state).Error
	if err != nil {
		return err
	}
	return nil
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
			"order by created_at desc limit 1) as temp)", from, to).
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
			"order by created_at desc limit 1) as temp)", from, to).
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
			PetType:    friendDto.PetType,
			Gender:     friendDto.Gender,
			Age:        friendDto.Age,
			Location:   friendDto.Location,
			Occupation: friendDto.Occupation,
			Images:     image,
		}

		friends = append(friends, friend)

	}

	return friends
}

func friendDtoToFriendToMeResponse(friendsDtos *[]dto.FriendDto) []response.FriendToMeResponse {
	var friends []response.FriendToMeResponse
	var ids []int

	for _, friendDto := range *friendsDtos {
		var username string
		var image string

		if utils.Contains(ids, friendDto.Id) {
			continue
		}

		ids = append(ids, friendDto.Id)

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

		friend := response.FriendToMeResponse{
			Id:       friendDto.Id,
			UserName: username,
			PetName:  friendDto.PetName,
			Images:   image,
			Message:  friendDto.Message,
			State:    friendDto.State,
		}

		friends = append(friends, friend)

	}

	return friends
}

func friendDtoToFriendToMyFriendRequest(friendsDtos *[]dto.FriendDto) []response.MyFriendRequest {
	var friends []response.MyFriendRequest
	var ids []int

	for _, friendDto := range *friendsDtos {
		var username string
		var image string

		if utils.Contains(ids, friendDto.Id) {
			continue
		}

		ids = append(ids, friendDto.Id)

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

		friend := response.MyFriendRequest{
			Id:        friendDto.Id,
			UserName:  username,
			PetName:   friendDto.PetName,
			Images:    image,
			State:     friendDto.State,
			CreatedAt: friendDto.CreatedAt.String()[0:10],
		}

		friends = append(friends, friend)

	}

	return friends
}
