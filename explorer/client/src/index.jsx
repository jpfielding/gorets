import React from 'react';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';
import ReactDOM from 'react-dom';
import App from 'components/App';
import Home from 'components/containers/Home';
import Login from 'components/containers/Login';
import NotFound from 'components/containers/NotFound';

import 'styles/app.css';

ReactDOM.render((
  <Router history={browserHistory}>
    <Route path="/" component={App}>
      <IndexRoute component={Home} />
      <Route path="login" component={Login} />
      <Route path="*" component={NotFound} />
    </Route>
  </Router>
  ), document.getElementById('app')
);
