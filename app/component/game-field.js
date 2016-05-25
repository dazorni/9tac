import React from 'react';
import GameCell from './game-cell';

class GameField extends React.Component {
  constructor() {
    super();

    this.state = {
      cells: []
    }
  }

  render() {
    return(
      <div className={this._getClassName()} data-field={this.props.field}>
        {this._getGameCells()}
      </div>
    )
  }

  componentWillMount() {
    this.setState({cells: this.props.cells});
  }

  _getClassName() {
    let className = 'game-field';

    if (this.props.won) {
      let playerClass = 'game-field-won-player-one';

      if (! this.props.isStartingPlayer) {
        playerClass = 'game-field-won-player-two';
      }

      className = className + ' game-field-won ' + playerClass;
    }

    if (this.props.next) {
      className = className + ' game-field-next';
    }

    return className
  }

  _getGameCells() {
    return this.state.cells.map((cell) => {
      return <GameCell
        key={cell.position}
        position={cell.position}
        field={this.props.field}
        onTurn={this.props.onTurn}
        marked={cell.marked}
        isLastTurn={cell.isLastTurn}
        isStartingPlayer={cell.isStartingPlayer} />
    });
  }
}

GameField.propTypes = {
  field: React.PropTypes.any.isRequired,
  onTurn: React.PropTypes.func.isRequired,
  isStartingPlayer: React.PropTypes.bool.isRequired
};

export default GameField;
