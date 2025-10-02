package repository

import (
	"Project/internal/app/ds"
)

func (r *Repository) GetAllChats() ([]ds.Chat, error) {
	// тут мы пользуемся ORM
	var chats []ds.Chat
	err := r.db.Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *Repository) SearchChatsByName(name string) ([]ds.Chat, error) {
	var chats []ds.Chat
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}
