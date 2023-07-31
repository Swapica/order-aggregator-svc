package notifications

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (n *NotificationsClient) getNotificationPayload(title, body, msgType string) NotificationPayload {
	return NotificationPayload{
		Notification: Notification{
			Title: title,
			Body:  body,
		},
		Data: Data{
			AMsg: body,
			ASub: title,
			Type: msgType,
		},
	}
}

func (n *NotificationsClient) getVerificationProof(payload NotificationPayload) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal notification payload to JSON")
	}

	message := "2+" + string(body)

	signature, err := SignTypedData(apitypes.TypedData{
		PrimaryType: "Data",
		Types: apitypes.Types{
			"Data": []apitypes.Type{
				{Name: "data", Type: "string"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
		},
		Domain: apitypes.TypedDataDomain{
			Name:              "EPNS COMM V1",
			ChainId:           math.NewHexOrDecimal256(n.chainId),
			VerifyingContract: n.pushCommAddress,
		},
		Message: apitypes.TypedDataMessage{
			"data": message,
		},
	}, n.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign typed data")
	}

	uuidInstance, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate uuid")
	}

	proof := "eip712v2:" + "0x" + hex.EncodeToString(signature) + "::uid::" + uuidInstance.String()

	return proof, nil
}

func getSourceFromChainId(chainId int64) string {
	switch chainId {
	case 1:
		return "ETH_MAINNET"
	case 5:
		return "ETH_TEST_GOERLI"
	case 137:
		return "POLYGON_MAINNET"
	case 80001:
		return "POLYGON_TEST_MUMBAI"
	case 56:
		return "BSC_MAINNET"
	case 97:
		return "BSC_TESTNET"
	case 10:
		return "OPTIMISM_MAINNET"
	case 420:
		return "OPTIMISM_TESTNET"
	case 1442:
		return "POLYGON_ZK_EVM_TESTNET"
	case 1101:
		return "POLYGON_ZK_EVM_MAINNET"
	}
	return ""
}

func addressToCAIP(addr string, chainId int64) string {
	return "eip155:" + strconv.FormatInt(chainId, 10) + ":" + addr
}
