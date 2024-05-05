package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/boxes"
	"github.com/yyewolf/rwbyadv3/internal/commands/bugs"
	"github.com/yyewolf/rwbyadv3/internal/commands/general"
	"github.com/yyewolf/rwbyadv3/internal/commands/rewards"
	"github.com/yyewolf/rwbyadv3/internal/commands/system"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func RegisterCommands(app interfaces.App) *builder.MenuStore {
	ms := builder.NewMenuStore(app)

	general.NewMenu(ms, app)
	boxes.NewMenu(ms, app)
	rewards.NewMenu(ms, app)
	system.NewMenu(ms, app)
	bugs.NewMenu(ms, app)

	createCommands, err := ms.RegisterCommands()
	if err != nil {
		logrus.Fatal("Couldn't load commands")
	}

	app.Client().Rest().SetGlobalCommands(
		app.Client().ApplicationID(),
		createCommands,
	)

	return ms
}
