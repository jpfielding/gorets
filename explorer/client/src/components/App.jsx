import React from 'react';
import ConfigService from 'services/ConfigService';
import Autocomplete from 'react-autocomplete';
import Configs from 'components/containers/Configs';
import Server from 'components/containers/Server';
import TabSection from 'components/containers/TabSection';

export default class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      configs: [],
      active: {},
      selected: {},
      configAutocompleteField: '',
    };
    this.updateConfigList = this.updateConfigList.bind(this);
    this.createTabList = this.createTabList.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.removeTab = this.removeTab.bind(this);
  }

  componentDidMount() {
    this.updateConfigList();
  }

  updateConfigList() {
    ConfigService
      .getConfigList()
      .then(res => res.json())
      .then((json) => {
        this.setState({ configs: json.result.configs });
      });
  }

  createTabList() {
    const rtn = ['New Config'];
    Object.keys(this.state.active).map(id => (rtn.push(id)));
    return rtn;
  }

  createTabs() {
    const rtn = [<Configs updateCallback={this.updateConfigList} />];
    Object.keys(this.state.active).map(id => (rtn.push(this.state.active[id])));
    return rtn;
  }

  removeTab(tab) {
    const active = Object.assign({}, this.state.active);
    delete active[tab];
    this.setState({ active });
  }

  render() {
    return (
      <div className="helvetica">
        <nav className="pa3 bg-black">
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
              value={this.state.configAutocompleteField}
              inputProps={{
                placeholder: 'Available configs',
                name: 'configs autocomplete',
                id: 'configs-autocomplete',
              }}
              items={(this.state.configs == null ? [] : this.state.configs)}
              shouldItemRender={(item, value) =>
                (item.id.toLowerCase().indexOf(value.toLowerCase()) !== -1)
              }
              onChange={(event, value) => this.setState({ configAutocompleteField: value })}
              onSelect={(value, config) => {
                console.log('Selected', value, 'from Autocomplete');
                const { active } = this.state;
                active[config.id] = (<Server config={config} />);
                this.setState({
                  configAutocompleteField: '',
                  selected: config,
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
