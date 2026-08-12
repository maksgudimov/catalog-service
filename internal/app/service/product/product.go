package sproduct

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/maksgudimov/catalog-service/internal/app/entity"
	"github.com/maksgudimov/catalog-service/internal/app/repository"
	"github.com/maksgudimov/catalog-service/internal/app/service"
)

type srv struct {
	repoCategory repository.Category
	repoProduct  repository.Product
}

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Product {
	return &srv{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	existing, err := s.repoProduct.List(ctx, &req.Name, nil)
	if err != nil {
		return entity.Product{}, err
	}
	if len(existing) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}

	category, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{req.CategoryGUID})
	if err != nil {
		return entity.Product{}, err
	}
	if len(category) == 0 {
		return entity.Product{}, entity.ErrNotFound
	}

	now := time.Now()
	product := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repoProduct.Create(ctx, product); err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *srv) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	products, err := s.repoProduct.GetByGUIDs(ctx, guids)
	if err != nil {
		return []entity.Product{}, err
	}
	if len(products) == 0 {
		return []entity.Product{}, nil
	}
	return products, nil
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return entity.Product{}, err
	}
	if len(product) == 0 {
		return entity.Product{}, entity.ErrNotFound
	}

	if req.Name != nil {
		existing, err := s.repoProduct.List(ctx, req.Name, nil)
		if err != nil {
			return entity.Product{}, err
		}

		for _, val := range existing {
			if val.GUID != guid {
				return entity.Product{}, entity.ErrAlreadyExists
			}
		}
		product[0].Name = *req.Name

	}

	if req.CategoryGUID != nil {
		category, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{*req.CategoryGUID})
		if err != nil {
			return entity.Product{}, err
		}
		if len(category) == 0 {
			return entity.Product{}, entity.ErrNotFound
		}
		product[0].CategoryGUID = *req.CategoryGUID

	}

	if req.Price != nil {
		product[0].Price = *req.Price
	}

	if req.Description != nil {
		product[0].Description = req.Description
	}

	now := time.Now()
	product[0].UpdatedAt = now

	err = s.repoProduct.Update(ctx, product[0])
	if err != nil {
		return entity.Product{}, err
	}

	return product[0], nil

}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	product, err := s.repoProduct.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return err
	}
	if len(product) == 0 {
		return entity.ErrNotFound
	}

	err = s.repoProduct.Delete(ctx, guid)
	if err != nil {
		return err
	}

	return nil

}

func (s *srv) List(ctx context.Context, req entity.RequestProductList) ([]entity.Product, error) {
	products, err := s.repoProduct.List(ctx, nil, req.CategoryGUID)
	if err != nil {
		return []entity.Product{}, err
	}
	if len(products) == 0 {
		return []entity.Product{}, nil
	}
	return products, err
}
