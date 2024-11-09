# Go game 


The purpose of the repo is to explore and demo potential of [flowstate](https://github.com/makasim/flowstate) library.

For now a player can create a game or find vacant game to join. 
Once the game has two players, the game starts. 
The server randomly decided who starts the game.
The players do moves in turns withing given time frame (30sec).
Either player can resign the game at any time.
There is no actual game logic implemented yet.

More to come.

## How to run

Start server:
```bash
go run main/main.go 
```

Start demo client: 
```bash
go run play.go
```

To build the UI read [its readme](./ui/README.md)

## Credits

- [Otrego](https://github.com/otrego) community for building [clamshell](https://github.com/otrego/clamshell) that implement go board rules.
