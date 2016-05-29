import React from 'react';
import {Link} from 'react-router';

class LandingPage extends React.Component {
	render() {
    return (
			<div className="landingpage">
				<section className="jumbotron text-xs-center">
		      <div className="container">
		        <h1 className="jumbotron-heading">9Tac</h1>
		        <p className="lead text-muted">This is an variation of the classic tic tac toe game. The board is composed of a 3x3 grid of Big Squares. Inside of each Big Square is another 3x3 grid, so there are 9 Small Squares inside each Big Square</p>
		        <p>
		          <Link to="/game" className="btn btn-block btn-primary">Start</Link>
							<Link to="/game-join" className="btn btn-block btn-secondary">Join</Link>
		        </p>
		      </div>
	    	</section>
			</div>
    )
  }
}

export default LandingPage
