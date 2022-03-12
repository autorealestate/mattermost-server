package wsapi

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v6/app"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

func (api *API) InitSubscription() {
	api.Router.Handle("subscribe", api.APIWebSocketHandler(api.subscribe))
	api.Router.Handle("unsubscribe", api.APIWebSocketHandler(api.unsubscribe))
}

func subscriptionIDFromRequest(req *model.WebSocketRequest) (model.WebsocketSubscriptionID, *model.AppError) {
	const paramKey = "subscription_id"

	paramVal, has := req.Data[paramKey]
	if !has {
		mlog.Debug("missing JSON field", mlog.String("key", paramKey))
		return "", NewInvalidWebSocketParamError(req.Action, paramKey)
	}

	subscriptionID := model.WebsocketSubscriptionID(fmt.Sprintf("%s", paramVal))
	if !subscriptionID.IsValid() {
		mlog.Debug("invalid JSON field", mlog.String("key", paramKey), mlog.Any("value", paramVal))
		return "", NewInvalidWebSocketParamError(req.Action, paramKey)
	}

	return subscriptionID, nil
}

func (api *API) subscribe(req *model.WebSocketRequest, conn *app.WebConn) (map[string]interface{}, *model.AppError) {
	subscriptionID, err := subscriptionIDFromRequest(req)
	if err != nil {
		return nil, err
	}

	conn.Subscribe(subscriptionID)

	return nil, nil
}

func (api *API) unsubscribe(req *model.WebSocketRequest, conn *app.WebConn) (map[string]interface{}, *model.AppError) {
	subscriptionID, err := subscriptionIDFromRequest(req)
	if err != nil {
		return nil, err
	}

	conn.Unsubscribe(subscriptionID)

	return nil, nil
}
