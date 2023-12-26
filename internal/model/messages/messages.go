package messages

import (
	"fmt"
	"github.com/roman-clancy/ho4uha-bot/internal/model/types"
	"strings"
)

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

var btnStart = []types.TgRowButtons{
	{
		types.TgInlineButton{DisplayName: "Добавить категорию", Value: "/add_cat"},
		types.TgInlineButton{DisplayName: "Добавить xотелку", Value: "/add_item"},
	},
	{
		types.TgInlineButton{DisplayName: "Показать мои категории", Value: "/show_cat"},
		types.TgInlineButton{DisplayName: "Показать мои хотелки", Value: "/show_item"},
	},
}
var cancelBtn = []types.TgRowButtons{
	{types.TgInlineButton{DisplayName: "Отмена", Value: "/cancel"}},
}

const (
	txtStart          = "Привет. Я могу помочь составить Вишлист и поделиться им с друзьями. Выберите действие."
	txtChooseCmd      = "Выберите действие."
	txtUnknownCommand = "К сожалению, данная команда мне неизвестна. Для начала работы введите /start"
	txtCatAdd         = "Введите название категории"
	txtItemAdd        = "Введите название хотелки"
	txtItemUrl        = "Добавьте ссылку на вашу хотелку"
	txtItemShow       = "Ваши хотелки:"
	txtCatChoose      = "Выберите категорию хотелки"
	txtCatShow        = "Ваши категории:"
	txtCatShowErr     = "Ошибка при формировании списка категорий"
	txtAddDone        = "Сохранение успешно"
)

func New(userStorage UserStorage, sender MessageSender) *BotModel {
	return &BotModel{
		UserStorage:      userStorage,
		MessageSender:    sender,
		lastUserCmd:      map[int64]string{},
		lastUserCat:      map[int64]string{},
		lastUserItemName: map[int64]string{},
	}
}

func (m *BotModel) OnMessage(msg Message) error {
	lastUserCmd := m.lastUserCmd[msg.UserID]
	m.lastUserCmd[msg.UserID] = ""
	if isNeedReturn, err := checkNewCategoryAdded(m, msg, lastUserCmd); isNeedReturn || err != nil {
		return err
	}
	if isNeedReturn, err := checkNewItemAdded(m, msg); isNeedReturn || err != nil {
		return err
	}
	if isNeedReturn, err := checkNewItemNameAdded(m, msg); isNeedReturn || err != nil {
		return err
	}
	if isNeedReturn, err := checkNewItemUrlAdded(m, msg); isNeedReturn || err != nil {
		return err
	}
	if isNeedReturn, err := checkBotCommands(m, msg); isNeedReturn || err != nil {
		return err
	}
	return m.MessageSender.SendMessage(msg.UserID, txtUnknownCommand)
}

func checkNewCategoryAdded(m *BotModel, msg Message, lastCmd string) (bool, error) {
	if lastCmd == "/add_cat" {
		_, err := m.UserStorage.AddUserCategory(msg.UserID, msg.Text)
		if err != nil {
			return true, err
		}
		return true, m.MessageSender.ShowButtons(msg.UserID, txtAddDone, btnStart)
	}
	return false, nil
}

func checkNewItemAdded(m *BotModel, msg Message) (bool, error) {
	if msg.IsCallback {
		if strings.Contains(msg.Text, "/cat ") {
			cat := strings.Replace(msg.Text, "/cat ", "", -1)
			m.lastUserCat[msg.UserID] = cat
			return true, m.MessageSender.SendMessage(msg.UserID, txtItemAdd)
		}
	}
	return false, nil
}

func checkNewItemNameAdded(m *BotModel, msg Message) (bool, error) {
	if m.lastUserCat[msg.UserID] != "" && m.lastUserItemName[msg.UserID] == "" && msg.Text != "" {
		m.lastUserItemName[msg.UserID] = msg.Text
		return true, m.MessageSender.SendMessage(msg.UserID, txtItemUrl)
	}
	return false, nil
}

