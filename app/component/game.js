import React from 'react';
import GameField from './game-field';
import IO from 'socket.io-client';

class Game extends React.Component {
	constructor() {
		super();

		this.state = {
			socket: IO.connect(),
			playedCells: [],
			gameStarted: false,
			opponnent: null,
			gameCode: null,
			startingPlayer: null,
			previousPlayer: null,
			nextField: 0,
			gameEnded: false,
			winner: null,
			fields: this._createFields()
		}

		this.state.socket.on('game:turn:draw', this._drawTurn.bind(this));
		this.state.socket.on('game:new', this._createGame.bind(this));
		this.state.socket.on('game:join', this._joinGame.bind(this));
	}

	_createFields() {
		const fields = [];

		for (let field = 0; field < 9; field++) {
			const cells = [];

	    for (let inc = 0; inc <= 8; inc++) {
	        cells.push({position: this._getPosition(field, inc), marked: false, isStartingPlayer: false});
	    }

			fields.push({ position: field, won: false, player: null, cells: cells, next: false})
		}

		return fields
	}

	_getPosition(field, position) {
    const resultX = Math.floor((position % 3) + (field * 3) % 9)
    const fieldY = Math.floor(field / 3);
    const resultY = Math.floor(position / 3) + (fieldY * 3);

    return resultX + (resultY * 9)
  }

	render() {
		return (this._play())
	};

	_play() {
		if (this.state.gameStarted) {
			return (
				<div className="game">
					<div>GameCode: {this.state.gameCode}</div>
					<div>{this.props.username} vs. {this.state.opponnent}</div>

					{this._getFields()}

					{this._whichTurn()}
					{this._winningCheck()}
				</div>
			)
		}

		return (
			<div>GameCode: {this.state.gameCode}</div>
		)
	}

	_getFields() {return this.state.fields.map(field => {
			return(<GameField
				onTurn={this._onTurn.bind(this)}
				field={field.position}
				cells={field.cells}
				won={field.won}
				isStartingPlayer={ field.winningPlayer == this.state.startingPlayer }
				key={field.position}
				next={this.state.nextField == field.position} />)
		})
	}

	_winningCheck() {
		console.log(this.state.gameEnded);

		if (this.state.gameEnded) {
			if (this.state.winner == this.props.username) {
				return ('You won');
			}

			return ('You lost the game');
		}

		return ('');
	}

	_whichTurn() {
		if (this.state.previousPlayer && this.state.previousPlayer != this.props.username) {
				return "Your turn";
		}

		if (! this.state.previousPlayer && this.state.startingPlayer == this.props.username) {
			return "Your turn";
		}

		return `${this.state.opponnent} is next`;
	}

	componentWillMount() {
		this._buildGame();
	};

	_createGame(gameCode) {
		this.setState({gameCode: gameCode});
	}

	_joinGame(firstPlayer, secondPlayer, gameCode, startingField) {
		let opponnent = firstPlayer;

		if (this.props.username != secondPlayer) {
			opponnent = secondPlayer;
		}

		this.setState({
			opponnent: opponnent,
			gameCode: gameCode,
			gameStarted: true,
			startingPlayer: firstPlayer,
			nextField: startingField
		});
	}

	_onTurn(position, field) {
		if (this.state.playedCells.indexOf(position) != -1) {
			return;
		}

		if (this.state.nextField != field) {
			return;
		}

		if (this.state.previousPlayer == this.props.username) {
			return;
		}

		if (! this.state.	previousPlayer && this.state.startingPlayer != this.props.username) {
			return;
		}

		if (this.state.gameEnded) {
			return;
		}

		this.state.socket.emit('game:turn', this.props.username, this.state.gameCode, position);
	}

	_drawTurn(turn, username) {
		const fields = this.state.fields;

		fields.map(field => {
			if (field.position == turn.Field) {
				field.cells.map(cell => {
					if (cell.position == turn.Position) {
						cell.marked = true
						cell.isStartingPlayer = this.state.startingPlayer == username
					}

					return cell
				});

				if (turn.WonField) {
					field.won = turn.WonField;
					field.winningPlayer = username;
				}
			}

			return field
		});

		const playedCells = this.state.playedCells;
		playedCells.push(turn.Position);

		if (turn.WonGame) {
			this.setState({wonGame: true, winner: username, gameEnded: true});
		}

		this.setState({
			previousPlayer: username,
			nextField: turn.NextField,
			fields: fields,
			playedCells: playedCells
		});
	};

	_buildGame() {
		if (this.props.gameCode) {
			this.state.socket.emit('game:join', this.props.username, this.props.gameCode);
		} else {
			this.state.socket.emit('game:new', this.props.username);
		}
	}
}

Game.propTypes = {
	username: React.PropTypes.string.isRequired,
	gameCode: React.PropTypes.string
}

export default Game
