# slack-circuit-service

A Slack bot that will provide automation services for circuit members.

- Invite people who submit form entries to the Slack request form after a quick look-over by an ASA board member
- Allow Slack guests to join any "open" #circuit- channels that may exist
- Organize and allow ASA board members to quickly see who is part of which team. Eventually tying this to Ekta membership would be godlike
- Potentially automatically create a "team channel" if any of the captains would like that
- Automatically move captains who have graduated into the #circuit-alumni channel and out of the #circuit-officers channel

## Tokens:

https://api.slack.com/apps/A01GA4V0C05/oauth?

- `slackUserToken`: user token
- `slackBotToken`: bot token

## Slackbot commands

### Available to guests

| Command                 | Behavior                               |
| ----------------------- | -------------------------------------- |
| `list`                  | list all circuit-accessible channels   |
| `join <channel>`        | join a circuit-accessible channel      |
| `directors <A3 \| ASA>` | list Slack info of ASA or A3 directors |
