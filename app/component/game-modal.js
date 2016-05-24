import React from 'react';

class GameModal extends React.Component {
  constructor() {
    super();

    this.state = {
      hide: false
    }
  }

  render() {
    if (this.state.hide) {
      return (null);
    }

    return(
      <div className="game-modal-container">
        <div className="game-modal">
          <button type="button" className="close" aria-label="Close" onClick={this._hide.bind(this)}>
            <span aria-hidden="true">&times;</span>
          </button>
          <div className={this._getIconClass()}></div>
          <div className="game-modal-text">{this._getText()}</div>
        </div>
      </div>
    );
  }

  _getIconClass() {
    let iconClass = 'game-modal-icon-lost';

    if (this.props.won) {
      iconClass = 'game-modal-icon-won';
    }

    return 'game-modal-icon ' + iconClass;
  }

  _getText() {
    if (this.props.won) {
      return ('You won that game!');
    }

    return ('You lost that game...');
  }

  _hide() {
    this.setState({hide: true});
  }
}

GameModal.propTypes = {
  won: React.PropTypes.bool.isRequired,
}

export default GameModal;
