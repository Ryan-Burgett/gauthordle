# gauthordle
Given a commit message, can you guess who wrote it?

https://github.com/JosephNaberhaus/gauthordle/assets/26124819/d43e249f-e7eb-4f32-9558-54f3303aeb35

New games come out at midnight Central Time. Each game is deterministically generated from your git history so everyone can play the same puzzle each day.

## How to play
To play, follow one of the installation guides below. After that, you can just navigate to any git repository, `git chekout` the main development branch, and then run `gauthordle` to start the game.

You might also want to `git pull` to ensure that your git history is up-to-date. Otherwise, you may end up playing the wrong game for the day. The program doesn't do this automatically because I didn't want it to make any changes to the file system.

### Installation from source (recommended)
With any v1.21 or greater of [Golang](https://go.dev/) you can easily install from source:

```shell
go install github.com/josephnaberhaus/gauthordle@latest
```

### Installation from prebuilt binary
Download and extract the appropriate binary for your platform on the [releases page](https://github.com/JosephNaberhaus/gauthordle/releases).
