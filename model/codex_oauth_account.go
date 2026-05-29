package model

import (
	"github.com/songquanpeng/one-api/common/helper"
	"gorm.io/gorm"
)

const CodexOAuthProvider = "codex_oauth"

type CodexOAuthAccount struct {
	Id              int    `json:"id"`
	AccountID       string `json:"account_id" gorm:"uniqueIndex;size:191"`
	Email           string `json:"email" gorm:"size:191"`
	RefreshToken    string `json:"-" gorm:"type:text"`
	AuthenticatedAt int64  `json:"authenticated_at" gorm:"bigint"`
	IsDefault       bool   `json:"is_default" gorm:"default:false"`
	CreatedTime     int64  `json:"created_time" gorm:"bigint"`
	UpdatedTime     int64  `json:"updated_time" gorm:"bigint"`
}

func GetCodexOAuthAccounts() ([]*CodexOAuthAccount, error) {
	var accounts []*CodexOAuthAccount
	err := DB.Order("is_default desc, id asc").Find(&accounts).Error
	return accounts, err
}

func GetDefaultCodexOAuthAccount() (*CodexOAuthAccount, error) {
	account := CodexOAuthAccount{}
	err := DB.Where("is_default = ?", true).First(&account).Error
	if err == nil {
		return &account, nil
	}
	err = DB.Order("id asc").First(&account).Error
	return &account, err
}

func GetCodexOAuthAccountByAccountID(accountID string) (*CodexOAuthAccount, error) {
	account := CodexOAuthAccount{}
	err := DB.Where("account_id = ?", accountID).First(&account).Error
	return &account, err
}

func UpsertCodexOAuthAccount(account *CodexOAuthAccount) error {
	if account.RefreshToken != "" {
		account.RefreshToken = codexOAuthProtectRefreshToken(account.RefreshToken)
	}
	now := helper.GetTimestamp()
	if account.CreatedTime == 0 {
		account.CreatedTime = now
	}
	account.UpdatedTime = now
	if account.AuthenticatedAt == 0 {
		account.AuthenticatedAt = now
	}

	var existing CodexOAuthAccount
	err := DB.Where("account_id = ?", account.AccountID).First(&existing).Error
	if err == nil {
		account.Id = existing.Id
		account.CreatedTime = existing.CreatedTime
		return DB.Model(&existing).Updates(account).Error
	}
	return DB.Create(account).Error
}

func SetDefaultCodexOAuthAccount(accountID string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&CodexOAuthAccount{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&CodexOAuthAccount{}).Where("account_id = ?", accountID).Update("is_default", true).Error
	})
}

func DeleteCodexOAuthAccount(accountID string) error {
	return DB.Where("account_id = ?", accountID).Delete(&CodexOAuthAccount{}).Error
}

var codexOAuthProtectRefreshToken = func(token string) string {
	return token
}

func SetCodexOAuthRefreshTokenProtector(protector func(string) string) {
	if protector == nil {
		codexOAuthProtectRefreshToken = func(token string) string { return token }
		return
	}
	codexOAuthProtectRefreshToken = protector
}
