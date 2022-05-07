package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"webapp_gin/utils/date"
)

func MakeSquareQueryIn3Day(uid int, limit int, sq SquareSetting) (string, error, int) {
	// 拼接sql查询
	var resQuery string
	rule, err, errorCode := RuleToQuery(&sq)
	if err != nil {
		return "", errors.New("广场条件设置错误"), errorCode
	}
	startTime := date.GetDateByOffsetDay(-3, time.Now())
	endTime := date.GetDateByOffsetDay(0, time.Now())
	resQuery = fmt.Sprintf("id != %s and created_at > %s and created_at < %s and %s order by rand() limit %s", strconv.Itoa(uid), startTime, endTime, rule, strconv.Itoa(limit))

	return resQuery, nil, 0

}

func MakeSquareQueryBefore3Day(uid int, limit int, sq SquareSetting) (string, error, int) {
	var resQuery string
	rule, err, errorCode := RuleToQuery(&sq)
	if err != nil {
		return "", errors.New("广场条件设置错误"), errorCode
	}
	endTime := date.GetDateByOffsetDay(-3, time.Now())
	resQuery = fmt.Sprintf("id != %s and created_at  < %s and %s order by rand() limit %s", strconv.Itoa(uid), endTime, rule, strconv.Itoa(limit))
	return resQuery, nil, 0
}
