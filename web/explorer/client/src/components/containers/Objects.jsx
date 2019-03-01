import React from 'react';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ObjectsService from 'services/ObjectsService';
import StorageCache from 'util/StorageCache';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import ContentHistory from 'components/containers/History';
import RouteLink from 'components/elements/RouteLink';
import _ from 'underscore';

const Base64 = require('js-base64').Base64;

class Objects extends React.Component {

  static propTypes = {
    shared: React.PropTypes.any,
    init: React.PropTypes.any,
    addTab: React.PropTypes.Func,
    pushWirelog: React.PropTypes.Func,
    idprefix: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    const objectsForm = createValue({
      value: (this.props.init && this.props.init.query ? this.props.init.query : ObjectsService.params),
      onChange: this.searchInputsChange.bind(this),
    });
    this.state = {
      objectsParams: (this.props.init && this.props.init.query ? this.props.init.query : ObjectsService.params),
      objectsForm,
      objectsHistory: [],
      objects: {},
      searching: false,
      infoOut: '',
      errorOut: '',
      resultCount: 1,
      tabName: '',
      objectsHistoryName: '',
    };
    this.getObjects = this.getObjects.bind(this);
    this.getObjectsByType = this.getObjectsByType.bind(this);
    this.removeHistoryElement = this.removeHistoryElement.bind(this);

    this.bindTabNameChange = this.bindTabNameChange.bind(this);
    this.bindQueryNameChange = this.bindQueryNameChange.bind(this);

    this.createNewTab = this.createNewTab.bind(this);
  }

  componentWillMount() {
    const ock = `${this.props.shared.source}-${this.props.shared.connection.id}-object-history`;
    const objectsHistory = StorageCache.getFromCache(ock) || [];
    this.setState({
      objectsHistory,
    });
  }

  componentDidMount() {
    if (this.props.init && this.props.init.launch === 'auto') {
      console.log('Auto Searching');
      this.getObjects();
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.shared.class['METADATA-TABLE']) {
      const resource = nextProps.shared.resource.ResourceID;
      const ids = this.state.objectsForm.value.ids;
      const objectsForm = createValue({
        value: { resource, ids, ...ObjectsService.params },
        onChange: this.searchInputsChange.bind(this),
      });
      this.setState({ objectsForm });
    }
  }

  getResource() {
    if (!this.state.objectsForm) {
      return [];
    }
    const rs = this.props.shared.metadata.System['METADATA-RESOURCE'].Resource.filter(
      r => (r.ResourceID === this.state.objectsForm.value.resource)
    );
    if (rs.length === 0) {
      return null;
    }
    return rs[0];
  }

  getObjectTypes() {
    if (!this.state.objectsForm) {
      return [];
    }
    const r = this.getResource();
    if (r == null || !r['METADATA-OBJECT']['Object']) {
      this.state.infoOut = `No Object Types found for ${this.state.objectsForm.value.resource}`;
      return [];
    }
    this.state.infoOut = '';
    return r['METADATA-OBJECT']['Object'].map(o => o.ObjectType) || [];
  }

  getObjectsByType(type) {
    const connection = this.props.shared.connection.id;
    const { resource, ids, location } = this.state.objectsForm.value;
    const name = this.state.objectsHistoryName;
    this.setState({ objectsParams: { resource, type, ids, connection, name, location } }, this.getObjects);
  }

  getKeyFieldColumn() {
    const { searchResultColumns } = this.state;
    const keyField = this.getResource().KeyField;
    const keyFieldCols = searchResultColumns.filter(c => (c.name === keyField));
    if (keyFieldCols.length === 0) {
      return null;
    }
    return keyFieldCols[0];
  }

