package notifications

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/config"
)

type NotificationsClient struct {
	client          *http.Client
	baseUrl         string
	chainId         int64
	channelAddress  string
	pushCommAddress string
	privateKey      *ecdsa.PrivateKey
	source          string
}

func NewNotificationsClient(cfg config.Notifications, chainId int64) *NotificationsClient {
	return &NotificationsClient{
		client:          http.DefaultClient,
		baseUrl:         cfg.PushURL,
		chainId:         chainId,
		source:          getSourceFromChainId(chainId),
		channelAddress:  addressToCAIP(cfg.ChannelAddress, chainId),
		pushCommAddress: cfg.PushCommAddress,
		privateKey:      cfg.PrivateKey,
	}
}
