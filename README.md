# tmuxify
simple session manager for tmux.

## Installtion
```bash
go install github.com/amirrezaask/tmuxify@latest
```

## How to use
- create a yaml file like below:
    ```yaml
            name: your session name
            cwd: your cwd
            windows: 
                - window 1
                - window 2
    ```

- run `tmuxify [your file name, defaults to .tmuxify.yml]
