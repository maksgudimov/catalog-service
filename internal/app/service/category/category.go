package scategory

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

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Category {
	return &srv{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error) {
	existing, err := s.repoCategory.List(ctx, &req.Name)
	if err != nil {
		return entity.Category{}, err
	}
	if len(existing) > 0 {
		return entity.Category{}, entity.ErrAlreadyExists
	}

	now := time.Now()
	category := entity.Category{
		GUID:      uuid.Must(uuid.NewV4()),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repoCategory.Create(ctx, category); err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error) {
	categories, err := s.repoCategory.GetByGUIDs(ctx, guids)
	if err != nil {
		return []entity.Category{}, err
	}
	if len(categories) == 0 {
		return []entity.Category{}, nil
	}
	return categories, nil
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error) {
	category, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return entity.Category{}, err
	}
	if len(category) == 0 {
		return entity.Category{}, entity.ErrNotFound
	}
	categories, err := s.repoCategory.List(ctx, &req.Name)
	if err != nil {
		return entity.Category{}, err
	}

	for _, val := range categories {
		if val.GUID != guid {
			return entity.Category{}, entity.ErrAlreadyExists
		}
	}

	category[0].Name = req.Name

	err = s.repoCategory.Update(ctx, category[0])
	if err != nil {
		return entity.Category{}, err
	}

	return category[0], nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	category, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return err
	}
	if len(category) == 0 {
		return entity.ErrNotFound
	}
	products, err := s.repoProduct.List(ctx, nil, &guid)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return entity.ErrCategoryHasProducts
	}

	err = s.repoCategory.Delete(ctx, guid)
	if err != nil {
		return err
	}

	return nil
}

func (s *srv) List(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repoCategory.List(ctx, nil)
	if err != nil {
		return []entity.Category{}, err
	}
	if len(categories) == 0 {
		return []entity.Category{}, nil
	}
	return categories, err
}
