import React from 'react';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';
import ReactDOM from 'react-dom';
import App from 'components/App';
import Home from 'components/containers/Home';
import Connections from 'components/containers/Connections';
import Explorer from 'components/containers/Explorer';
import Search from 'components/containers/Search';
import SearchMetadata from 'components/containers/SearchMetadata';
import NotFound from 'components/containers/NotFound';

import 'styles/app.css';

ReactDOM.render((
  <Router history={hashHistory}>
    <Route path="/" component={App}>
      <IndexRoute component={Home} />
      <Route path="connections" component={Connections} />
      <Route path="explorer" component={Explorer}>
        <Route path=":connection" component={Explorer} />
      </Route>
      <Route path="search" component={Search}>
        <Route path="metadata" component={SearchMetadata} />
      </Route>
      <Route path="*" component={NotFound} />
    </Route>
  </Router>
  ), document.getElementById('app')
);
