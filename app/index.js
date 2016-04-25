import ReactDom from 'react-dom';
import React from 'react';
import { Router, Route, hashHistory, IndexRoute } from 'react-router';
import App from './app';
import LandingPage from './page/landing';
import GameStartPage from './page/game/start';
import GameJoinPage from './page/game/join';
import GamePlayPage from './page/game/play';

const app = (
  <Router history={hashHistory}>
    <Route path="/" component={App}>
			<IndexRoute component={LandingPage} />
			<Route path="game" component={GameStartPage} />
      <Route path="game-join" component={GameJoinPage} />
      <Route path="game/:gameCode" component={GamePlayPage} />
    </Route>
  </Router>
)

ReactDom.render(app, document.getElementById('app'));
