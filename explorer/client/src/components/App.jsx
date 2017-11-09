import React from 'react';
import ConnectionService from 'services/ConnectionService';
import Autocomplete from 'react-autocomplete';
import Connections from 'components/containers/Connections';
import Server from 'components/containers/Server';
import TabSection from 'components/containers/TabSection';
import _ from 'underscore';

export default class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      connections: [],
      active: {},
      selected: {},
      connectionAutocompleteField: '',
    };
    this.addTab = this.addTab.bind(this);
    this.updateConnectionList = this.updateConnectionList.bind(this);
    this.createTabList = this.createTabList.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.removeTab = this.removeTab.bind(this);
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

  createTabList() {
    const rtn = ['New Connection'];
    Object.keys(this.state.active).map(id => (rtn.push(id)));
    return rtn;
  }

  createTabs() {
    const rtn = [<Connections addTab={this.addTab} />];
    Object.keys(this.state.active).map(id => (rtn.push(this.state.active[id])));
    return rtn;
  }

  removeTab(tab) {
    const active = _.clone(this.state.active);
    delete active[tab];
    this.setState({ active });
  }

  addTab(connection) {
    const active = _.clone(this.state.active);
    active[connection.id] = <Server connection={connection} />;
    this.setState({ active });
  }

  render() {
    return (
      <div className="helvetica">
        <nav className="pa3 bg-black" style={{ paddingBottom: '0px' }}>
          <h1 className=" fw2 red b f4 dib mr3 nonclickable"> RETS Explorer </h1>
          <div
            style={{
              position: 'relative',
              zIndex: '100',
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
                console.log('Selected', value, 'from Autocomplete');
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
                  className="clickable"
                >
                  {item.id}
                </div>
              )}
            />
          </div>
        </nav>
        <TabSection
          className="customTabElementA"
          names={this.createTabList()}
          components={this.createTabs()}
          enableRemove
          onRemove={this.removeTab}
          removeOffset={1}
        />
      </div>
    );
  }
}
