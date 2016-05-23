import React from 'react';

class GamePlayerInfo extends React.Component {
	render() {
		return(
			<div className={this._getInfoClass()}>
				<div className={this._getIconClass(this.props.username)}></div>
				<div className="username">{this.props.username}</div>
			</div>
		);
	};

	_getInfoClass() {
		let className = "player-info";

		if (this.props.isOpponent) {
			className = className + " player-info-opponent";
		}

		if (this.props.isCurrentPlayer) {
			className = className + " player-info-active";
		}

		return className;
	}

	_getIconClass() {
		let className = "icon-player-one";

		if (! this.props.isStartingPlayer) {
			className = "icon-player-two";
		}

		return "icon " + className;
	};
}

GamePlayerInfo.propTypes = {
	username: React.PropTypes.string.isRequired,
	isStartingPlayer: React.PropTypes.bool.isRequired,
	isCurrentPlayer: React.PropTypes.bool.isRequired,
	isOpponent: React.PropTypes.bool.isRequired
}

export default GamePlayerInfo;
