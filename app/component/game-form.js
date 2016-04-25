import React from 'react';

class GameForm extends React.Component {
  render() {
    return (
      <div className="game-form-wrapper">
        <form className="game-form" onSubmit={this.props.handleSubmit.bind(this)} >
          {this.props.children}
          <button type="submit" className="btn btn-block btn-primary">{this.props.submitTrans}</button>
        </form>
      </div>
    )
  }
}

GameForm.propTypes = {
  handleSubmit: React.PropTypes.func.isRequired,
  submitTrans: React.PropTypes.string.isRequired,
  isSubmitable: React.PropTypes.bool.isRequired
}

export default GameForm;
