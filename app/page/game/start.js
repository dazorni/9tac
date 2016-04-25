import React from 'react';
import Game from '../../component/game';
import GameForm from '../../component/game-form';

class GameStartPage extends React.Component {
	constructor() {
		super();

		this.state = {
			isFormSubmitable: false
		};
	}

	render() {
    return (
			<div>
				{this._startGame()}
			</div>
    )
  }

	_startGame() {
		if (this.state.username) {
				return (<Game username={this.state.username} />);
		}

		return(
			<GameForm handleSubmit={this._handleSubmit.bind(this)} submitTrans="Start game" isSubmitable={this.state.isFormSubmitable}>
				<fieldset className="form-group">
    			<input type="text" className="form-control" id="username" placeholder="Enter Username" ref={c => this._username = c} onChange={this._handleFormChange.bind(this)} />
  			</fieldset>
			</GameForm>
		);
	}

	_handleSubmit(event) {
		event.preventDefault();

		this.setState({username: this._username.value});
	}

	_handleFormChange() {
		if (! this._username.value) {
			return;
		}

		this.setState({
			isFormSubmitable: true
		});
	}
}

export default GameStartPage
