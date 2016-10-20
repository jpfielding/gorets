import React from 'react';
import { Link } from 'react-router';
import ConnectionService from 'services/ConnectionService';

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

  componentDidMount() {
    ConnectionService
      .getConnectionList()
      .then(res => res.json())
      .then((json) => {
        console.log('Connection List: ', json);
      });
  }

  render() {
    return (
      <div className="helvetica">
        <nav className="pa3 bg-dark-blue">
          <Link
            to="/"
            title="Home"
            activeStyle={{ color: 'white' }}
            className="link fw2 moon-gray b f4 dib mr3"
          >
            RETS Explorer
          </Link>
          <Link
            to="/connections"
            title="Connections"
            activeStyle={{ color: 'white' }}
            className="link moon-gray f6 dib mr3"
          >
            Connections
          </Link>
          <Link
            to="/explorer"
            title="Explorer"
            activeStyle={{ color: 'white' }}
            className="link moon-gray f6 dib mr3"
          >
            Explorer
          </Link>
        </nav>
        <div className="pv2 pl3 bb v-mid flex flex-row align-center">
          <span className="f6 mr3">Connections: </span>
          {this.state.connections.map(connection =>
            <Link to={`/explorer/${connection.id}`} title={connection.id} className="link f6 dib mr2">
              {connection.id}
            </Link>
          )}
        </div>
        {this.props.children}
      </div>
    );
  }
}
