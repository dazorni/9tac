import React from 'react';
import Game from '../../component/game';
import GameForm from '../../component/game-form';

class GameJoinPage extends React.Component {
	constructor() {
		super();

		this.state = {
			username: null,
			gameCode: null,
			isFormSubmitable: false
		}
	}

	render() {
    return this._joinGame()
  }

	_joinGame() {
		if (this.state.gameCode) {
			return(<Game gameCode={this.state.gameCode} username={this.state.username} />)
		}

		return (
			<div className="container">
				<GameForm handleSubmit={this._handleSubmit.bind(this)} submitTrans="Join" isSubmitable={this.state.isFormSubmitable}>
					<fieldset className="form-group"><input type="text" className="form-control" id="username" placeholder="Enter Username" ref={c => this._username = c} onChange={this._handleFormChange.bind(this)} maxlength="20" /></fieldset>
					<fieldset className="form-group"><input type="text" className="form-control" id="gameCode" placeholder="Enter GameCode" ref={c => this._gameCode = c} onChange={this._handleFormChange.bind(this)} /></fieldset>
				</GameForm>
			</div>
		)
	}

	_handleSubmit(event) {
		event.preventDefault();

		if (this._username.value.length > 20) {
			alert("Choose a shorter username");
		}

		if (this._username.value.length < 3) {
			alert("Choose a longer username");
		}

		this.setState({
			username: this._username.value,
			gameCode: this._gameCode.value
		});

		this._username.value = '';
		this._gameCode.value = '';
	}

	_handleFormChange() {
		if (! this._username.value && ! this._gameCode.value) {
			return;
		}

		this.setState({isFormSubmitable: true});
	}
}

export default GameJoinPage
