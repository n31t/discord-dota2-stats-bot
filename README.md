
# Project Title

This is a Discord bot written in Go. It uses the discordgo and godotenv libraries.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.16 or later)
- A Discord account and a server to test the bot
- Discord BOT token 
- STRATZ API token

### Installation

Clone the repository to your local machine.

### Commands
This Discord bot responds to commands that start with !neit. Here are the available commands:

!neit ping
The bot responds with "pong". This command is useful for checking if the bot is online and responsive.

!neit embed
The bot sends an embedded message with a title and description. This command is useful for testing the bot's ability to send embedded messages.

!neit stratz <id>
The bot fetches and displays data about a Dota 2 player from the STRATZ API. Replace <id> with the ID of the player you want to look up. The bot displays the following information:

- ID
- Last Active Time
- Dota Plus Subscriber status
- Smurf Flag
- Match Count
- Win Count
- Behavior Score
- Season Rank
- The bot also displays the player's rank history.

Please note that you need a STRATZ API token to use this command. You can get your STRATZ API token on https://stratz.com/api.

