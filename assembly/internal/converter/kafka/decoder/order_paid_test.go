package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	eventsV1 "github.com/massodo1993/service-example/shared/pkg/proto/events/v1"
)

func TestOrderPaidDecodeSuccess(t *testing.T) {
	data, err := proto.Marshal(&eventsV1.OrderPaid{EventUuid: "event", OrderUuid: "order", UserUuid: "user", PaymentMethod: "CARD", TransactionUuid: "tx"})
	require.NoError(t, err)

	got, err := NewOrderPaidDecoder().Decode(data)
	require.NoError(t, err)
	require.Equal(t, "order", got.OrderUUID)
	require.Equal(t, "user", got.UserUUID)
	require.Equal(t, "CARD", got.PaymentMethod)
	require.Equal(t, "tx", got.TransactionUUID)
}

func TestOrderPaidDecodeInvalidData(t *testing.T) {
	_, err := NewOrderPaidDecoder().Decode([]byte{0xff, 0xff, 0xff})
	require.Error(t, err)
}
