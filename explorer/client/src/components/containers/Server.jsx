import React from 'react';
import { withRouter } from 'react-router';
import Metadata from 'components/containers/Metadata';
import Search from 'components/containers/Search';
import Objects from 'components/containers/Objects';
import StorageCache from 'util/StorageCache';
import MetadataService from 'services/MetadataService';
import TabSection from 'components/containers/TabSection';
import _ from 'underscore';

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
      SystemDescription: '',
      SystemID: '',
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
      },
      tabs: {},
      errorOut: '',
    };
    this.getMetadata = this.getMetadata.bind(this);
    this.onMetadataSelected = this.onMetadataSelected.bind(this);
    this.onMetadataDeselected = this.onMetadataDeselected.bind(this);
    this.onClassSelected = this.onClassSelected.bind(this);

    this.removeTab = this.removeTab.bind(this);
    this.addTab = this.addTab.bind(this);

    this.errorOut = this.errorOut.bind(this);
  }

  componentWillMount() {
    this.getMetadata(m => {
      console.log('Setting ', m);
      const shared = this.state.shared;
      shared.metadata = m;
      this.setState({ shared });
    });
  }

  onMetadataSelected(rows) {
    console.log('Rows selected:', rows);
    const shared = _.clone(this.state.shared);
    shared.fields = shared.fields.concat(rows);
    this.setState({ shared });
    console.log('Rows:', shared.fields);
  }

  onMetadataDeselected(rows) {
    console.log('Rows deselected:', rows);
    const shared = _.clone(this.state.shared);
    shared.fields = shared.fields.filter(i => rows.map(r => r.row).indexOf(i.row) === -1);
    this.setState({ shared });
    console.log('Rows:', shared.fields);
  }

  onClassSelected(res, cls) {
    console.log('Class selected:', res, cls);
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
      console.log('[SERVER] Loaded metadata from local cache', md);
      onFound(md);
      return;
    }
    const args = {
      extraction: 'COMPACT', // TODO configurable?
    };
    console.log('no metadata cached, pulling', args.extraction);
    MetadataService
      .get(this.state.shared.connection, args)
      .then(response => response.json())
      .then(json => {
        if (json.error !== null) {
          this.errorOut(`[ERROR]  ${json.error}`);
          return;
        }
        console.log('[SERVER] Metadata pulled via json request');
        onFound(json.result.Metadata);
        StorageCache.putInCache(ck, json.result.Metadata, 60);
      });
  }

  errorOut(errorOut) {
    this.setState({ errorOut });
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
        addTab={this.addTab}
      />,
      <Objects
        shared={this.state.shared}
        addTab={this.addTab}
      />
    );
    return (
      <div>
        <div className={`bg-dark-red white br1 pa4 w-100 tc ${this.state.errorOut.length === 0 ? 'dn' : 'db'}`}>
          {this.state.errorOut}
        </div>
        <div className={`${this.state.shared.metadata.System.SystemID.length === 0 ? 'dn' : 'db'}`}>
          <TabSection
            names={names}
            components={components}
            enableRemove
            onRemove={this.removeTab}
            removeOffset={3}
          />
        </div>
        <div className={`loading-wrap ${this.state.shared.metadata.System.SystemID.length !== 0 ? 'dn' : 'db'}`}>
          <div className="loading">LOADING METADATA</div>
        </div>
      </div>
    );
  }
}

export default withRouter(Server);
