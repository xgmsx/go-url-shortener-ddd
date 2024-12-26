package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
)

// Database stub

type fakeDatabase struct {
	Storage []entity.Link
	Err     error
}

func newFakeDatabase(storage []entity.Link) *fakeDatabase {
	if storage != nil {
		return &fakeDatabase{Storage: storage}
	}
	return &fakeDatabase{Storage: []entity.Link{}}
}

func (f *fakeDatabase) CreateLink(_ context.Context, l entity.Link) error {
	if f.Err != nil {
		return f.Err
	}

	for _, v := range f.Storage {
		if v == l {
			return nil
		}
	}
	f.Storage = append(f.Storage, l)
	return nil
}

func (f *fakeDatabase) GetLink(_ context.Context, alias, url string) (entity.Link, error) {
	if f.Err != nil {
		return entity.Link{}, f.Err
	}

	for _, v := range f.Storage {
		if alias != "" && alias == v.Alias {
			return v, nil
		}
		if url != "" && url == v.URL {
			return v, nil
		}
	}
	return entity.Link{}, entity.ErrNotFound
}

// Redis stub

type fakeCache struct {
	Storage []entity.Link
	Err     error
}

func newFakeCache(storage []entity.Link) *fakeCache {
	if storage != nil {
		return &fakeCache{Storage: storage}
	}
	return &fakeCache{Storage: []entity.Link{}}
}

func (f *fakeCache) PutLink(_ context.Context, l entity.Link) error {
	if f.Err != nil {
		return f.Err
	}

	for _, v := range f.Storage {
		if v == l {
			return nil
		}
	}
	f.Storage = append(f.Storage, l)
	return nil
}

func (f *fakeCache) GetLink(_ context.Context, alias string) (entity.Link, error) {
	if f.Err != nil {
		return entity.Link{}, f.Err
	}

	for _, v := range f.Storage {
		if v.Alias == alias {
			return v, nil
		}
	}
	return entity.Link{}, entity.ErrNotFound
}

// Kafka stub

type fakeBroker struct {
	Storage []entity.Link
	Err     error
}

func newFakeBroker(storage []entity.Link) *fakeBroker {
	if storage != nil {
		return &fakeBroker{Storage: storage}
	}
	return &fakeBroker{Storage: []entity.Link{}}
}

func (f *fakeBroker) CreateEvent(_ context.Context, l entity.Link) error {
	if f.Err != nil {
		return f.Err
	}

	f.Storage = append(f.Storage, l)
	return nil
}

func TestNew(t *testing.T) {
	// Arrange
	var (
		db     = &fakeDatabase{}
		cache  = &fakeCache{}
		broker = &fakeBroker{}
	)

	// Act
	u := usecase.New(db, cache, broker)

	// Assert
	assert.NotNil(t, u)
}
