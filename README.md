# Hexapawn

## Description

The game Hexapawn consists of chess pawns on an m-by-n chess board. The pawns begin on a single row and move forward one square at a time capturing diagonally. The game is won by the side that gets a pawn to the other side of the board or when all opponent pawns have been captured. Stalemate occurs when a side cannot make a legal move. As in chess, white moves first. Unlike chess, pawns move only one square at a time, even on a pawn's first move and capturing en passant is not allowed.

### Game Play Example

+-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+
|b|b|b|     |b|b|b|     |b|b| |     | |b| |     | |b| |
+-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+
| | | | --> | |w| | --> | |b| | --> | |b| | --> | |w| | --> *stalemate*
+-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+
|w|w|w|     |w| |w|     |w| |w|     |w| | |     | | | |
+-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+     +-+-+-+

## Game Modes

The following game modes are available:

* **cvc**: two npcs play against each other.
* **cvp**: play as black against a trained white npc.
* **pvc**: play as white against a trained black npc.
* **pvp**: two people play each other.

## Training an NPC

The agent consists of a set of positions it has seen before with a list of available actions. An action is selected at random, but the probability of selecting an action is determined by a weight that is adjusted by a learning rate during training. When a game is won, actions that contributed to winning are incremented and all other actions are decremented. When a game is lost, actions that contributed to losing are decremented and all other actions are incremented. The learning rate `r` on the range `(0,1)` for a selected action is a constant, but the learning rate `p` for all other `n-1` actions in a position defined as `p := r/(n-1)` when `n>1` and `p := 1` for `n < 2`.
