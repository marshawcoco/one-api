package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/model"
	"github.com/songquanpeng/one-api/services/codexoauth"
)

type codexOAuthStatusResponse struct {
	Provider         string                     `json:"provider"`
	Authenticated    bool                       `json:"authenticated"`
	DefaultAccountID string                     `json:"default_account_id,omitempty"`
	Accounts         []*model.CodexOAuthAccount `json:"accounts"`
}

func GetCodexOAuthStatus(c *gin.Context) {
	accounts, err := model.GetCodexOAuthAccounts()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	defaultAccountID := ""
	for _, account := range accounts {
		if account.IsDefault {
			defaultAccountID = account.AccountID
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": codexOAuthStatusResponse{
			Provider:         model.CodexOAuthProvider,
			Authenticated:    len(accounts) > 0,
			DefaultAccountID: defaultAccountID,
			Accounts:         accounts,
		},
	})
}

func ListCodexOAuthAccounts(c *gin.Context) {
	accounts, err := model.GetCodexOAuthAccounts()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": accounts})
}

func StartCodexOAuthLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"message": "Codex OAuth login requires a configured OAuth exchange backend and is not implemented in this build",
	})
}

func PollCodexOAuthAccount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"message": "Codex OAuth polling requires a configured OAuth exchange backend and is not implemented in this build",
	})
}

func SetDefaultCodexOAuthAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "account_id is required"})
		return
	}
	if _, err := model.GetCodexOAuthAccountByAccountID(accountID); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := model.SetDefaultCodexOAuthAccount(accountID); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func DeleteCodexOAuthAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "account_id is required"})
		return
	}
	if err := model.DeleteCodexOAuthAccount(accountID); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	codexoauth.DefaultManager.ClearAccountToken(accountID)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}
