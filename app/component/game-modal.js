import React from 'react';
import Modal from './modal';

class GameModal extends React.Component {
  render() {
    return(<Modal iconClass={this._getIconClass()} text={this._getText()} />);
  }

  _getIconClass() {
    if (this.props.won) {
      return 'game-modal-icon-won';
    }

    return 'game-modal-icon-lost';
  }

  _getText() {
    if (this.props.won) {
      return ('You won that game!');
    }

    return ('You lost that game...');
  }
}

GameModal.propTypes = {
  won: React.PropTypes.bool.isRequired,
}

export default GameModal;
