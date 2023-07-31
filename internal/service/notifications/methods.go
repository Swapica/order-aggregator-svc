package notifications

import (
	"encoding/json"
	"io"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (n *NotificationsClient) NotifyUser(msgTitle, msgBody, address string) error {
	payload := n.getNotificationPayload(msgTitle, msgBody, "3")
	payload.Recipients = addressToCAIP(address, n.chainId)

	proof, err := n.getVerificationProof(payload)
	if err != nil {
		return errors.Wrap(err, "failed to get verification proof")
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal notification payload to JSON")
	}

	apiPayload := ApiPayload{
		VerificationProof: proof,
		Identity:          "2+" + string(payloadJson),
		Sender:            n.channelAddress,
		Source:            n.source,
		Recipient:         payload.Recipients,
	}

	body, err := json.Marshal(apiPayload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal notification payload to JSON")
	}

	resp, err := n.post("/v1/payloads/", body)
	if err != nil {
		return errors.Wrap(err, "failed to send notification")
	}
	if resp.StatusCode > 204 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.New(string(b))
	}
	return nil
}

func (n *NotificationsClient) NotifyUsers(msgTitle, msgBody string, addresses []string) error {
	// TODO
	return nil
}

func (n *NotificationsClient) NotifyAll(msgTitle, msgBody string) error {
	payload := n.getNotificationPayload(msgTitle, msgBody, "1")
	payload.Recipients = n.channelAddress

	proof, err := n.getVerificationProof(payload)
	if err != nil {
		return errors.Wrap(err, "failed to get verification proof")
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal notification payload to JSON")
	}

	apiPayload := ApiPayload{
		VerificationProof: proof,
		Identity:          "2+" + string(payloadJson),
		Sender:            n.channelAddress,
		Source:            n.source,
		Recipient:         n.channelAddress, // channel address for broadcast notification type
	}

	body, err := json.Marshal(apiPayload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal notification payload to JSON")
	}

	resp, err := n.post("/v1/payloads/", body)
	if err != nil {
		return errors.Wrap(err, "failed to send notification")
	}
	if resp.StatusCode > 204 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.New(string(b))
	}
	return nil
}
