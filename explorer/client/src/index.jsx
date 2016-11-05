import React from 'react';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';
import ReactDOM from 'react-dom';
import App from 'components/App';
import Home from 'components/containers/Home';
import Connections from 'components/containers/Connections';
import Metadata from 'components/containers/Metadata';
import Search from 'components/containers/Search';
import Objects from 'components/containers/Objects';
import NotFound from 'components/containers/NotFound';

import 'styles/app.css';

ReactDOM.render((
  <Router history={hashHistory}>
    <Route path="/" component={App}>
      <IndexRoute component={Home} />
      <Route path="connections" component={Connections} />
      <Route path="metadata" component={Metadata}>
        <Route path=":connection" component={Metadata} />
      </Route>
      <Route path="search">
        <IndexRoute component={Search} />
      </Route>
      <Route path="objects" component={Objects} />
      <Route path="*" component={NotFound} />
    </Route>
  </Router>
  ), document.getElementById('app')
);
