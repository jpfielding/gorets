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
    idprefix: React.PropTypes.any,
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
    const args = props.init ? props.init.args || {
      extraction: 'COMPACT',
      oldest: 360,
    } : {
      extraction: 'COMPACT',
      oldest: 360,
    };
    this.state = {
      shared: {
        connection: props.connection.data,
        metadata: Server.emptyMetadata,
        resource: {},
        class: {},
        fields: [],
        source: props.connection.config,
        args,
      },
      wirelog: [],
      errorOut: '',
      validMetadata: true,
      metadataIssue: '',
      loading: true,
    };

    this.getMetadata = this.getMetadata.bind(this);
    this.onMetadataSelected = this.onMetadataSelected.bind(this);
    this.onMetadataDeselected = this.onMetadataDeselected.bind(this);
    this.onMetadataCleared = this.onMetadataCleared.bind(this);
    this.onClassSelected = this.onClassSelected.bind(this);

    this.updateConnection = this.updateConnection.bind(this);

    this.errorOut = this.errorOut.bind(this);
    this.pushWirelog = this.pushWirelog.bind(this);

    this.validateMetadata = this.validateMetadata.bind(this);
    this.failMetadata = this.failMetadata.bind(this);
    this.successMetadata = this.successMetadata.bind(this);

    this.renderLoading = this.renderLoading.bind(this);
    this.renderError = this.renderError.bind(this);
    this.renderInvalid = this.renderInvalid.bind(this);
    this.renderStandard = this.renderStandard.bind(this);
    this.renderSelect = this.renderSelect.bind(this);
  }

  componentWillMount() {
    this.getMetadata((m, log, extra) => {
      const shared = _.clone(this.state.shared);
      const wirelog = _.clone(this.state.wirelog);
      shared.metadata = m;
      wirelog.unshift({ tag: 'Metadata', log, extra });
      this.setState({ shared, wirelog, loading: false });
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

  onMetadataCleared() {
    console.log('rows deselected:', this.state.shared.fields);
    const shared = _.clone(this.state.shared);
    shared.fields = [];
    this.setState({ shared });
    console.log('rows:', shared.fields);
  }

  onClassSelected(res, cls) {
    console.log('class selected:', res, cls);
    const shared = _.clone(this.state.shared);
    shared.resource = res;
    shared.class = cls;
    shared.fields = [];
    this.setState({ shared });
  }

  getMetadata(onFound) {
    // Check local
    const ck = `${this.state.shared.source}-${this.state.shared.connection.id}-metadata`;
    const md = StorageCache.getFromCache(ck);
    if (md) {
      console.log('loaded metadata from local cache', md);

      const log = 'No wirelog available, metadata was loaded from local cashe.\n' +
        'To force a new pull press \'Update Changes\' in the \'Conncetion Config\' pannel.';
      const extra = { type: 'Info' };

      if (!this.validateMetadata(md)) {
        const shared = _.clone(this.state.shared);
        const wirelog = _.clone(this.state.wirelog);
        shared.metadata = md;
        wirelog.unshift({ tag: 'Metadata', log, extra });
        this.setState({ shared, wirelog });
        onFound(md, log, extra);
        return;
      }

      onFound(md, log, extra);
      return;
    }

    const args = this.state.shared.args;
    console.log('no metadata cached, pulling', args.extraction);

    // Make api request
    MetadataService
      .get(this.state.shared.connection, args)
      .then(response => {
        console.log(response);
        return response.json();
      })
      .then(json => {
        // Error out
        if (json.error && json.error !== null) {
          this.errorOut(json.error);
          return;
        }
        console.log('metadata pulled via json request', json);

        // Determins if the request has a wirelog or not.
        // If not it fills in one for it.
        let log = '';
        let extra;

        if (json.result.wirelog) {
          log = Base64.decode(json.result.wirelog);
        } else {
          log = 'Metadata recived without wirelog.\n' +
            'This is mostly because it was pulled from a cashe in the provider and no the source itself.\n' +
            'To force a new pull press \'Update Changes\' in the \'Conncetion Config\' pannel.';
          extra = { type: 'Info' };
        }

        if (!this.validateMetadata(json.result.Metadata)) {
          const shared = _.clone(this.state.shared);
          const wirelog = _.clone(this.state.wirelog);
          shared.metadata = json.result.Metadata;
          wirelog.unshift({ tag: 'Metadata', log, extra });
          this.setState({ shared, wirelog });
          onFound(json.result.Metadata, log, extra);
          return;
        }

        onFound(json.result.Metadata, log, extra);
        StorageCache.putInCache(ck, json.result.Metadata, 60);
      })
      .catch(e => {
        this.errorOut(`Error Parsing Responce: ${e.message}`);
      });
  }

  // Passed to subclasses for them to call to add wirelogs
  pushWirelog(e) {
    const wirelog = _.clone(this.state.wirelog);
    wirelog.unshift(e);
    this.setState({ wirelog });
  }

  // Called when updating connection info
  updateConnection(connection, args) {
    // Clear local memory
    const sck = `${this.state.shared.source}-${this.state.shared.connection.id}-search-history`;
    const ock = `${this.state.shared.source}-${this.state.shared.connection.id}-object-history`;
    const mck = `${this.state.shared.source}-${this.state.shared.connection.id}-metadata`;
    StorageCache.remove(sck);
    StorageCache.remove(ock);
    StorageCache.remove(mck);

    // Reset shared variables
    const shared = {
      connection,
      metadata: Server.emptyMetadata,
      resource: {},
      class: {},
      fields: [],
      source: this.props.connection.config,
      args,
    };
    const wirelog = [];
    this.setState({ shared, errorOut: '', validMetadata: true, loading: true }, () => {
      // Retrive new Metadata
      this.getMetadata((m, log, extra) => {
        // If successful this function is called with the metadata
        console.log('Setting ', m);
        shared.metadata = m;
        wirelog.unshift({ tag: 'Metadata', log, extra });
        this.setState({ shared, wirelog, loading: false });
      });
    });
  }

  // Validates that all required peices of the Metadata available
  validateMetadata(data) {
    if (data.System && !data.System.SystemID) {
      data.System.SystemID = 'Missing ID'; // eslint-disable-line no-param-reassign
    }
    if (!data.System) {
      return this.failMetadata('System is not found');
    } else if (!data.System['METADATA-RESOURCE']) {
      return this.failMetadata('METADATA-RESOURCE is not found');
    } else if (!data.System['METADATA-RESOURCE'].Resource) {
      return this.failMetadata('Resource is not found');
    } else if (!data.System['METADATA-RESOURCE'].Resource.length) {
      return this.failMetadata('Resouce is empty');
    }
    return this.successMetadata();
  }


  // Series of shorthand functions that format a setState
  successMetadata() {
    this.setState({ validMetadata: true, metadataIssue: '' });
    return true;
  }


  failMetadata(issue) {
    this.setState({ validMetadata: false, metadataIssue: issue });
    return false;
  }

  errorOut(errorOut) {
    this.setState({ errorOut });
  }

  renderLoading() {
    return (
      <div>
        {this.renderError()}
        <div className={'loading-wrap db'}>
          <div className="loading" id={`${this.props.idprefix}-loading`}>LOADING METADATA</div>
        </div>
      </div>
    );
  }

  renderError() {
    return (
      <div
        className={`bg-dark-red white br1 pa4 w-100 tc ${this.state.errorOut.length === 0 ? 'dn' : 'db'}`}
        id={`${this.props.idprefix}-error`}
      >
        {this.state.errorOut}
      </div>
    );
  }

  renderInvalid() {
    const pages =
      [
        <Explore
          shared={this.state.shared}
          idprefix={`${this.props.idprefix}-Explore`}
        />,
        <Wireloger
          wirelog={this.state.wirelog}
          idprefix={`${this.props.idprefix}-Wirelog`}
        />,
      ];
    return (
      <div>
        <div
          className={'bg-dark-red white br1 pa4 w-100 tc'}
          id={`${this.props.idprefix}-error`}
        >
          Invalid Metadata Form: {this.state.metadataIssue}
        </div>
        <TabSection
          className="customTabElementB"
          components={[
            {
              id: 'Explore',
              page: pages[0],
              idprefix: 'Explore',
            },
            {
              id: 'Wirelog',
              page: pages[1],
              idprefix: 'Wirelog',
            },
          ]}
          initID={this.props.init ? this.props.init.tab : null}
          tag={this.props.idprefix}
        />
      </div>
    );
  }

  renderStandard() {
    const pages =
      [
        <Metadata
          shared={this.state.shared}
          onRowsSelected={this.onMetadataSelected}
          onRowsDeselected={this.onMetadataDeselected}
          onRowsCleared={this.onMetadataCleared}
          onClassSelected={this.onClassSelected}
          idprefix={`${this.props.idprefix}-Metadata`}
        />,
        <Search
          shared={this.state.shared}
          addTab={this.addTab}
          pushWirelog={this.pushWirelog}
          init={this.props.init && this.props.init.tab === 'Search' ? this.props.init : null}
          idprefix={`${this.props.idprefix}-Search`}
        />,
        <Objects
          shared={this.state.shared}
          addTab={this.addTab}
          pushWirelog={this.pushWirelog}
          init={this.props.init && this.props.init.tab === 'Objects' ? this.props.init : null}
          idprefix={`${this.props.idprefix}-Objects`}
        />,
        <Explore
          shared={this.state.shared}
          idprefix={`${this.props.idprefix}-Explore`}
        />,
        <Wireloger
          wirelog={this.state.wirelog}
          idprefix={`${this.props.idprefix}-Wirelog`}
        />,
      ];
    return (
      <div>
        {this.renderError()}
        <div>
          <TabSection
            className="customTabElementB"
            components={[
              {
                id: 'Metadata',
                page: pages[0],
                idprefix: 'Metadata',
              },
              {
                id: 'Search',
                page: pages[1],
                idprefix: 'Search',
              },
              {
                id: 'Objects',
                page: pages[2],
                idprefix: 'Objects',
              },
              {
                id: 'Explore',
                page: pages[3],
                idprefix: 'Explore',
              },
              {
                id: 'Wirelog',
                page: pages[4],
                idprefix: 'Wirelog',
              },
            ]}
            removeOffset={pages.length}
            initID={this.props.init ? this.props.init.tab : null}
            tag={this.props.idprefix}
          />
        </div>
      </div>
    );
  }

  renderSelect() {
    if (this.state.loading) {
      return this.renderLoading();
    }
    if (!this.state.validMetadata) {
      return this.renderInvalid();
    }
    return this.renderStandard();
  }

  render() {
    return (
      <div>
        <div className="fr">
          <div className="customHoverSection" id={`${this.props.idprefix}-config-hover`}>
            <button className="fr ma-3 customHoverBar" id={`${this.props.idprefix}-config`}> Connection Config </button>
            <div className="customHoverBody">
              <ConnectionForm
                updateConnection={this.updateConnection}
                args={this.state.shared.args}
                connection={this.state.shared.connection}
                location={this.props.location}
                idprefix={`${this.props.idprefix}-config`}
              />
            </div>
          </div>
        </div>
        {this.renderSelect()}
      </div>
    );
  }
}

export default withRouter(Server);
