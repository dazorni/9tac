import React from 'react';

class UserList extends React.Component {
	constructor() {
		super();

		let userData = [
			{"id": 1, "username": "dazorni"},
			{"id": 2, "username": "nino"},
		];

		this.state = {
			data: userData
		};
	}

	render() {
		let userNodes = this.state.data.map(user => <User key={user.id} username={user.username} avatar={this._getAvatarUrl(user.username)}/>);

		return (
			<div>UserList</div>
			{userNodes}
		)
	}

	_getAvatarUrl(username) {
		return "https://api.adorable.io/avatar/" + username;
	}
};

export default UserList
