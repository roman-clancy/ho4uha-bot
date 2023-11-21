package messages

import "github.com/roman-clancy/ho4uha-bot/internal/model/types"

type WishItem struct {
	Name string
	URL  string
}

type UserStorage interface {
	AddNewUser(userId int64) (bool, error)
	AddUserCategory(userId int64, catName string) (bool, error)
	AddWishItem(userId int64, item WishItem) (bool, error)
	AddWishItemToCategory(userId int64, catName string, item WishItem) (bool, error)
	GetWishListByCategory(userId int64) map[string][]WishItem
	GetCategories(userId int64) []string
}

type MessageSender interface {
	SendMessage(userId int64, text string) error
	ShowButtons(userId int64, text string, buttons []types.TgRowButtons) error
}

type Message struct {
	Text          string
	UserID        int64
	UserName      string
	IsCallback    bool
	CallbackMsgID string
}

type BotModel struct {
	UserStorage      UserStorage
	MessageSender    MessageSender
	lastUserCmd      map[int64]string
	lastUserCat      map[int64]string
	lastUserItemName map[int64]string
}
