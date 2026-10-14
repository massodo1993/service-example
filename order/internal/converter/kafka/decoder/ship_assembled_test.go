package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	eventsV1 "github.com/massodo1993/service-example/shared/pkg/proto/events/v1"
)

func TestShipAssembledDecodeSuccess(t *testing.T) {
	data, err := proto.Marshal(&eventsV1.ShipAssembled{EventUuid: "event", OrderUuid: "order", UserUuid: "user", BuildTimeSec: 7})
	require.NoError(t, err)

	got, err := NewShipAssembledDecoder().Decode(data)
	require.NoError(t, err)
	require.Equal(t, "order", got.OrderUUID)
	require.Equal(t, "user", got.UserUUID)
	require.Equal(t, int64(7), got.BuildTimeSec)
}

func TestShipAssembledDecodeInvalidData(t *testing.T) {
	_, err := NewShipAssembledDecoder().Decode([]byte{0xff, 0xff, 0xff})
	require.Error(t, err)
}
