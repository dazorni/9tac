import React from 'react';

class App extends React.Component {
	render() {
		return (
			<div>
				<nav className="navbar navbar-fixed-top navbar-light bg-faded">
					<a className="navbar-brand navbar-brand-center" href="#">9tac</a>
				</nav>

				{this.props.children}
			</div>
		)
	}
}

export default App