  getObjects() {
    this.setState({ searching: true, errorOut: '', objects: [] });
    const resource = this.state.objectsParams.resource;
    const type = this.state.objectsParams.type;
    const connection = this.state.objectsParams.connection;
    const ids = this.state.objectsParams.ids;
    const location = this.state.objectsParams.location;

    const objectsForm = createValue({
      value: { resource, type, connection, ids, location },
      onChange: this.searchInputsChange.bind(this),
    });

    this.setState({ objectsForm });

    const ock = `${this.props.shared.source}-${this.props.shared.connection.id}-object-history`;
    const objectsHistory = StorageCache.getFromCache(ock) || [];
    const objectsParams = this.state.objectsParams;
    if (!objectsParams.ids) {
      this.setState({ searching: false, errorOut: 'No object IDs found' });
      return;
    }

    objectsParams.ids = objectsParams.ids.split(',').map(
      (i) => {
        if (i.indexOf(':') > -1) {
          return i;
        }
        return [i, '*'].join(':');
      }
    ).join(',');

    ObjectsService
      .getObjects(this.props.shared.connection, objectsParams)
      .then((res) => res.json())
      .then((json) => {
        console.log('Object Response', json);
        if (json.error && json.error !== null) {
          this.props.pushWirelog({ tag: 'Objects', log: json.error.message, extra: { type: 'Error' } });
          this.setState({ searching: false, errorOut: json.error.message });
          return;
        }
        let log = Base64.decode(json.result.wirelog);
        if (log === undefined) log = 'Unable to parse wirelog';
        this.props.pushWirelog({ tag: 'Objects', log });
        objectsParams.submited = true;

        if (!some(objectsHistory, objectsParams)) {
          objectsHistory.unshift(objectsParams);
          StorageCache.putInCache(ock, objectsHistory, 720);
        }
        console.log(json);
        this.setState({
          objects: json,
          searching: false,
          objectsHistory,
        });
      });
  }

  createNewTab() {
    let tabName = this.state.tabName;
    if (tabName === '') {
      tabName = `O${this.state.resultCount}`;
      const resultCount = this.state.resultCount + 1;
      this.setState({ resultCount });
    }
    console.log(`[OBJECT] Creating new tab of name ${tabName}`);
    const { objects } = this.state;
    const hasResult = (objects.result && objects.result['Objects'].length > 0);
    this.props.addTab(tabName, (
      <ul>
        {hasResult
          ? (
            objects.result['Objects'].map(obj =>
              this.renderPicture(obj)
            )
          )
          : null
        }
      </ul>
    ));
  }

  bindTabNameChange(tabName) {
    this.setState({ tabName });
  }

  bindQueryNameChange(objectsHistoryName) {
    this.setState({ objectsHistoryName });
  }

  removeHistoryElement(params) {
    console.log('[OBJECT] Removing history element');
    const sck = `${this.props.shared.source}-${this.props.shared.connection.id}-object-history`;
    const objectsHistory = StorageCache.getFromCache(sck) || [];
    if (some(objectsHistory, params)) {
      objectsHistory.splice(objectsHistory.findIndex(i => _.isEqual(i, params)), 1);
      StorageCache.putInCache(sck, objectsHistory, 720);
    }
    this.setState({
      objectsHistory,
    });
  }

  search(params) {
    this.setState({
      objectsParams: params,
    }, this.getObjects);
  }

  searchInputsChange(objectsForm) {
    let location = Number.parseInt(objectsForm.value.location, 10);
    if (isNaN(location)) {
      location = 0;
    }

    this.setState({ objectsForm: createValue({
      value: { ...objectsForm.value, location },
      onChange: this.searchInputsChange.bind(this),
    }) });
  }

  renderObjectInfo(obj, i) {
    const rows = Object.keys(obj)
      .filter(k => k !== 'Blob')
      .filter(k => obj[k] !== null)
      .map(k => {
        if (k === 'Location') {
          return (
            <tr><td>{k}</td><td><a href={obj[k]}>{obj[k]}</a></td></tr>
          );
        }
        return (
          <tr><td>{k}</td><td>{obj[k]}</td></tr>
        );
      });
    return (
      <table id={`${this.props.idprefix}-result-${i}-table`}>
        {rows}
      </table>
    );
  }

  renderPicture(obj, i) {
    if (obj.RetsError) {
      return (<div><div
        className="b mv3 bg-dark-red white br1 pa2 dib" id={`${this.props.idprefix}-result-${i}-error`}
      >An error occured, {obj.RetsError}</div></div>);
    }
    if (obj.Blob) {
      if (obj.ContentType.startsWith('text')) {
        return (
          <li className="pa0 ma0 no-list-style">
            {this.renderObjectInfo(obj, i)}
            <div><div>{'Blob'}</div><pre>{atob(obj['Blob'])}</pre></div>
          </li>
        );
      }
      return (
        <li className="pa0 ma0 no-list-style">
          {this.renderObjectInfo(obj, i)}
          <img
            src={`data:image/png;base64,${obj.Blob}`}
            id={`${this.props.idprefix}-result-${i}-image`} alt="pic"
          />
        </li>
      );
    }
    return this.renderObjectInfo(obj, i);
  }

