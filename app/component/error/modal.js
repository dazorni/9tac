import React from 'react';
import Modal from '../modal';

class ErrorModal extends React.Component {
  render() {
    return (
      <Modal
        text={this.props.message}
        iconClass={this._getIconClass()}
        hidable={this._isHideable()} />
    )
  }

  _getIconClass() {
    return "modal-icon-" + this.props.type;
  }

  _isHideable() {
    if (! this.props.type == "error") {
      return true;
    }

    return false;
  }
}

ErrorModal.propTypes = {
  message: React.PropTypes.string.isRequired,
  type: React.PropTypes.oneOf(['error', 'warning', 'info'])
}

ErrorModal.defaultProps = {
  type: "error"
}

export default ErrorModal;
