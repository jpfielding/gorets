import React from 'react';
import { Link } from 'react-router';

export default class App extends React.Component {

  static propTypes = {
    children: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      connections: [{
        id: 'aaarnc',
      }, {
        id: 'test',
      }],
    };
  }

  render() {
    return (
      <div className="helvetica">
        <nav className="pa3 bg-dark-blue">
          <Link to="/" title="Home" className="link fw2 moon-gray b f4 dib mr3">RETS Explorer</Link>
          <Link to="/connections" title="Connections" className="link moon-gray f6 dib mr3">Connections</Link>
          <Link to="/explorer" title="Explorer" className="link moon-gray f6 dib mr3">Explorer</Link>
        </nav>
        <div className="pv2 pl3 bb">
          <span>Connections: </span>
          {this.state.connections.map(connection =>
            <Link to={`/explorer/${connection.id}`} title={connection.id} className="link f6 dib mr3">
              {connection.id}
            </Link>
          )}
        </div>
        {this.props.children}
      </div>
    );
  }
}
