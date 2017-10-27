import React from 'react';
import { withRouter } from 'react-router';
import Metadata from 'components/containers/Metadata';
import Search from 'components/containers/Search';
import Objects from 'components/containers/Objects';
import StorageCache from 'util/StorageCache';
import MetadataService from 'services/MetadataService';
import TabSection from 'components/containers/TabSection';

class Server extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
    location: React.PropTypes.any,
    router: React.PropTypes.any,
    connection: React.PropTypes.any,
  }

  static emptyMetadata = {
    System: {
      'METADATA-RESOURCE': {
        Resource: [],
      },
      SystemDescription: 'Loading metadata...',
      SystemID: 'Loading...',
    },
  };
  constructor(props) {
    super(props);
    this.state = {
      shared: {
        connection: props.connection,
        metadata: Server.emptyMetadata,
        resource: {},
        class: {},
        fields: [],
        data: [],
        tabs: {},
      },
      tab: 0,
    };
    this.getMetadata = this.getMetadata.bind(this);
    this.onMetadataSelected = this.onMetadataSelected.bind(this);
    this.onMetadataDeselected = this.onMetadataDeselected.bind(this);
    this.onDataSelected = this.onDataSelected.bind(this);
    this.onDataDeselected = this.onDataDeselected.bind(this);
    this.onClassSelected = this.onClassSelected.bind(this);
    this.removeTab = this.removeTab.bind(this);
    this.addTab = this.addTab.bind(this);
  }

  componentWillMount() {
    this.getMetadata(m => {
      console.log('setting ', m);
      const shared = this.state.shared;
      shared.metadata = m;
      this.setState({ shared });
    });
  }

  onMetadataSelected(rows) {
    console.log('rows selected:', rows);
    const shared = this.state.shared;
    shared.fields = shared.fields.concat(rows);
    this.setState({ shared });
    console.log('rows:', shared.fields);
  }

  onMetadataDeselected(rows) {
    console.log('rows deselected:', rows);
    const shared = this.state.shared;
    shared.fields = shared.fields.filter(i => rows.map(r => r.row).indexOf(i.row) === -1);
    this.setState({ shared });
    console.log('rows:', shared.fields);
  }

  onDataSelected(rows) {
    console.log('rows selected:', rows);
    const shared = this.state.shared;
    shared.data = shared.fields.concat(rows);
    this.setState({ shared });
    console.log('rows:', shared.data);
  }

  onDataDeselected(rows) {
    console.log('rows deselected:', rows);
    const shared = this.state.shared;
    shared.data = shared.data.filter(i => rows.map(r => r.row).indexOf(i.row) === -1);
    this.setState({ shared });
    console.log('rows:', shared.data);
  }

  onClassSelected(res, cls) {
    console.log('class selected:', res, cls);
    const shared = this.state.shared;
    shared.resource = res;
    shared.class = cls;
    this.setState({ shared });
    this.forceUpdate();
  }

  getMetadata(onFound) {
    const ck = `${this.state.shared.connection.id}-metadata`;
    const md = StorageCache.getFromCache(ck);
    if (md) {
      console.log('loaded metadata from local cache', md);
      onFound(md);
      return;
    }
    console.log('no metadata cached');
    MetadataService
      .get(this.state.shared.connection.id)
      .then(response => response.json())
      .then(json => {
        if (json.error !== null) {
          return;
        }
        console.log('metadata pulled via json request');
        onFound(json.result.Metadata);
        StorageCache.putInCache(ck, json.result.Metadata, 60);
      });
  }

  addTab(key, value) {
    const tabs = Object.assign({}, this.state.tabs);
    tabs[key] = value;
    this.setState({ tabs });
  }

  removeTab(tab) {
    const tabs = Object.assign({}, this.state.tabs);
    delete tabs[tab];
    this.setState({ tabs });
  }

  render() {
    const tabs = this.state.tabs || {};
    const names = Object.keys(tabs);
    const components = Object.keys(tabs).map((key) => tabs[key]);
    names.unshift('Metadata', 'Search', 'Objects');
    components.unshift(
      <Metadata
        shared={this.state.shared}
        onRowsSelected={this.onMetadataSelected}
        onRowsDeselected={this.onMetadataDeselected}
        onClassSelected={this.onClassSelected}
      />,
      <Search
        shared={this.state.shared}
        onRowsSelected={this.onDataSelected}
        onRowsDeselected={this.onDataDeselected}
        addTab={this.addTab}
      />,
      <Objects
        shared={this.state.shared}
        addTab={this.addTab}
      />
    );
    return (
      <TabSection
        names={names}
        components={components}
        enableRemove
        onRemove={this.removeTab}
        removeOffset={3}
      />
    );
  }
}

export default withRouter(Server);
