package bugs

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/repo"
)

type ReportCommand struct {
	s *state.State
	c *env.Config

	menu *BugsMenu

	interfaces.ContextGenerator
}

func newReportCommand(m *BugsMenu) interfaces.Command {
	return &ReportCommand{
		c:                m.cr.Config(),
		s:                m.cr.State(),
		menu:             m,
		ContextGenerator: m.cr.ContextGenerator(),
	}
}

func (cmd *ReportCommand) GetName() string {
	return "report"
}

func (cmd *ReportCommand) GetDescription() string {
	return "Report a bug or an issue with the bot"
}

func (cmd *ReportCommand) RegisterCommand() (*api.CreateCommandData, error) {
	cmd.s.AddHandler(cmd.HandleCommand())

	cmd.s.AddHandler(cmd.HandleResponse())

	return &api.CreateCommandData{
		Name:        cmd.GetName(),
		Description: cmd.GetDescription(),
		Type:        discord.ChatInputCommand,

		// subcommands for bugs and issues
		Options: []discord.CommandOption{
			&discord.SubcommandOption{
				OptionName:  "bug",
				Description: "Report a bug",
			},
			&discord.SubcommandOption{
				OptionName:  "issue",
				Description: "Report an issue",
			},
		},
	}, nil
}

func (cmd *ReportCommand) HandleCommand() func(*gateway.InteractionCreateEvent) {
	return func(e *gateway.InteractionCreateEvent) {
		switch e.Data.(type) {
		default:
			return
		case *discord.CommandInteraction:
		}

		d := e.Data.(*discord.CommandInteraction)

		if d.Name != cmd.GetName() {
			return
		}

		// Get type
		var t string
		for _, o := range d.Options {
			if o.Name == "bug" || o.Name == "issue" {
				t = o.Name
				break
			}
		}

		// Reply with modal
		err := cmd.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
			Type: api.ModalResponse,
			Data: &api.InteractionResponseData{
				Title: option.NewNullableString("Report a new " + t),
				Components: discord.ComponentsPtr(
					&discord.TextInputComponent{
						CustomID:    "title",
						Placeholder: "A brief title of the issue",
						Label:       "Title",
						Style:       discord.TextInputShortStyle,
						Required:    true,
					},
					&discord.TextInputComponent{
						CustomID:    "description",
						Placeholder: "A detailed description of the issue",
						Label:       "Description",
						Style:       discord.TextInputParagraphStyle,
						Required:    true,
					},
				),
				CustomID: option.NewNullableString("modal_report_" + t),
			},
		})
		if err != nil {
			logrus.WithError(err).WithField("command", "bug").Error("Failed to respond to interaction")
			return
		}
	}
}

func (cmd *ReportCommand) HandleResponse() func(*gateway.InteractionCreateEvent) {
	return func(e *gateway.InteractionCreateEvent) {
		switch e.Data.(type) {
		default:
			return
		case *discord.ModalInteraction:
		}

		d := e.Data.(*discord.ModalInteraction)
		if d.CustomID != "modal_report_bug" && d.CustomID != "modal_report_issue" {
			return
		}

		reportType := string(d.CustomID[13:])
		reportTitle := d.Components.Find("title").(*discord.TextInputComponent).Value
		reportDescription := d.Components.Find("description").(*discord.TextInputComponent).Value

		var response api.InteractionResponse

		issue, err := cmd.menu.cr.NewGithubIssue(repo.NewIssueParams{
			Title:       fmt.Sprintf("New %s: %s", reportType, reportTitle),
			Description: reportDescription,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create issue")
			response = api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString("Failed to create issue"),
					Flags:   discord.EphemeralMessage,
				},
			}
		} else {
			response = api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Embeds: &[]discord.Embed{
						{
							Title: "Reported " + reportType,
							Description: "Thank you for the report.\n" +
								"Your report can be seen [here](" + issue.GetHTMLURL() + ")\n" +
								"You can complement your issue by logging in and completing it.",
						},
					},
				},
			}
		}

		// Reply with modal
		err = cmd.s.RespondInteraction(e.ID, e.Token, response)
		if err != nil {
			logrus.WithError(err).WithField("command", "bug").Error("Failed to respond to interaction")
			return
		}
	}
}
