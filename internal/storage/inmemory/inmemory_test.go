package inmemory

import (
	"github.com/roman-clancy/ho4uha-bot/internal/model/messages"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func (s *Storage) getUserData(userId int64) *UserData {
	return s.users[userId]
}

func TestStorage_AddNewUser(t *testing.T) {
	storage, err := New()
	require.NoError(t, err)

	t.Run("Should add new user", func(t *testing.T) {
		userId := int64(1)
		inserted, err := storage.AddNewUser(userId)
		require.True(t, inserted)
		require.NoError(t, err)
		data := storage.getUserData(userId)
		require.Equalf(t, userId, data.userId, "Should store user with id = %d", userId)
	})

	t.Run("Should add new user with default category", func(t *testing.T) {
		userId := int64(2)
		expectedCategoryCnt := 1
		inserted, err := storage.AddNewUser(userId)
		require.True(t, inserted, "New user should be marked as 'inserted'")
		require.NoError(t, err)

		data := storage.getUserData(userId)
		require.Equalf(t, userId, data.userId, "Should store user with id = %d", userId)
		require.Equalf(t, expectedCategoryCnt, len(data.categories), "New user should have %d category", expectedCategoryCnt)
		require.Equalf(t, "default", data.categories[0].name, "Category should have '%s' name", "default")
	})

	t.Run("Should not add new user twice", func(t *testing.T) {
		userId := int64(3)
		inserted, err := storage.AddNewUser(userId)
		require.True(t, inserted)
		require.NoError(t, err)
		inserted, err = storage.AddNewUser(userId)
		require.Falsef(t, inserted, "Insert user with existing Id should return 'false' flag")
		data := storage.getUserData(userId)
		require.Equalf(t, userId, data.userId, "Should store user with id = %d", userId)
	})
}

func TestStorage_AddUserCategory(t *testing.T) {
	storage, err := New()
	userId := int64(1)
	inserted, err := storage.AddNewUser(userId)
	require.True(t, inserted)
	require.NoError(t, err)
	t.Run("Should add new category with empty wishlist", func(t *testing.T) {
		newCategoryName := "Board Games"
		expectedCategorySize := 2
		added, err := storage.AddUserCategory(userId, newCategoryName)
		require.NoError(t, err, "Adding new category shouldn't cause error")
		require.True(t, added, "Adding new category should return 'true' flag")
		data := storage.getUserData(userId)
		categoriesSize := len(data.categories)
		require.Equalf(t, expectedCategorySize, categoriesSize, "Expected categories new size is %d", expectedCategorySize)
		newCategoryIndex := slices.IndexFunc(data.categories, func(category *Category) bool {
			return category.name == newCategoryName
		})
		require.NotEqualf(t, -1, newCategoryIndex, "IndexFunc should find new category")
		category := data.categories[newCategoryIndex]
		require.Equalf(t, newCategoryName, category.name, "New category name should be %s, got %s", newCategoryName, category.name)
		wishlistSize := len(category.items)
		require.Equalf(t, 0, wishlistSize, "New category should have empty wishlist, got %d", wishlistSize)
	})

	t.Run("Shouldn't add new category when user doesn't exists", func(t *testing.T) {
		newCategoryName := "Board Games"
		notExistingUserId := int64(2)
		added, err := storage.AddUserCategory(notExistingUserId, newCategoryName)
		require.NoError(t, err, "Adding new category shouldn't cause error")
		require.Falsef(t, added, "Adding new category should return 'false' flag when user doesn't exists")
	})
}

func TestStorage_AddWishItem(t *testing.T) {
	storage, err := New()
	userId := int64(1)
	notExistingUserId := int64(2)
	inserted, err := storage.AddNewUser(userId)
	require.True(t, inserted)
	require.NoError(t, err)
	t.Run("Should add new wish item to default wishlist", func(t *testing.T) {
		wishItemName := "Wish Item Name"
		wishItemUrl := "Wish URL"
		newItem := messages.WishItem{
			Name: wishItemName,
			URL:  wishItemUrl,
		}
		added, err := storage.AddWishItem(userId, newItem)
		require.NoErrorf(t, err, "Add new wish item shouldn't cause error")
		require.True(t, added, "Add new wish item should return 'true' flag")
		data := storage.getUserData(userId)
		defaultCategory := data.categories[0]
		expectedWishlistSize := 1
		require.Equalf(t, expectedWishlistSize, len(defaultCategory.items), "Default category should have wishlist with size = %d")
		wishItem := defaultCategory.items[0]
		require.Equalf(t, wishItemName, wishItem.Name, "New item should have Name = %s, got %s", wishItemName, wishItem.Name)
		require.Equalf(t, wishItemUrl, wishItem.URL, "New item should have URL = %s, got %s", wishItemUrl, wishItem.URL)
	})

	t.Run("Shouldn't add new wish item to default wishlist when user doesn't exists", func(t *testing.T) {
		wishItemName := "Wish Item Name"
		wishItemUrl := "Wish URL"
		newItem := messages.WishItem{
			Name: wishItemName,
			URL:  wishItemUrl,
		}
		added, err := storage.AddWishItem(notExistingUserId, newItem)
		require.NoErrorf(t, err, "Add new wish item shouldn't cause error")
		require.Falsef(t, added, "Add new wish item should return 'false' flag when user doesn't exists")
	})
}

func TestStorage_AddWishItemToCategory(t *testing.T) {
	storage, err := New()
	userId := int64(1)
	// notExistingUserId := int64(2)
	categoryName := "Table Games"
	// notExistingCategory := "Not exists"

	inserted, err := storage.AddNewUser(userId)
	require.True(t, inserted)
	require.NoError(t, err)

	inserted, err = storage.AddUserCategory(userId, categoryName)
	require.True(t, inserted)
	require.NoError(t, err)

	t.Run("Should add new wish item to specific category", func(t *testing.T) {
		wishItemName := "Wish Item Name"
		wishItemUrl := "Wish URL"
		newItem := messages.WishItem{
			Name: wishItemName,
			URL:  wishItemUrl,
		}
		added, err := storage.AddWishItemToCategory(userId, categoryName, newItem)
		require.NoErrorf(t, err, "Add new wish item shouldn't cause error")
		require.True(t, added, "Add new wish item should return 'true' flag")
		data := storage.getUserData(userId)
		idx := slices.IndexFunc(data.categories, func(category *Category) bool {
			return category.name == categoryName
		})
		require.NotEqualf(t, -1, idx, "Should find category '%s'", categoryName)
		category := data.categories[idx]
		expectedWishlistSize := 1
		require.Equalf(t, expectedWishlistSize, len(category.items), "Category should have wishlist with size = %d")
		wishItem := category.items[0]
		require.Equalf(t, wishItemName, wishItem.Name, "New item should have Name = %s, got %s", wishItemName, wishItem.Name)
		require.Equalf(t, wishItemUrl, wishItem.URL, "New item should have URL = %s, got %s", wishItemUrl, wishItem.URL)
	})

}

func TestStorage_GetWishListByCategory(t *testing.T) {
	storage, err := New()
	userId := int64(1)
	inserted, err := storage.AddNewUser(userId)
	require.True(t, inserted)
	require.NoError(t, err)

	t.Run("", func(t *testing.T) {
		storage.GetWishListByCategory(userId)
	})
}
