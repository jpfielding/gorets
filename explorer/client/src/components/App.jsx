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

    /*
      State Setup
      connections: array of objects of format
      {
        name: <name of config connection is from
        data: <actual connection data>
      }
      used for holding list of currently avaiable connections and their source
      activeTabs: array of objects of format
      {
        id: <id to be displayed as title header>
        [title: <overrides the id with a diffrent id>]
        page: <component object to be rendered when tab is active>
        [tabs: ... ]
      }
      connectionAutocompleteField: autocomplete field for Connections autocomplete
      component
      configAutocompleteField: autocomplete field for Configs autocomplete component
      popout: <component holder for an overlaying component>
      Use to force a component infront of the current display and block action to
      other components
    */
    this.state = {
      connections: [],
      activeTabs: [],
      connectionAutocompleteField: '',
      configAutocompleteField: '',
      popout: null,
      configs: {
        active: {},
        server: (config.configURLS == null ? [] : config.configURLS),
        stored: (stored == null ? [] : stored),
      },
      color: {
        init: [
          '#F44336', '#E91E63', '#9C27B0',
          '#673AB7', '#3F51B5', '#2196F3',
          '#03A9F4', '#00BCD4', '#009688',
          '#4CAF50', '#FF5722',
        ],
        available: [],
      },
    };

    this.addTab = this.addTab.bind(this);
    this.submitConfig = this.submitConfig.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.removeTab = this.removeTab.bind(this);
    this.getNewColor = this.getNewColor.bind(this);

    this.selectConnection = this.selectConnection.bind(this);

    this.renderConnectionAutocomplete = this.renderConnectionAutocomplete.bind(this);
    this.renderConnectionItem = this.renderConnectionItem.bind(this);
    this.renderConnectionMenu = this.renderConnectionMenu.bind(this);

    this.renderConfigAutocomplete = this.renderConfigAutocomplete.bind(this);
  }

  getNewColor() {
    if (this.state.color.available.length === 0) {
      this.state.color.available = _.clone(this.state.color.init);
    }
    const index = Math.floor(Math.random() * this.state.color.available.length);
    const color = this.state.color.available[index];
    this.state.color.available.splice(index, 1);
    return color;
  }

  createTabs() {
    const rtn = [{ id: 'newcon', name: 'New Connection', page: (<Connections addTab={this.addTab} />) }];
    this.state.activeTabs.forEach(e => {
      console.log(e);
      if (e.config) {
        rtn.push(
          {
            id: e.config + e.id,
            name: e.id,
            tags: [{ name: e.config, color: this.state.configs.active[e.config] }],
            page: e.page,
          }
        );
      } else {
        rtn.push(
          {
            id: e.id,
            page: e.page,
          }
        );
      }
    });
    return rtn;
  }

  removeTab(tab) {
    const activeTabs = _.clone(this.state.activeTabs);
    console.log(tab);
    activeTabs.forEach((e, i) => {
      if (e.config + e.id === tab) {
        activeTabs.splice(i, 1);
        console.log('FOUND');
      }
    });
    this.setState({ activeTabs });
  }

  addTab(connection) {
    const activeTabs = _.clone(this.state.activeTabs);
    activeTabs.push({
      id: connection.id,
      page: (<Server connection={{ name: 'newcon', data: connection }} />) });
    this.setState({ activeTabs });
  }

  submitConfig(e) {
    console.log('Submiting Config', e);

    const activeIDs = Object.keys(this.state.configs.active);
    if (activeIDs.indexOf(e.name) > -1) {
      console.log('Rejected. ', e.name, ' is a already an active source name');
      return;
    }

    ConfigService
      .getConfigList(e.url)
      .then(res => res.json())
      .then((json) => {
        console.log(json);

        const configs = _.clone(this.state.configs);
        configs.active[e.name] = this.getNewColor();

        if (!_.contains(configs.stored, e) && !_.contains(configs.server, e)) {
          configs.stored.push(e);
          StorageCache.putInCache('configs', configs.stored, 720);
        }

        const r = json.result.configs.map((el) => (
          {
            name: e.name,
            data: el,
          }
          ));
        const connections = _.clone(this.state.connections).concat(r);

        this.setState({ connections, configs });
      });
  }

  selectConnection(value, connection) {
    console.log('Selected', value, 'from Autocomplete', connection);
    const { activeTabs } = this.state;
    let unique = true;
    activeTabs.forEach((e) => {
      if (e.id === value && e.config === connection.name) {
        unique = false;
      }
    });
    if (unique) {
      activeTabs.push(
        {
          config: connection.name,
          id: value,
          page: (<Server connection={connection} />),
        }
       );
      this.setState({
        connectionAutocompleteField: '',
        selected: connection,
        activeTabs,
      });
    }
  }

  renderConnectionAutocomplete() {
    return (
      <div className="navbarAutocomplete">
        <Autocomplete
          value={this.state.connectionAutocompleteField}
          inputProps={{
            placeholder: 'Available connections',
            name: 'connections autocomplete',
            id: 'connections-autocomplete',
          }}
          items={(this.state.connections == null ? [] : this.state.connections)}
          onChange={(event, value) => this.setState({ connectionAutocompleteField: value })}
          onSelect={this.selectConnection}
          shouldItemRender={(item) => (item.data.id.includes(this.state.connectionAutocompleteField))}
          sortItems={(a, b) => {
            if (a.data.id.toLowerCase() === b.data.id.toLowerCase()) {
              const k = Object.keys(this.state.configs.active);
              return k.indexOf(a.name) - k.indexOf(b.name);
            }
            return (a.data.id.toLowerCase() <= b.data.id.toLowerCase() ? -1 : 1);
          }}
          getItemValue={(item) => item.data.id}
          renderMenu={this.renderConnectionMenu}
          renderItem={this.renderConnectionItem}
        />
      </div>
    );
  }

  renderConfigAutocomplete() {
    const fullList = this.state.configs.server.concat(this.state.configs.stored);
    return (
      <div className="navbarAutocomplete">
        <Autocomplete
          value={this.state.configAutocompleteField}
          inputProps={{
            placeholder: 'Source URLS',
            name: 'config autocomplete',
            id: 'config-autocomplete',
          }}
          items={fullList}
          onChange={(event, value) => this.setState({ configAutocompleteField: value })}
          onSelect={(value, item) => {
            console.log('Selected', value, 'from Config Autocomplete');
            this.setState({
              configAutocompleteField: '',
            }, () => this.submitConfig(item));
          }}
          renderMenu={(items, value, style) => (
            <div style={{ ...style, padding: '0px', position: 'fixed' }}>
              <div style={{ padding: '0px' }} className="titleSelectBox" children={items} />
              <div className="titleBottom" >
                <button
                  onClick={() => {
                    const popout = (
                      <NewUrl submit={this.submitConfig} close={() => this.setState({ popout: null })} />
                    );
                    this.setState({ popout });
                  }}
                > Create new Source URL </button>
              </div>
            </div>
            )}
          getItemValue={(item) => item.name}
          shouldItemRender={(item) => (item.name.includes(this.state.configAutocompleteField))}
          renderItem={(item, isHighlighted) => (
            <div
              style={isHighlighted ? { backgroundColor: '#e8e8e8' } : { backgroundColor: 'white' }}
              key={item.name}
              className="clickable flex"
            >
              { (Object.keys(this.state.configs.active).indexOf(item.name) > -1) ?
                <div
                  style={{
                    width: '3px',
                    display: 'inline-block',
                    backgroundColor: this.state.configs.active[item.name],
                    marginRight: '5px',
                  }}
                />
                : null}
              <div style={{ flex: '1', padding: '0px' }}>
                {`${item.name}: `}
                <span className="moon-gray">
                  {item.url}
                </span>
              </div>
            </div>
          )}
        />
      </div>
    );
  }

  renderConnectionMenu(items, value, style) {
    return (
      <div style={{ ...style, padding: '0px', position: 'fixed' }}>
        <div style={{ padding: '0px' }} className="titleSelectBox" children={items} />
        <div className="titleBottom" >
          {Object.keys(this.state.configs.active).map((e) => (
            <div style={{ backgroundColor: this.state.configs.active[e] }} className="activeFullTag">
              {e}
            </div>
          ))}
        </div>
      </div>
    );
  }

  renderConnectionItem(item, isHighlighted) {
    return (
      <div
        style={isHighlighted ? { backgroundColor: '#e8e8e8' } : { backgroundColor: 'white' }}
        key={item.name + item.data.id}
        className="clickable"
      >
        <div style={{ backgroundColor: this.state.configs.active[item.name] }} className="activeStartTag" />
        {item.data.id}
      </div>
    );
  }

  render() {
    return (
      <div className="helvetica">

        {this.state.popout ? (
          <div className="cover" key="popout">
            {this.state.popout}
          </div>
        ) : null}

        <nav className="pa3 bg-black flex" style={{ paddingBottom: '0px' }}>
          <h1 className=" fw2 red b f4 dib mr3 nonclickable"> RETS Explorer </h1>
          {this.renderConnectionAutocomplete()}
          <div style={{ flex: '1' }} />
          {this.renderConfigAutocomplete()}
        </nav>

        <TabSection
          className="customTabElementA"
          components={this.createTabs()}
          enableRemove
          onRemove={this.removeTab}
          removeOffset={1}
        />

      </div>
    );
  }
}
