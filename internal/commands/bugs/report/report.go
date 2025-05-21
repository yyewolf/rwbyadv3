package report

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/repo"
)

const (
	commandName        = "report"
	commandDescription = "Report a bug or an issue with the bot"

	modalBugId   = "modal_report_bug"
	modalIssueId = "modal_report_issue"
)

type reportCommand struct {
	app interfaces.App
}

func ReportCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd reportCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/report/bug", cmd.HandleCommand("bug"))
			h.Command("/report/issue", cmd.HandleCommand("issue"))
			h.Modal("/"+modalBugId, cmd.HandleResponse("bug"))
			h.Modal("/"+modalIssueId, cmd.HandleResponse("issue"))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			// subcommands for bugs and issues
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "bug",
					Description: "Report a bug",
				},
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "issue",
					Description: "Report an issue",
				},
			},
		}),
	)
}

func (cmd *reportCommand) HandleCommand(reportType string) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		return event.Modal(discord.NewModalCreateBuilder().
			SetCustomID("modal_report_" + reportType).
			SetTitle("Report a new " + reportType).
			AddContainerComponents([]discord.ContainerComponent{
				discord.NewActionRow().AddComponents(
					discord.NewShortTextInput("title", "Title").
						WithRequired(true),
				),
				discord.NewActionRow().AddComponents(
					discord.NewParagraphTextInput("description", "Description").
						WithRequired(true),
				),
			}...).
			Build(),
		)
	}
}

func (cmd *reportCommand) HandleResponse(reportType string) handler.ModalHandler {
	return func(event *handler.ModalEvent) error {
		form := event.Data

		reportType := string(form.CustomID[13:])
		reportTitle := form.Text("title")
		reportDescription := form.Text("description")

		issue, err := cmd.app.Github().NewGithubIssue(repo.NewIssueParams{
			Title:       fmt.Sprintf("New %s: %s - %s", reportType, reportTitle, event.User().ID),
			Description: reportDescription,
		})
		if err != nil {
			logrus.WithField("user_id", event.User().ID).WithError(err).Error("Failed to create issue")
			return event.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetContent("Failed to create the bug report.").
					SetEphemeral(true).
					Build(),
			)
		}

		return event.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Reported "+reportType).
						SetDescriptionf(
							"Thank you for the report.\n"+
								"Your report can be seen [here](%s)\n"+
								"You can complement your issue by logging in and completing it.",
							issue.GetHTMLURL(),
						).
						SetEmbedFooter(cmd.app.Footer()).
						SetColor(cmd.app.Config().App.BotColor).
						Build(),
				).
				Build(),
		)
	}
}
