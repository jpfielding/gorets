import React from 'react';
import { withRouter } from 'react-router';
import Metadata from 'components/containers/Metadata';
import Search from 'components/containers/Search';
import Objects from 'components/containers/Objects';
import Explore from 'components/containers/Explore';
import StorageCache from 'util/StorageCache';
import MetadataService from 'services/MetadataService';
import TabSection from 'components/containers/TabSection';
import ConnectionForm from 'components/containers/ConnectionForm';
import Wireloger from 'components/containers/Wireloger';
import _ from 'underscore';

const Base64 = require('js-base64').Base64;

class Server extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
    connection: React.PropTypes.any,
    init: React.PropTypes.any,
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
        connection: props.connection.data,
        metadata: Server.emptyMetadata,
        resource: {},
        class: {},
        fields: [],
        source: props.connection.config,
      },
      wirelog: [],
      tabs: [],
      errorOut: '',
      args: {
        extraction: 'COMPACT',
        oldest: 360,
      },
    };

    this.getMetadata = this.getMetadata.bind(this);
    this.onMetadataSelected = this.onMetadataSelected.bind(this);
    this.onMetadataDeselected = this.onMetadataDeselected.bind(this);
    this.onClassSelected = this.onClassSelected.bind(this);

    this.updateConnection = this.updateConnection.bind(this);

    this.removeTab = this.removeTab.bind(this);
    this.addTab = this.addTab.bind(this);

    this.errorOut = this.errorOut.bind(this);
    this.pushWirelog = this.pushWirelog.bind(this);
  }

  componentWillMount() {
    this.getMetadata((m, log, extra) => {
      const shared = _.clone(this.state.shared);
      const wirelog = _.clone(this.state.wirelog);
      shared.metadata = m;
      wirelog.unshift({ tag: 'Metadata', log, extra });
      this.setState({ shared, wirelog });
    });
  }

  onMetadataSelected(rows) {
    console.log('rows selected:', rows);
    const shared = _.clone(this.state.shared);
    shared.fields = shared.fields.concat(rows);
    this.setState({ shared });
    console.log('rows:', shared.fields);
  }

  onMetadataDeselected(rows) {
    console.log('rows deselected:', rows);
    const shared = _.clone(this.state.shared);
    shared.fields = shared.fields.filter(i => rows.map(r => r.row).indexOf(i.row) === -1);
    this.setState({ shared });
    console.log('rows:', shared.fields);
  }

  onClassSelected(res, cls) {
    console.log('class selected:', res, cls);
    const shared = _.clone(this.state.shared);
    shared.resource = res;
    shared.class = cls;
    this.setState({ shared });
  }

  getMetadata(onFound) {
    const ck = `${this.state.shared.source}-${this.state.shared.connection.id}-metadata`;
    const md = StorageCache.getFromCache(ck);
    if (md) {
      console.log('loaded metadata from local cache', md);
      onFound(md, 'No wirelog available, metadata was loaded from local cashe.\n' +
        'To force a new pull press \'Update Changes\' in the \'Conncetion Config\' pannel.',
         { type: 'Info' });
      return;
    }
    const args = this.state.args;
    console.log('no metadata cached, pulling', args.extraction);
    MetadataService
      .get(this.state.shared.connection, args)
      .then(response => response.json())
      .then(json => {
        if (json.error !== null) {
          this.errorOut(json.error);
          return;
        }
        console.log('metadata pulled via json request', json);
        if (json.result.wirelog) {
          onFound(json.result.Metadata, Base64.decode(json.result.wirelog));
        } else {
          onFound(json.result.Metadata, 'Metadata recived without wirelog.\n' +
            'This is mostly because it was pulled from a cashe in the provider and no the source itself.\n' +
            'To force a new pull press \'Update Changes\' in the \'Conncetion Config\' pannel.',
             { type: 'Info' });
        }
        StorageCache.putInCache(ck, json.result.Metadata, 60);
      });
  }

  pushWirelog(e) {
    const wirelog = _.clone(this.state.wirelog);
    wirelog.unshift(e);
    this.setState({ wirelog });
  }

  updateConnection(connection, args) {
    const sck = `${this.state.shared.source}-${this.state.shared.connection.id}-search-history`;
    const ock = `${this.state.shared.source}-${this.state.shared.connection.id}-object-history`;
    const mck = `${this.state.shared.source}-${this.state.shared.connection.id}-metadata`;
    StorageCache.remove(sck);
    StorageCache.remove(ock);
    StorageCache.remove(mck);
    const shared = _.clone(this.state.shared);
    shared.connection = connection;
    shared.metadata = Server.emptyMetadata;
    this.setState({ shared, args, errorOut: '' }, () => {
      this.getMetadata((m, log, extra) => {
        console.log('Setting ', m);
        shared.metadata = m;
        shared.wirelog.unshift({ tag: 'Metadata', log, extra });
        this.setState({ shared });
      });
    });
  }

  errorOut(errorOut) {
    this.setState({ errorOut });
  }

  addTab(key, value) {
    const tabs = _.clone(this.state.tabs);
    tabs.push({
      id: key,
      page: value,
    });
    this.setState({ tabs });
  }

  removeTab(t) {
    const tabs = _.clone(this.state.tabs);
    tabs.forEach((tab, i) => {
      if (tab.id === t) {
        tabs.splice(i, 1);
      }
    });
    this.setState({ tabs });
  }

  render() {
    const tabs = _.clone(this.state.tabs);
    const pages = [
      <Metadata
        shared={this.state.shared}
        onRowsSelected={this.onMetadataSelected}
        onRowsDeselected={this.onMetadataDeselected}
        onClassSelected={this.onClassSelected}
      />,
      <Search
        shared={this.state.shared}
        addTab={this.addTab}
        pushWirelog={this.pushWirelog}
        init={this.props.init && this.props.init.tab === 'Search' ? this.props.init : null}
      />,
      <Objects
        shared={this.state.shared}
        addTab={this.addTab}
        pushWirelog={this.pushWirelog}
        init={this.props.init && this.props.init.tab === 'Objects' ? this.props.init : null}
      />,
      <Explore
        shared={this.state.shared}
      />,
      <Wireloger wirelog={this.state.wirelog} />,
    ];
    tabs.unshift(
      {
        id: 'Metadata',
        page: pages[0],
      },
      {
        id: 'Search',
        page: pages[1],
      },
      {
        id: 'Objects',
        page: pages[2],
      },
      {
        id: 'Explore',
        page: pages[3],
      },
      {
        id: 'Wirelog',
        page: pages[4],
      }

    );
    return (
      <div>
        <div className="fr">
          <div className="customHoverSection">
            <div className="fr ma-3 customHoverBar"> Connection Config </div>
            <div className="customHoverBody">
              <ConnectionForm
                updateConnection={this.updateConnection}
                connection={this.state.shared.connection}
                location={this.props.location}
              />
            </div>
          </div>
        </div>
        <div className={`bg-dark-red white br1 pa4 w-100 tc ${this.state.errorOut.length === 0 ? 'dn' : 'db'}`}>
          {this.state.errorOut}
        </div>
        <div className={`${this.state.shared.metadata.System.SystemID.length === 0 ? 'dn' : 'db'}`}>
          <TabSection
            className="customTabElementB"
            components={tabs}
            enableRemove
            onRemove={this.removeTab}
            removeOffset={pages.length}
            initID={this.props.init ? this.props.init.tab : null}
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