func checkNewItemUrlAdded(m *BotModel, msg Message) (bool, error) {
	if m.lastUserCat[msg.UserID] != "" && m.lastUserItemName[msg.UserID] != "" {
		cat := m.lastUserCat[msg.UserID]
		m.lastUserCat[msg.UserID] = ""
		itemName := m.lastUserItemName[msg.UserID]
		m.lastUserItemName[msg.UserID] = ""
		_, err := m.UserStorage.AddWishItemToCategory(msg.UserID, cat, WishItem{
			Name: itemName,
			URL:  msg.Text,
		})
		if err != nil {
			return true, err
		}
		return true, m.MessageSender.ShowButtons(msg.UserID, txtAddDone, btnStart)
	}
	return false, nil
}

func checkBotCommands(model *BotModel, msg Message) (bool, error) {
	switch msg.Text {
	case "/start":
		if _, err := model.UserStorage.AddNewUser(msg.UserID); err != nil {
			return true, err
		}
		return true, model.MessageSender.ShowButtons(msg.UserID, txtStart, btnStart)
	case "/add_cat":
		model.lastUserCmd[msg.UserID] = "/add_cat"
		return true, model.MessageSender.ShowButtons(msg.UserID, txtCatAdd, cancelBtn)
	case "/add_item":
		model.lastUserCmd[msg.UserID] = "/add_item"
		var categoryButtons = getCategoryButtons(model.UserStorage.GetCategories(msg.UserID))
		return true, model.MessageSender.ShowButtons(msg.UserID, txtCatChoose, categoryButtons)
	case "/show_cat":
		categoriesString, err := getCategoryList(model, msg.UserID)
		if err != nil {
			return true, model.MessageSender.SendMessage(msg.UserID, txtCatShowErr)
		}
		return true, model.MessageSender.ShowButtons(msg.UserID, categoriesString, btnStart)
	case "/show_item":
		list, err := getItemList(model, msg.UserID)
		if err != nil {
			return false, err
		}
		return true, model.MessageSender.ShowButtons(msg.UserID, list, btnStart)
	case "/cancel":
		model.lastUserCmd[msg.UserID] = ""
		model.lastUserCat[msg.UserID] = ""
		model.lastUserItemName[msg.UserID] = ""
		return true, model.MessageSender.ShowButtons(msg.UserID, txtChooseCmd, btnStart)
	}
	return false, nil
}

func getCategoryButtons(categoryList []string) []types.TgRowButtons {
	var categoryButtons = []types.TgRowButtons{}
	for i, cat := range categoryList {
		categoryButtons = append(categoryButtons, types.TgRowButtons{})
		categoryButtons[i] = append(categoryButtons[i], types.TgInlineButton{
			DisplayName: cat,
			Value:       "/cat " + cat,
		})
	}
	categoryButtons = append(categoryButtons, types.TgRowButtons{})
	categoryButtons[len(categoryList)] = append(categoryButtons[len(categoryList)], types.TgInlineButton{
		DisplayName: "Без категории",
		Value:       "/cat default",
	})
	return categoryButtons
}

func getCategoryList(model *BotModel, userId int64) (string, error) {
	var result strings.Builder
	result.WriteString(txtCatShow + "\n")
	for i, cat := range model.UserStorage.GetCategories(userId) {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, cat))
	}
	return result.String(), nil
}

func getItemList(model *BotModel, userId int64) (string, error) {
	var result strings.Builder
	result.WriteString(txtItemShow + "\n")
	for cat, items := range model.UserStorage.GetWishListByCategory(userId) {
		result.WriteString(fmt.Sprintf("Категория '%s'\n", cat))
		for i, item := range items {
			result.WriteString(fmt.Sprintf("%d. %s. Сайт: %s\n", i+1, item.Name, item.URL))
		}
	}
	return result.String(), nil
}
