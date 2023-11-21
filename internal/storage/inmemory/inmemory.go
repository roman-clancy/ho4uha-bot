package inmemory

import (
	"github.com/roman-clancy/ho4uha-bot/internal/model/messages"
)

type UserData struct {
	userId     int64
	categories []*Category
}

type Category struct {
	name  string
	items []messages.WishItem
}

type Storage struct {
	users map[int64]*UserData
}

func New() (*Storage, error) {
	return &Storage{
		users: make(map[int64]*UserData),
	}, nil
}

func (s *Storage) AddNewUser(userId int64) (bool, error) {
	_, ok := s.users[userId]
	if !ok {
		userData := &UserData{
			userId:     userId,
			categories: make([]*Category, 0),
		}
		userData.categories = append(userData.categories,
			&Category{
				name:  "default",
				items: make([]messages.WishItem, 0),
			})
		s.users[userId] = userData
		return true, nil
	}
	return false, nil
}

func (s *Storage) AddUserCategory(userId int64, catName string) (bool, error) {
	data, ok := s.users[userId]
	if ok {
		data.categories = append(data.categories, &Category{
			name:  catName,
			items: make([]messages.WishItem, 0),
		})
		return true, nil
	}
	return false, nil
}

func (s *Storage) AddWishItemToCategory(userId int64, catName string, item messages.WishItem) (bool, error) {
	if data, ok := s.users[userId]; ok {
		for _, cat := range data.categories {
			if cat.name == catName {
				cat.items = append(cat.items, item)
			}
		}
		return true, nil
	}
	return false, nil
}

func (s *Storage) AddWishItem(userId int64, item messages.WishItem) (bool, error) {
	return s.AddWishItemToCategory(userId, "default", item)
}

func (s *Storage) GetWishListByCategory(userId int64) map[string][]messages.WishItem {
	userData, ok := s.users[userId]
	if ok {
		result := make(map[string][]messages.WishItem)
		for _, category := range userData.categories {
			var a = make([]messages.WishItem, len(category.items))
			copy(a, category.items)
			result[category.name] = a
		}
		return result
	}
	return nil
}

func (s *Storage) GetCategories(userId int64) []string {
	data, ok := s.users[userId]
	result := make([]string, 0, 10)
	if ok {
		for _, cat := range data.categories {
			if cat.name != "default" {
				result = append(result, cat.name)
			}
		}
	}
	return result
}
