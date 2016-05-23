import React from 'react';
import GamePlayerInfo from './game-player-info';

class GamePlayerInfoBox extends React.Component {
  render() {
    return(
      <div className="player-info-box">
        {this._getPlayerInfo()}
      </div>
    );
  };

  _getPlayerInfo() {
    const player = this.props.player;

		return player.map(player => {
      return (
        <GamePlayerInfo
          username={player.username}
          isCurrentPlayer={player.isCurrentPlayer}
          isStartingPlayer={player.isStartingPlayer}
          isOpponent={player.isOpponent}
          key={player.username}
        />
        );
		});
  }
}

GamePlayerInfoBox.propTypes = {
    player: React.PropTypes.array.isRequired
};

export default GamePlayerInfoBox;