  render() {
    const { objects } = this.state;
    const hasResult = (objects.result && objects.result['Objects'].length > 0);
    return (
      <div
        className="flex"
        style={{
          maxWidth: '1500px',
          margin: 'auto',
        }}
      >
        <div className="fl w-100 w-20-ns pa3 pr0">
          <div className="customResultsCompactSet">
            <div className="customResultsTitle">
              <div className="b nonclickable">Current Object Params</div>
            </div>
            <div className="customResultsBody">
              <ContentHistory
                params={this.state.objectsParams} preventClose
                idprefix={`${this.props.idprefix}-history-current`}
              />
            </div>
          </div>
          <div className="customResultsCompactSet">
            <div className="customResultsTitle">
              <div className="b nonclickable">Objects History</div>
            </div>
            <div className="customResultsBody">
              <ul className="pa0 ma0 no-list-style customListSpacing">
                {this.state.objectsHistory.map((params, i) =>
                  <li>
                    <ContentHistory
                      params={params}
                      removeHistoryElement={() => this.removeHistoryElement(params)}
                      clickHandle={() => this.search(params)}
                      idprefix={`${this.props.idprefix}-history-${i}`}
                    />
                  </li>
                  )}
              </ul>
            </div>
          </div>
        </div>
        <div className="fl w-100 w-80-ns pa3">
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <div className="customTitle">
                  Query:
              </div>
              <input
                className="customInputBar pt1"
                placeholder="Query Name"
                onChange={(e) => this.bindQueryNameChange(e.target.value)}
                value={this.state.objectsHistoryName}
                id={`${this.props.idprefix}-query-name`}
              />
              <RouteLink
                type={'full'}
                connection={this.props.shared.connection}
                args={this.props.shared.args}
                init={{ tab: 'Objects', query: this.state.objectsForm.value }}
                style={{ float: 'right' }}
                idprefix={`${this.props.idprefix}-query-link`}
              />
            </div>
            <div className="customResultsBody">
              <Fieldset formValue={this.state.objectsForm}>
                <Field select="resource" label="Resource">
                  <Input
                    className="w-30 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-resource`}
                  />
                </Field>
                <Field select="ids" label="IDs">
                  <Input
                    className="w-30 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-ids`}
                  />
                </Field>
                <Field select="location" label="Location">
                  <Input
                    className="w-30 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-location`}
                  />
                </Field>
              </Fieldset>
            </div>
            <div className="customResultsFoot">
              <div className="customButtonSet">
                {this.getObjectTypes().map(type =>
                  <button
                    onClick={() => this.getObjectsByType(type)}
                    disabled={this.state.searching}
                    id={`${this.props.idprefix}-submit-${type}`}
                  >
                    {type}
                  </button>
                )}
              </div>
              <div className={`bg-dark-red white br1 pa2 ${this.state.infoOut.length === 0 ? 'dn' : 'dib'}`}>
                {this.state.infoOut}
              </div>
            </div>
          </div>
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <RouteLink
                type={'fullAuto'}
                connection={this.props.shared.connection}
                args={this.props.shared.args}
                init={{ tab: 'Objects', query: this.state.objectsForm.value }}
                style={{ float: 'right' }}
                idprefix={`${this.props.idprefix}-result-link`}
              />
              <div className="customTitle">
                  Results:
              </div>
            </div>
            <div className="customResultsBody">
              <div className={`bg-dark-red white br1 pa2 ${this.state.errorOut.length === 0 ? 'dn' : 'dib'}`}>
                {this.state.errorOut}
              </div>
              <ul>
                {hasResult
                  ? (
                    objects.result['Objects'].map((obj, i) =>
                      this.renderPicture(obj, i)
                    )
                  )
                  : null
                }
              </ul>
            </div>
          </div>
        </div>
      </div>
    );
  }

}

export default withRouter(Objects);
