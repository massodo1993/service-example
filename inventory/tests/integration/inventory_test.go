//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			uuid, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
			partUUID = uuid.String()
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(Equal(partUUID))
			Expect(resp.GetPart().GetName()).ToNot(BeEmpty())
		})

		It("должен возвращать ошибку для несуществующего UUID", func() {
			_, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: "00000000-0000-0000-0000-000000000000",
			})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("ListParts", func() {
		var partUUID string

		BeforeEach(func() {
			uuid, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
			partUUID = uuid.String()
		})

		It("должен возвращать деталь по фильтру uuids", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{partUUID},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(1))
			Expect(resp.GetParts()[0].GetUuid()).To(Equal(partUUID))
		})

		It("должен возвращать пустой список для несуществующего фильтра", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{"00000000-0000-0000-0000-000000000000"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(BeEmpty())
		})
	})
})
