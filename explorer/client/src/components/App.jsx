import React from 'react';
import { Link } from 'react-router';

export default class App extends React.Component {

  static propTypes = {
    children: React.PropTypes.any,
  }

  render() {
    return (
      <div className="helvetica">
        <nav className="pa3 bg-dark-blue">
          <Link to="/" title="Home" className="link fw2 moon-gray b f4 dib mr3">RETS Explorer</Link>
          <Link to="/login" title="Login" className="link moon-gray f6 dib mr3">Login</Link>
        </nav>
        {this.props.children}
      </div>
    );
  }
}
