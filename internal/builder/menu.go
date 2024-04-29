package builder

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type MenuStore struct {
	Menus    []*Menu
	Commands []*Command

	h *handler.Mux
}

func NewMenuStore(app interfaces.App) *MenuStore {
	return &MenuStore{
		h: app.Handler(),
	}
}
func (ms *MenuStore) RegisterCommands() ([]discord.ApplicationCommandCreate, error) {
	var data []discord.ApplicationCommandCreate

	for _, m := range ms.Menus {
		val, err := m.RegisterCommands()
		if err != nil {
			return data, err
		}
		data = append(data, val...)
	}

	return data, nil
}

type Menu struct {
	Name  string
	Emoji discord.Emoji

	Commands []*Command

	h *handler.Mux
}

func (ms *MenuStore) NewMenu(opts ...MenuOption) *Menu {
	var menu = &Menu{
		h: ms.h,
	}

	for _, opt := range opts {
		opt(menu)
	}

	ms.Menus = append(ms.Menus, menu)
	ms.Commands = append(ms.Commands, menu.Commands...)

	return menu
}

func (m *Menu) RegisterCommands() ([]discord.ApplicationCommandCreate, error) {
	var data []discord.ApplicationCommandCreate

	for _, c := range m.Commands {
		err := c.Register()
		if err != nil {
			return nil, err
		}
		val := c.createCommand
		data = append(data, val)
	}

	return data, nil
}
