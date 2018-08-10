import React from 'react';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';
import ReactDOM from 'react-dom';
import App from 'components/App';

import 'styles/app.css';

ReactDOM.render((
  <Router history={hashHistory}>
    <Route path="/" component={App}>
      <IndexRoute component={App} />
      <Route path="*" component={App} />
    </Route>
  </Router>
  ), document.getElementById('app')
);
