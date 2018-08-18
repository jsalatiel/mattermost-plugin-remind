package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

func (p *Plugin) registerCommand(teamId string) error {

	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone or ~channel] [what] [when]",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		DisplayName:      "Remind Plugin Command",
		Description:      "Set a reminder",
	}); err != nil {
		p.API.LogError(
			"failed to register command",
			"error", err.Error(),
		)
	}
	
	p.run()

	return nil
}

func (p *Plugin) unregisterCommand(teamId string) error {

	p.API.UnregisterCommand(teamId, CommandTrigger);

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	p.API.LogDebug("ExecuteCommand")

	user, err := p.API.GetUser(args.UserId)
	
	if err != nil {
		p.API.LogError("failed to query user %s", args.UserId)
	}

	if strings.HasSuffix(args.Command, "help") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf(HelpText),
		}, nil
	}

	if strings.HasSuffix(args.Command, "list") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("todo"),
		}, nil
	}

	if strings.HasSuffix(args.Command, "version") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf(Version),
		}, nil
	}

	if strings.HasSuffix(args.Command, "debug") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("* %s\n * %s\n * %s\n * %s\n * %s\n * %s\n", 
				args.Command, 
				args.TeamId,
				args.SiteURL,
				user.Username, 
				user.Id,  
				user.Timezone["automaticTimezone"]),
		}, nil
	}

	payload := strings.Trim(strings.Replace(args.Command, "/"+CommandTrigger, "", -1),"")
	commandSplit := strings.Split(payload," ")

	if len(commandSplit) == 0 {	

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf(ExceptionText),
		}, nil
	}

	if commandSplit[1] == "me" ||
		strings.HasPrefix(commandSplit[1][:1], "@") ||
		strings.HasPrefix(commandSplit[1][:1], "~") {

		request := ReminderRequest{args.TeamId, user.Username, payload, Reminder{}}
		response, err := p.scheduleReminder(request)

		if err != nil {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text: fmt.Sprintf(ExceptionText),
			}, nil
		}

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("%s",response),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: fmt.Sprintf(ExceptionText),
	}, nil

}
