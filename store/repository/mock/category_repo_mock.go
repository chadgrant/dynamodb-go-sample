package mock

import "sort"

type MockCategoryRepository struct {
	categories map[string]interface{}
}

func NewCategoryRepository(categories ...string) *MockCategoryRepository {
	r := &MockCategoryRepository{
		categories: make(map[string]interface{}),
	}
	for _, c := range categories {
		if err := r.Upsert(c); err != nil {
			panic(err)
		}
	}
	return r
}

func (r *MockCategoryRepository) GetAll() ([]string, error) {
	arr := make([]string, 0)
	for k := range r.categories {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return arr, nil
}

func (r *MockCategoryRepository) Upsert(category string) error {
	r.categories[category] = nil
	return nil
}

func (r *MockCategoryRepository) Delete(category string) error {

	delete(r.categories, category)

	return nil
}
