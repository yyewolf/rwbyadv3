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

func ReportCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
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

func (cmd *reportCommand) HandleCommand(t string) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		return e.Modal(discord.NewModalCreateBuilder().
			SetCustomID("modal_report_" + t).
			SetTitle("Report a new " + t).
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

func (cmd *reportCommand) HandleResponse(t string) handler.ModalHandler {
	return func(e *handler.ModalEvent) error {
		d := e.Data

		reportType := string(d.CustomID[13:])
		reportTitle := d.Text("title")
		reportDescription := d.Text("description")

		issue, err := cmd.app.NewGithubIssue(repo.NewIssueParams{
			Title:       fmt.Sprintf("New %s: %s", reportType, reportTitle),
			Description: reportDescription,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create issue")
			return e.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetContent("Failed to create issue").
					SetEphemeral(true).
					Build(),
			)
		}
		return e.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Reported"+reportType).
						SetDescriptionf(
							"Thank you for the report.\n"+
								"Your report can be seen [here](%s)\n"+
								"You can complement your issue by logging in and completing it.",
							issue.GetHTMLURL(),
						).
						Build(),
				).
				Build(),
		)
	}
}
