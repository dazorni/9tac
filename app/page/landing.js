import React from 'react';
import {Link} from 'react-router';

class LandingPage extends React.Component {
	render() {
    return (
			<div className="landingpage">
				<section className="jumbotron text-xs-center">
		      <div className="container">
		        <h1 className="jumbotron-heading">9Tac</h1>
		        <p className="lead text-muted">A simple game </p>
		        <p>
		          <Link to="/game" className="btn btn-primary">Start</Link>
							<Link to="/game-join" className="btn btn-secondary">Join</Link>
		        </p>
		      </div>
	    	</section>
			</div>
    )
  }
}

export default LandingPage
