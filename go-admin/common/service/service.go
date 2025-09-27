package service

import (
	"fmt"

	"github.com/go-admin-team/go-admin-core/logger"
	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Msg   string
	MsgID string
	Log   *logger.Helper
	Error error
}

// AddError 方法用於將錯誤添加到 Service 結構中。
// 如果當前的 Error 為空，則將傳入的錯誤設置為當前錯誤。
func (db *Service) AddError(err error) error {
	if db.Error == nil {
		db.Error = err
	} else if err != nil {
		db.Error = fmt.Errorf("%v; %w", db.Error, err)
	}
	return db.Error
}
