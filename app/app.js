import React from 'react';

class App extends React.Component {
	render() {
		return (
			<div>
				<nav className="navbar navbar-light bg-faded">
					<a className="navbar-brand navbar-brand-center" href="#">9tac</a>
				</nav>

				<div className="container">{this.props.children}</div>
			</div>
		)
	}
}

export default App
