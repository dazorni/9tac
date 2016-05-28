import React from 'react';
import Game from '../../component/game';
import GameForm from '../../component/game-form';

class GameJoinPage extends React.Component {
	constructor(props) {
		super(props);

		let gameCode;

		if (props.params && props.params.gameCode) {
			gameCode = props.params.gameCode;
		}

		this.state = {
			username: null,
			gameCode: gameCode,
			isFormSubmitable: false
		}
	}

	render() {
    return this._joinGame()
  }

	_joinGame() {
		if (this.state.gameCode && this.state.username) {
			return(<Game gameCode={this.state.gameCode} username={this.state.username} />)
		}

		return (
			<div className="container">
				<GameForm handleSubmit={this._handleSubmit.bind(this)} submitTrans="Join" isSubmitable={this.state.isFormSubmitable}>
					{this._getInputFields()}
				</GameForm>
			</div>
		)
	}

	_getInputFields () {
		if (this.state.gameCode) {
			return(
				<fieldset className="form-group"><input type="text" className="form-control" id="username" placeholder="Enter Username" ref={c => this._username = c} onChange={this._handleFormChange.bind(this)} maxlength="20" /></fieldset>
			);
		}

		return (
			<div>
				<fieldset className="form-group"><input type="text" className="form-control" id="username" placeholder="Enter Username" ref={c => this._username = c} onChange={this._handleFormChange.bind(this)} maxlength="20" /></fieldset>
				<fieldset className="form-group"><input type="text" className="form-control" id="gameCode" placeholder="Enter GameCode" ref={c => this._gameCode = c} onChange={this._handleFormChange.bind(this)} /></fieldset>
			</div>
		);
	}

	_handleSubmit(event) {
		event.preventDefault();

		if (this._username.value.length > 20) {
			alert("Choose a shorter username");
		}

		if (this._username.value.length < 3) {
			alert("Choose a longer username");
		}

		this.setState({username: this._username.value});
		this._username.value = '';

		if (this._gameCode && this._gameCode.value) {
			this.setState({gameCode: this._gameCode.value});
			this._gameCode.value = '';
		}
	}

	_handleFormChange() {
		if (! this._username.value) {
			return;
		}

		if (this._gameCode && ! this._gameCode.value) {
			return;
		}

		this.setState({isFormSubmitable: true});
	}
}

export default GameJoinPage
