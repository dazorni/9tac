import React from 'react';
import GameField from './game-field';
import IO from 'socket.io-client';
import GamePlayerInfoBox from './game-player-info-box';
import GameModal from './game-modal';
import ErrorModal from './error/modal';
import SharingScreen from './game/sharing-screen';

class Game extends React.Component {
	constructor() {
		super();

		this.state = {
			socket: IO.connect(),
			playedCells: [],
			gameStarted: false,
			opponent: null,
			gameCode: null,
			startingPlayer: null,
			previousPlayer: null,
			nextField: 0,
			isNextFieldRandom: false,
			gameEnded: false,
			winner: null,
			error: null,
			fields: this._createFields()
		}

		this.state.socket.on('game:turn:draw', this._drawTurn.bind(this));
		this.state.socket.on('game:new', this._createGame.bind(this));
		this.state.socket.on('game:join', this._joinGame.bind(this));
		this.state.socket.on('error:error', this._showError.bind(this));
	}

	_createFields() {
		const fields = [];

		for (let field = 0; field < 9; field++) {
			const cells = [];

	    for (let inc = 0; inc <= 8; inc++) {
	        cells.push({position: this._getPosition(field, inc), marked: false, isStartingPlayer: false, isLastTurn: false});
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
				<div id="game-container">
					<div className="game">
						{this._getFields()}
					</div>

					{this._winningModal()}

					{this._errorModal()}

					<nav className="navbar navbar-fixed-bottom navbar-light bg-faded">
						{this._getPlayerInfo()}
					</nav>
				</div>
			)
		}

    if (this.state.gameCode) {
      return (
        <SharingScreen gameCode={this.state.gameCode} />
  		)
    }

    return (
      <div>
				{this._errorModal()}
				Create game...
			</div>
    )
	}

	_getFields() {return this.state.fields.map(field => {
			let isNextField = false;

			if (this.state.isNextFieldRandom || this.state.nextField == field.position) {
				isNextField = true
			}

			return(<GameField
				onTurn={this._onTurn.bind(this)}
				field={field.position}
				cells={field.cells}
				won={field.won}
				isStartingPlayer={ field.winningPlayer == this.state.startingPlayer }
				key={field.position}
				next={isNextField} />)
		})
	}

	_errorModal() {
		if (! this.state.error) {
			return ('');
		}

		return (<ErrorModal message={this.state.error} />)
	}

	_winningModal() {
		if (this.state.gameEnded) {
			let won = false;

			if (this.state.winner == this.props.username) {
				won = true;
			}

			return (
				<GameModal won={won} />
			)
		}

		return ('');
	}

	_getPlayerInfo() {
		if (! this.state.gameStarted) {
			return;
		}

		const playerOne = {
			username: this.props.username,
			isCurrentPlayer: this._isCurrentPlayer(this.props.username),
			isStartingPlayer: this.state.startingPlayer == this.props.username,
			isOpponent: false
		};

		const playerTwo = {
			username: this.state.opponent,
			isCurrentPlayer: this._isCurrentPlayer(this.state.opponent),
			isStartingPlayer: this.state.startingPlayer == this.state.opponent,
			isOpponent: true
		}

		const player = [playerOne, playerTwo];

		return (<GamePlayerInfoBox player={player} />);
	};

	_isCurrentPlayer(username) {
		if (this.state.previousPlayer && this.state.previousPlayer != username) {
				return true;
		}

		if (! this.state.previousPlayer && this.state.startingPlayer == username) {
			return true;
		}

		return false;
	};

	componentWillMount() {
		this._buildGame();
	};

	_createGame(gameCode) {
		this.setState({gameCode: gameCode});
	}

	_joinGame(firstPlayer, secondPlayer, gameCode, startingField) {
		let opponent = firstPlayer;

		if (this.props.username != secondPlayer) {
			opponent = secondPlayer;
		}

		this.setState({
			opponent: opponent,
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

		if (this.state.fields && this.state.fields[field] && this.state.fields[field].won) {
			return;
		}

		if (this.state.isNextFieldRandom == false && this.state.nextField != field) {
			return;
		}

		if (this.state.previousPlayer == this.props.username) {
			return;
		}

		if (! this.state.previousPlayer && this.state.startingPlayer != this.props.username) {
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
			field.cells.map(cell => {
				cell.isLastTurn = false;

				return cell;
			});

			if (field.position == turn.Field) {
				field.cells.map(cell => {
					cell.isLastTurn = false;

					if (cell.position == turn.Position) {
						cell.marked = true;
						cell.isStartingPlayer = this.state.startingPlayer == username;
						cell.isLastTurn = true;
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

		let isNextFieldRandom = false;

		if (turn.RandomField) {
			isNextFieldRandom = true;
		}

		this.setState({
			previousPlayer: username,
			nextField: turn.NextField,
			fields: fields,
			playedCells: playedCells,
			isNextFieldRandom: isNextFieldRandom
		});
	};

	_buildGame() {
		if (this.props.gameCode) {
			this.state.socket.emit('game:join', this.props.username, this.props.gameCode);
		} else {
			this.state.socket.emit('game:new', this.props.username);
		}
	}

	_showError(errorMessage) {
		this.setState({error: errorMessage});
	}
}

Game.propTypes = {
	username: React.PropTypes.string.isRequired,
	gameCode: React.PropTypes.string
}

export default Game
