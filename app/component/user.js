import React from 'react';

class User extends React.Component {
  render() {
    return (
      <div className="user">
        <div className="user-username"></div>
        <div classNames="user-avatar">
          <img src="{this.props.avatarUrl}" />
        </div>
      </div>
    )
  }
}

User.propTypes = {
  username: React.PropTypes.string.isRequired,
  avatarUrl: React.PropTypes.string.isRequired
}

export default User
