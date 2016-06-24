import React from 'react';

class Modal extends React.Component {
  constructor() {
    super();

    this.state = {
      hidden: false
    }
  }

  render() {
    if (this.state.hidden) {
      return (null);
    }

    return(
      <div className="modal-container">
        <div className="modal-box">
          {this._getHidableButton()}
          <div className={"modal-icon " + this.props.iconClass}></div>
          <div className="modal-text">{this.props.text}</div>
        </div>
      </div>
    );
  }

  _getHidableButton() {
    if (! this.props.hidable) {
      return ('');
    }

    return (
      <button type="button" className="close" aria-label="Close" onClick={this._hide.bind(this)}>
        <span aria-hidden="true">&times;</span>
      </button>
    );
  }

  _hide() {
    this.setState({hidden: true});
  }
}

Modal.PropTypes = {
  iconClass: React.PropTypes.string.isRequired,
  text: React.PropTypes.string.isRequired,
  hideable: React.PropTypes.bool
}

Modal.DefaultProps = {
  hidable: true
}

export default Modal;
