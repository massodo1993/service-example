package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		partUUID = uuid.New()
		expected = model.Part{UUID: partUUID, Name: gofakeit.Name(), Price: gofakeit.Price(10, 1000)}
	)

	s.partRepository.On("GetPart", s.ctx, partUUID.String()).Return(expected, nil)

	part, err := s.service.GetPart(s.ctx, partUUID.String())
	s.Require().NoError(err)
	s.Require().Equal(expected, part)
}

func (s *ServiceSuite) TestGetPartRepoError() {
	var (
		partUUID = uuid.New().String()
		repoErr  = gofakeit.Error()
	)

	s.partRepository.On("GetPart", s.ctx, partUUID).Return(model.Part{}, repoErr)

	part, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(part)
}
