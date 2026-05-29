package codexoauth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
	"github.com/songquanpeng/one-api/relay/relaymode"
	"github.com/songquanpeng/one-api/services/codexoauth"
)

var _ adaptor.Adaptor = new(Adaptor)

const (
	channelName     = "codex_oauth"
	defaultBaseURL  = "https://chatgpt.com/backend-api/codex"
	responsesPath   = "/responses"
	originatorValue = "one-api"
)

var modelList = []string{"gpt-5-codex"}

type Adaptor struct{}

func (a *Adaptor) Init(meta *meta.Meta) {}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	baseURL := strings.TrimRight(meta.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return baseURL + responsesPath, nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	accountID := meta.Config.ManagedAccountIDFor("codex_oauth")
	token, err := codexoauth.DefaultManager.GetValidTokenForAccount(accountID)
	if err != nil {
		return err
	}
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("ChatGPT-Account-Id", token.AccountID)
	req.Header.Set("originator", originatorValue)
	setSessionHeaders(c, req, meta)
	return nil
}

func setSessionHeaders(c *gin.Context, req *http.Request, meta *meta.Meta) {
	sessionID := strings.TrimSpace(c.Request.Header.Get("session_id"))
	if sessionID == "" {
		sessionID = strings.TrimSpace(c.Request.Header.Get("x-client-request-id"))
	}
	if sessionID == "" {
		sessionID = strings.TrimSpace(c.GetString("id"))
	}
	if sessionID == "" {
		return
	}
	req.Header.Set("session_id", sessionID)
	req.Header.Set("x-client-request-id", sessionID)
	req.Header.Set("x-codex-window-id", sessionID+":0")
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	if relayMode != relaymode.Responses {
		return nil, fmt.Errorf("codex oauth only supports responses relay mode, got %d", relayMode)
	}
	return request, nil
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	return nil, errors.New("codex oauth image requests are not supported")
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Set(k, vv)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)
	if _, copyErr := io.Copy(c.Writer, resp.Body); copyErr != nil {
		return nil, &model.ErrorWithStatusCode{
			StatusCode: http.StatusInternalServerError,
			Error: model.Error{
				Message: copyErr.Error(),
			},
		}
	}
	return &model.Usage{}, nil
}

func (a *Adaptor) GetModelList() []string {
	return modelList
}

func (a *Adaptor) GetChannelName() string {
	return channelName
}
