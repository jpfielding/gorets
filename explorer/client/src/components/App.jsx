import React from 'react';
import ConfigService from 'services/ConfigService';
import Autocomplete from 'react-autocomplete';
import Connections from 'components/containers/Connections';
import StorageCache from 'util/StorageCache';
import Server from 'components/containers/Server';
import TabSection from 'components/containers/TabSection';
import NewUrl from 'components/containers/NewUrl';
import _ from 'underscore';

const Base64 = require('js-base64').Base64;

export default class App extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
  }

  constructor(props) {
    super(props);

    const stored = StorageCache.getFromCache('configs');

    /*
      State Setup
      connections:
        {
          name: <name of config connection is from>
          data: <actual connection data>
        }
      activeTabs:
        {
          id: <unique id>
          [name: <overrides the id with a diffrent displayed text>]
          page: <component object to be rendered when tab is active>
          [tags:[{ name: <displayed name>, color: <displayed background color> }, ...]]
        }
      connectionAutocompleteField: ...
      configAutocompleteField: ...
      popout: <component holder for an overlaying component>
      configs: {
        active {
          [<config source id>: <color>]
        }
        server: {
          [{
            url: <url>,
            name: <name for the source>,
          }, ...]
        }
        stored: {
          [{
            url: <url>,
            name: <name for the source>,
          }, ...]
        }
      }
      color: {
        ...
      }
    */

    /*
      TODO Imporove color support to allow lighter colors with black text
      Currently only supports well colors with white text
    */
    this.addSimpleTab = this.addSimpleTab.bind(this);

    this.state = {
      connections: [],
      activeTabs: [
        {
          id: 'newcon',
          name: 'New Connection',
          page: (<Connections addTab={this.addSimpleTab} />),
          idprefix: 'newcon',
        },
      ],
      connectionAutocompleteField: '',
      connectionErrorOut: '',
      configAutocompleteField: '',
      configErrorOut: '',
      popout: null,
      configs: {
        active: {},
        server: (config.configURLS == null ? [] : config.configURLS),
        stored: (stored == null ? [] : stored),
      },
      color: {
        init: [
          '#03A9F4',
          '#673AB7',
          '#4CAF50',
          '#F44336',
        ],
        available: [],
      },
      init: null,
    };

    this.addSimpleTabWithInit = this.addSimpleTabWithInit.bind(this);
    this.addFullTab = this.addFullTab.bind(this);
    this.addDirectTab = this.addDirectTab.bind(this);

    this.submitNewConfig = this.submitNewConfig.bind(this);
    this.submitConfig = this.submitConfig.bind(this);
    this.removeTab = this.removeTab.bind(this);

    this.getNewColor = this.getNewColor.bind(this);
    this.getColor = this.getColor.bind(this);

    this.selectConnection = this.selectConnection.bind(this);

    this.renderConnectionAutocomplete = this.renderConnectionAutocomplete.bind(this);
    this.renderConnectionItem = this.renderConnectionItem.bind(this);
    this.renderConnectionMenu = this.renderConnectionMenu.bind(this);

    this.renderConfigAutocomplete = this.renderConfigAutocomplete.bind(this);
    this.renderConfigItem = this.renderConfigItem.bind(this);
    this.renderConfigMenu = this.renderConfigMenu.bind(this);
  }

  componentWillMount() {
    if (this.props.location.query.s != null) {
      if (this.props.location.query.i != null) {
        const init = JSON.parse(Base64.decode(this.props.location.query.i));
        this.addSimpleTabWithInit(JSON.parse(Base64.decode(this.props.location.query.s)), init);
        this.setState({ init });
      } else {
        this.addSimpleTab(JSON.parse(Base64.decode(this.props.location.query.s)));
      }
    }
  }

  getNewColor() {
    if (this.state.color.available.length === 0) {
      this.state.color.available = _.clone(this.state.color.init);
    }
    return this.state.color.available.pop();
  }

  getColor(configSourceId) {
    return this.state.configs.active[configSourceId];
  }

  getConfig(source, target) {
    let rtn = null;
    source.server.forEach((e) => {
      if (rtn !== null) return;
      if (e.name === target.name && e.url === target.url) {
        rtn = e;
      }
    });
    source.stored.forEach((e) => {
      if (rtn !== null) return;
      if (e.name === target.name && e.url === target.url) {
        rtn = e;
      }
    });
    return rtn;
  }

  removeTab(tab) {
    const activeTabs = _.clone(this.state.activeTabs);
    console.log(tab);
    activeTabs.forEach((e, i) => {
      if (e.id === tab) {
        activeTabs.splice(i, 1);
        console.log('FOUND');
      }
    });
    this.setState({ activeTabs });
  }

  addSimpleTab(connection) {
    const activeTabs = _.clone(this.state.activeTabs);
    activeTabs.push({
      id: connection.id,
      idprefix: connection.id.replace(/\s|:/gi, ''),
      page: (<Server
        connection={{ config: 'simplecon', data: connection }}
        location={this.props.location}
        idprefix={connection.id.replace(/\s|:/gi, '')}
      />),
    });
    this.setState({ activeTabs });
  }

  addSimpleTabWithInit(connection, init) {
    const activeTabs = _.clone(this.state.activeTabs);
    activeTabs.push({
      id: connection.id,
      idprefix: connection.id.replace(/\s|:/gi, ''),
      page: (<Server
        connection={{ config: 'simplecon', data: connection }}
        location={this.props.location}
        init={init}
        idprefix={connection.id.replace(/\s|:/gi, '')}
      />),
    });
    this.setState({ activeTabs });
  }

  addFullTab(connection, configID, name) {
    const activeTabs = _.clone(this.state.activeTabs);
    activeTabs.push({
      id: configID + name,
      idprefix: (configID + name).replace(/\s|:/gi, ''),
      name,
      tags: [{ name: configID, color: this.getColor(configID) }],
      page: (<Server
        connection={connection}
        location={this.props.location}
        idprefix={(configID + name).replace(/\s|:/gi, '')}
      />),
    });
    this.setState({ activeTabs });
  }

  addDirectTab(newTab) {
    const activeTabs = _.clone(this.state.activeTabs);
    activeTabs.push({ newTab });
    this.setState({ activeTabs });
  }

  submitNewConfig(e) {
    this.setState({ configErrorOut: '' });
    const fullList = this.state.configs.server.concat(this.state.configs.stored);
    if (fullList.find((el) => (el.name === e.name))) {
      this.setState({ configErrorOut: `Config name ${e.name} already in use` });
      return;
    }
    this.submitConfig(e);
  }

  submitConfig(e) {
    console.log('Submiting Config.', e);
    this.setState({ configErrorOut: '' }, () => {
      if (!e.name || !e.url) {
        this.setState({ configErrorOut: 'Invalid config submited, both Name and URL fields are required' });
        return;
      }

      const activeIDs = Object.keys(this.state.configs.active);
      if (activeIDs.indexOf(e.name) > -1) {
        this.setState({ configErrorOut: `Config name ${e.name} already in use` });
        return;
      }

      ConfigService
        .getConfigList(e.url)
        .then(res => {
          console.log(res);
          return res.json();
        })
        .then((json) => {
          console.log('Config Responce:', json);
          const configs = _.clone(this.state.configs);

          if (json.error) {
            console.error(json.error);
            this.getConfig(configs, e).failed = true;
            this.setState({ configErrorOut: 'Error with config response', configs });
            return;
          }

          configs.active[e.name] = this.getNewColor();

          if (!_.contains(configs.stored, e) && !_.contains(configs.server, e)) {
            configs.stored.push(e);
            StorageCache.putInCache('configs', configs.stored, 720);
          }

          const r = json.result.configs.map((el) => (
            {
              config: e.name,
              data: el,
            }
          ));
          const connections = _.clone(this.state.connections).concat(r);
          this.setState({ connections, configs });
        })
        .catch((err) => {
          console.error(err);
          const configs = _.clone(this.state.configs);
          this.getConfig(configs, e).failed = true;
          this.setState({ configErrorOut: 'Invalid config', configs });
          return;
        });
    });
  }

  selectConnection(value, connection) {
    console.log('Selected', value, 'from Autocomplete', connection);
    this.setState({ connectionErrorOut: '' }, () => {
      const { activeTabs } = this.state;
      let unique = true;
      activeTabs.forEach((e) => {
        const id = connection.config + value;
        console.log(e.id, id);
        if (e.id === id) {
          unique = false;
        }
      });
      if (unique) {
        this.addFullTab(connection, connection.config, value);
        return;
      }
      this.setState({ connectionErrorOut: 'Failed to create tab, Does that tab already exist?' });
    });
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
          count={(this.state.connections == null ? 0 : this.state.connections.length)}
          getItemValue={(item) => item.data.id}
          onChange={(event, value) => this.setState({ connectionAutocompleteField: value })}
          onSelect={this.selectConnection}
          sortItems={(a, b) => {
            if (a.data.id.toLowerCase() === b.data.id.toLowerCase()) {
              const k = Object.keys(this.state.configs.active);
              return k.indexOf(a.name) - k.indexOf(b.name);
            }
            return (a.data.id.toLowerCase() <= b.data.id.toLowerCase() ? -1 : 1);
          }}
          shouldItemRender={(item) => (item.data.id.includes(this.state.connectionAutocompleteField))}
          renderMenu={this.renderConnectionMenu}
          renderItem={this.renderConnectionItem}
        />
        { this.state.connectionErrorOut !== '' ?
          <div className="error-out">
            {this.state.connectionErrorOut}
          </div>
          : null
        }
      </div>
    );
  }

  renderConfigAutocomplete() {
    const fullList = this.state.configs.server.concat(this.state.configs.stored);
    const activeCount = Object.keys(this.state.configs.active).length;
    return (
      <div className="navbarAutocomplete">
        <div className={`tag ${activeCount === 0 ? 'alert' : null}`}>
          {Object.keys(this.state.configs.active).length}
        </div>
        <Autocomplete
          value={this.state.configAutocompleteField}
          inputProps={{
            placeholder: 'Source URLS',
            name: 'config autocomplete',
            id: 'config-autocomplete',
          }}
          items={fullList}
          count={fullList.length}
          getItemValue={(item) => item.name}
          onChange={(event, value) => this.setState({ configAutocompleteField: value })}
          onSelect={(value, item) => {
            console.log('Selected', value, 'from Config Autocomplete');
            this.setState({
              configAutocompleteField: '',
            }, () => this.submitConfig(item));
          }}
          shouldItemRender={(item) => (item.name.includes(this.state.configAutocompleteField))}
          renderMenu={this.renderConfigMenu}
          renderItem={this.renderConfigItem}
        />
        { this.state.configErrorOut !== '' ?
          <div className="error-out">
            {this.state.configErrorOut}
          </div>
          : null
        }
      </div>
    );
  }

  renderConnectionMenu(items, value, style) {
    return (
      <div style={{ ...style, padding: '0px', position: 'fixed' }}>
        <div style={{ padding: '0px' }} className="titleSelectBox" children={items} id="connection-menu" />
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
        key={item.config + item.data.id}
        className="clickable"
        id={(item.config + item.data.id).replace(/\s|:/gi, '')}
      >
        <div style={{ backgroundColor: this.state.configs.active[item.config] }} className="activeStartTag" />
        {item.data.id}
      </div>
    );
  }

  renderConfigMenu(items, value, style) {
    return (
      <div style={{ ...style, padding: '0px', position: 'fixed', width: '400px' }}>
        <div style={{ padding: '0px' }} className="titleSelectBox" children={items} id="config-menu" />
        <div className="titleBottom" >
          <button
            className="default"
            onClick={() => {
              const popout = (
                <NewUrl submit={this.submitNewConfig} close={() => this.setState({ popout: null })} />
              );
              this.setState({ popout });
            }}
            id={'add-source-url'}
          > + Source URL </button>
        </div>
      </div>
    );
  }

  renderConfigItem(item, isHighlighted) {
    return (
      <div
        id={item.name.replace(/\s|:/gi, '')}
        style={isHighlighted ? { backgroundColor: '#e8e8e8' } : { backgroundColor: 'white' }}
        key={item.name}
        className="clickable"
      >
        <div className={item.failed ? 'o-50 flex' : 'flex'} >
          { (Object.keys(this.state.configs.active).indexOf(item.name) > -1) ?
            <div
              style={{
                backgroundColor: this.state.configs.active[item.name],
              }}
              className="activeStartTag"
            />
            : null
          }
          { item.failed ?
            <div className="failed-tag">
              X
            </div>
            : null
          }
          <div style={{ flex: '1', padding: '0px' }}>
            {`${item.name}: `}
            <span className="moon-gray">
              {item.url}
            </span>
          </div>
        </div>
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
          components={this.state.activeTabs}
          enableRemove
          onRemove={this.removeTab}
          removeOffset={1}
          initID={this.state.init != null ? this.state.init.id : null}
        />

      </div>
    );
  }
}
