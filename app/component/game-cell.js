import React from 'react';

class GameCell extends React.Component {
  render() {
    return (
      <div onClick={this._onClick.bind(this)} className={this._getClassName()} data-position={this.props.position}>
        <div className="content"></div>
      </div>
    )
  }

  _getClassName() {
    let className = 'game-cell';

    if (this.props.marked) {
      let playerClass = '-pone';

      if (! this.props.isStartingPlayer) {
        playerClass = '-ptwo';
      }

      className = className + ' game-cell-marked game-cell-marked' + playerClass;
    }

    return className
  }

  _onClick() {
    this.props.onTurn(this.props.position, this.props.field);
  }
}

GameCell.propTypes = {
  position: React.PropTypes.any.isRequired,
  onTurn: React.PropTypes.func.isRequired
}

export default GameCell;
