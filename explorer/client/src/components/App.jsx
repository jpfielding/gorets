import React from 'react';
import { Link } from 'react-router';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import ConnectionService from 'services/ConnectionService';
import Autocomplete from 'react-autocomplete';
import Connections from 'components/containers/Connections';
import Server from 'components/containers/Server';

export default class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      connections: [],
      active: {},
      selected: {},
      connectionAutocompleteField: '',
    };
    this.updateConnectionList = this.updateConnectionList.bind(this);
  }

  componentDidMount() {
    this.updateConnectionList();
  }

  updateConnectionList() {
    ConnectionService
      .getConnectionList()
      .then(res => res.json())
      .then((json) => {
        this.setState({ connections: json.result.connections });
      });
  }

  render() {
    return (
      <div className="helvetica">
        <nav className="pa3 bg-black">
          <Link
            to="/"
            title="Home"
            className="link fw2 red b f4 dib mr3"
          >
            RETS Explorer
          </Link>
          <div
            style={{
              position: 'relative',
              zIndex: '10',
              display: 'inline-block',
              width: '400px',
            }}
            className="titleInput"
          >
            <Autocomplete
              value={this.state.connectionAutocompleteField}
              inputProps={{
                placeholder: 'Available connections',
                name: 'connections autocomplete',
                id: 'connections-autocomplete',
              }}
              items={(this.state.connections == null ? [] : this.state.connections)}
              shouldItemRender={(item, value) =>
                (item.id.toLowerCase().indexOf(value.toLowerCase()) !== -1)
              }
              onChange={(event, value) => this.setState({ connectionAutocompleteField: value })}
              onSelect={(value, connection) => {
                console.log('selected', value);
                const { active } = this.state;
                active[connection.id] = (<Server connection={connection} />);
                this.setState({
                  connectionAutocompleteField: '',
                  selected: connection,
                  active,
                });
              }}
              sortItems={(a, b) => (a.id.toLowerCase() <= b.id.toLowerCase() ? -1 : 1)}
              getItemValue={(item) => item.id}
              renderItem={(item, isHighlighted) => (
                <div
                  style={isHighlighted ? { backgroundColor: '#e8e8e8' } : { backgroundColor: 'white' }}
                  key={item.id}
                >
                  {item.id}
                </div>
              )}
            />
          </div>
        </nav>
        <Tabs>
          <TabList>
            <Tab>New Connection</Tab>
            {Object.keys(this.state.active).map(id =>
              (<Tab>{id}</Tab>)
            )}
          </TabList>
          <TabPanel><Connections updateCallback={this.updateConnectionList} /></TabPanel>
          {Object.keys(this.state.active).map(id =>
            (<TabPanel>{this.state.active[id]}</TabPanel>)
          )}
        </Tabs>
      </div>
    );
  }
}
