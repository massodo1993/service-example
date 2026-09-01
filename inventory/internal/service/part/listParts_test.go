package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		filter = model.PartsFilter{UUIDs: []string{uuid.New().String()}}
		parts  = []model.Part{
			{UUID: uuid.New(), Name: gofakeit.Name()},
			{UUID: uuid.New(), Name: gofakeit.Name()},
		}
	)

	s.partRepository.On("ListParts", s.ctx, filter).Return(parts, nil)

	got, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(parts, got)
}

func (s *ServiceSuite) TestListPartsRepoError() {
	var (
		filter  = model.PartsFilter{}
		repoErr = gofakeit.Error()
	)

	s.partRepository.On("ListParts", s.ctx, filter).Return(nil, repoErr)

	got, err := s.service.ListParts(s.ctx, filter)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(got)
}
