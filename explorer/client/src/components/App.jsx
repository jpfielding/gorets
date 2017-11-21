import React from 'react';
import ConfigService from 'services/ConfigService';
import Autocomplete from 'react-autocomplete';
import Connections from 'components/containers/Connections';
import StorageCache from 'util/StorageCache';
import Server from 'components/containers/Server';
import TabSection from 'components/containers/TabSection';
import NewUrl from 'components/containers/NewUrl';
import _ from 'underscore';

export default class App extends React.Component {

  constructor(props) {
    super(props);
    const stored = StorageCache.getFromCache('configs');
    this.state = {
      connections: [],
      configList: (config.configURLS == null ? [] : config.configURLS),
      storedConfigList: (stored == null ? [] : stored),
      active: {},
      selected: {},
      connectionAutocompleteField: '',
      configAutocompleteField: '',
      configURL: '',
      popout: null,
      newURL: {
        name: 'Hey',
        url: '',
      },
    };

    this.addTab = this.addTab.bind(this);
    this.submitConfigURL = this.submitConfigURL.bind(this);
    this.submitNewUrl = this.submitNewUrl.bind(this);
    this.createTabList = this.createTabList.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.removeTab = this.removeTab.bind(this);
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

  updateConfigURL(configURL) {
    this.setState({ configURL });
  }

  handleURLKeyPress(e) {
    if (e.keyCode === 13) {
      this.submitConfigURL();
    }
  }

  submitConfigURL() {
    console.log('Pulling new source list from ', this.state.configURL);
    ConfigService
      .getConfigList(this.state.configURL)
      .then(res => res.json())
      .then((json) => {
        console.log(json);
        this.setState({ connections: json.result.configs });
      });
  }

  submitNewUrl(e) {
    console.log('Submit', e);
    ConfigService
      .getConfigList(e.url)
      .then(res => res.json())
      .then((json) => {
        console.log(json);
        const storedConfigList = _.clone(this.state.storedConfigList);
        storedConfigList.push(e);
        StorageCache.putInCache('configs', storedConfigList, 720);
        this.setState({ connections: json.result.configs, storedConfigList });
      });
  }

  render() {
    const fullList = this.state.configList.concat(this.state.storedConfigList);
    return (
      <div className="helvetica">
        {this.state.popout ? (
          <div className="cover" key="popout">
            {this.state.popout}
          </div>
        ) : null}
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
          <div
            style={{
              position: 'relative',
              zIndex: '100',
              width: '400px',
              float: 'right',
              marginTop: '10px',
              marginRight: '80px',
            }}
            className="titleInput"
          >
            <Autocomplete
              value={this.state.configAutocompleteField}
              inputProps={{
                placeholder: 'Source URLS',
                name: 'config autocomplete',
                id: 'config-autocomplete',
              }}
              items={fullList}
              onChange={(event, value) => this.setState({ configAutocompleteField: value })}
              onSelect={(value) => {
                console.log('Selected', value, 'from Config Autocomplete');
                this.setState({
                  configAutocompleteField: '',
                  configURL: value,
                }, () => this.submitConfigURL());
              }}
              renderMenu={(items, value, style) => (
                <div style={{ ...style, padding: '0px', position: 'fixed' }}>
                  <div style={{ padding: '0px' }} children={items} />
                  <div className="titleBottom" >
                    <button
                      onClick={() => {
                        const popout = (
                          <NewUrl submit={this.submitNewUrl} close={() => this.setState({ popout: null })} />
                        );
                        this.setState({ popout });
                      }}
                    > Create new Source URL </button>
                  </div>
                </div>
                )}
              getItemValue={(item) => item.url}
              renderItem={(item, isHighlighted) => (
                <div
                  style={isHighlighted ? { backgroundColor: '#e8e8e8' } : { backgroundColor: 'white' }}
                  key={item.name}
                  className="clickable"
                >
                  {`${item.name}: `}
                  <span className="moon-gray">
                    {item.url}
                  </span>
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
