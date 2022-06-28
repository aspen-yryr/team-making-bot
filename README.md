# team-making-bot

Discord bot for team making

## Requirement

### Discord bot token

If you don't have token get from [Discord Developer Portal](https://discord.com/developers/docs)

## Genarate .env file

Type below command
```sh
$ make run_debug
```

and PASTE your bot token along prompt. (Token don't displayed for security)
```sh
Please input your bot token to make .env.dev file.
>
```

## Run debug (Visual Studio Code)

Add below configuration to launch.json

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/app/main.go",
            "cwd": "${workspaceFolder}",
            "env": {},
            "args": [
                "-v=6",
                "-logtostderr",
                "--env_file=./env/.env.dev",
                "--greet=false"
            ]
        }
    ]
}
```

## Run debug (console)

```
go run cmd/app/main.go -v 6 -logtostderr --env_file=./env/.env.dev --greet=false
```