package codexoauth

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
	"github.com/songquanpeng/one-api/services/codexoauth"
)

var _ adaptor.Adaptor = new(Adaptor)

const channelName = "codex_oauth"

var modelList = []string{"gpt-5-codex"}

type Adaptor struct{}

func (a *Adaptor) Init(meta *meta.Meta) {}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	return "", errors.New("codex oauth adaptor is not implemented")
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	accountID := meta.Config.ManagedAccountIDFor("codex_oauth")
	accessToken, err := codexoauth.DefaultManager.GetValidTokenForAccount(accountID)
	if err != nil {
		return err
	}
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("ChatGPT-Account-Id", accountID)
	req.Header.Set("originator", "one-api")
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	return nil, errors.New("codex oauth adaptor is not implemented")
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	return nil, errors.New("codex oauth adaptor is not implemented")
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return nil, errors.New("codex oauth adaptor is not implemented")
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	return nil, &model.ErrorWithStatusCode{
		StatusCode: http.StatusNotImplemented,
		Error: model.Error{
			Message: "codex oauth adaptor is not implemented",
		},
	}
}

func (a *Adaptor) GetModelList() []string {
	return modelList
}

func (a *Adaptor) GetChannelName() string {
	return channelName
}
