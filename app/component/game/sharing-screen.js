import React from 'react';
import copy from 'copy-to-clipboard';

class SharingScreen extends React.Component {
  constructor() {
    super();

    this.state = {
      copied: false
    };
  }

  render() {
    return (
      <div className="container">
        <div className="sharing-screen">
          <i className="icon-loading"></i>
          <small className="loading-text">Waiting for opponent...</small>
          <kbd className="sharing-screen-link">{this._getUrl()}</kbd>
          <p className="sharing-screen-info">Share this link to your friend and start the game!</p>
          {this._getShareButton()}
        </div>
      </div>
    )
  }

  _getUrl() {
    const origin = window.location.origin;

    return origin + '/#/game-join/' + this.props.gameCode;
  }

  _getShareButton() {
    let text = 'Copy & Share';
    let className = 'btn-share';

    if (this.state.copied) {
      text = 'Copied to Clipboard √';
      className = className + ' btn-share-active';
    }

    setTimeout(() => this.setState({copied: false}), 5000);

    return (<button className={className} onClick={this._copyShareUrl.bind(this)}>{text}</button>);
  }

  _copyShareUrl() {
    copy(this._getUrl());

    this.setState({copied: true});
  }
}

SharingScreen.propTypes = {
  gameCode: React.PropTypes.string.isRequired
}

export default SharingScreen;
