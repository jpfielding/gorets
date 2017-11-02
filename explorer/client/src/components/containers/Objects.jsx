import React from 'react';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ObjectsService from 'services/ObjectsService';
import StorageCache from 'util/StorageCache';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import SearchHistory from 'components/containers/SearchHistory';
import _ from 'underscore';

class Objects extends React.Component {

  static propTypes = {
    shared: React.PropTypes.any,
    addTab: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);
    const objectsForm = createValue({
      value: {},
      onChange: this.searchInputsChange.bind(this),
    });
    this.state = {
      objectsParams: ObjectsService.params,
      objectsForm,
      objectsHistory: [],
      objects: {},
      searching: false,
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
    // TODO set inital state
    const ock = `${this.props.shared.connection.id}-object-history`;
    let objectsHistory = StorageCache.getFromCache(ock) || [];
    if (objectsHistory && objectsHistory.length > 0) {
      objectsHistory = objectsHistory.filter((i) => (i.ids));
    }
    this.setState({
      objectsHistory,
    });
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.shared.class['METADATA-TABLE']) {
      const resource = nextProps.shared.resource.ResourceID;
      const ids = this.state.objectsForm.value.ids;
      const objectsForm = createValue({
        value: { resource, ids },
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
      this.state.errorOut = `No Object Types found for ${this.state.objectsForm.value.resource}`;
      return [];
    }
    this.state.errorOut = '';
    return r['METADATA-OBJECT']['Object'].map(o => o.ObjectType) || [];
  }

  getObjectsByType(type) {
    const id = this.props.shared.connection.id;
    const { resource, ids } = this.state.objectsForm.value;
    const name = this.state.objectsHistoryName;
    this.setState({ objectsParams: { resource, type, ids, id, name } }, this.getObjects);
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
    this.setState({ searching: true });
    const resource = this.state.objectsParams.resource;
    const type = this.state.objectsParams.type;
    const id = this.state.objectsParams.id;
    const ids = this.state.objectsParams.ids;

    const objectsForm = createValue({
      value: { resource, type, id, ids },
      onChange: this.searchInputsChange.bind(this),
    });

    this.setState({ objectsForm });

    const ock = `${this.props.shared.connection.id}-object-history`;
    let objectsHistory = StorageCache.getFromCache(ock) || [];
    const objectsParams = this.state.objectsParams;
    if (!objectsParams.ids) {
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
        if (!some(objectsHistory, objectsParams)) {
          objectsHistory.unshift(objectsParams);
          StorageCache.putInCache(ock, objectsHistory, 720);
        }
        console.log(json);
        if (objectsHistory && objectsHistory.length > 0) {
          objectsHistory = objectsHistory.filter((i) => (i.ids));
        }
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
    const sck = `${this.props.shared.connection.id}-object-history`;
    let objectsHistory = StorageCache.getFromCache(sck) || [];
    if (some(objectsHistory, params)) {
      objectsHistory.splice(objectsHistory.findIndex(i => _.isEqual(i, params)), 1);
      StorageCache.putInCache(sck, objectsHistory, 720);
    }
    if (objectsHistory && objectsHistory.length > 0) {
      objectsHistory = objectsHistory.filter((i) => (i.ids));
    }
    this.setState({
      objectsHistory,
    });
  }

  searchInputsChange(objectsForm) {
    this.setState({ objectsForm });
  }

  renderObjectInfo(obj) {
    const rows = Object.keys(obj)
      .filter(k => k !== 'Blob')
      .filter(k => obj[k] !== null)
      .filter(k => k !== 'Location' && obj[k] !== 0)
      .map(k => (
        <tr><td>{k}</td><td>{obj[k]}</td></tr>
      ));
    return (
      <table>
        {rows}
      </table>
    );
  }

  renderPicture(obj) {
    if (obj.RetsError) {
      return <div className="b mv3">An error occured, ${obj.RetsError}</div>;
    }
    if (!obj.ContentType.startsWith('image/')) {
      return null;
    }
    if (obj.location) {
      return (
        <li className="pa0 ma0 no-list-style">
          {this.renderObjectInfo(obj)}
          <img src={`data:image/png;base64,${obj.location}`} alt="pic" />
        </li>
      );
    }
    if (obj.Blob) {
      return (
        <li className="pa0 ma0 no-list-style">
          {this.renderObjectInfo(obj)}
          <img src={`data:image/png;base64,${obj.Blob}`} alt="pic" />
        </li>
      );
    }
    return null;
  }

  render() {
    const { objects } = this.state;
    const hasResult = (objects.result && objects.result['Objects'].length > 0);
    return (
      <div className="flex">
        <div className="fl w-100 w-20-ns pa3 pr0">
          <div className="customResultsCompactSet">
            <div className="customResultsTitle">
              <div className="b nonclickable">Current Object Params</div>
            </div>
            <div className="customResultsBody">
              <SearchHistory params={this.state.objectsParams} preventClose />
            </div>
          </div>
          <div className="customResultsCompactSet">
            <div className="customResultsTitle">
              <div className="b nonclickable">Objects History</div>
            </div>
            <div className="customResultsBody">
              <ul className="pa0 ma0 no-list-style customListSpacing">
                {this.state.objectsHistory.map(params =>
                  <li>
                    <SearchHistory
                      params={params}
                      removeHistoryElement={() => this.removeHistoryElement(params)}
                      clickHandle={() => {
                        // TODO convert getObjects to accept params directly
                        this.setState({
                          objectsParams: params,
                        }, this.getObjects);
                      }}
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
              />
            </div>
            <div className="customResultsBody">
              <Fieldset formValue={this.state.objectsForm}>
                <Field select="resource" label="Resource">
                  <Input className="w-30 pa1 b--none outline-transparent" />
                </Field>
                <Field select="ids" label="IDs">
                  <Input className="w-30 pa1 b--none outline-transparent" />
                </Field>
              </Fieldset>
            </div>
            <div className="customResultsFoot">
              <div className="customButtonSet">
                {this.getObjectTypes().map(type =>
                  <button
                    onClick={() => this.getObjectsByType(type)}
                    disabled={this.state.searching}
                  >
                    {type}
                  </button>
                )}
              </div>
              <div className={`bg-dark-red white br1 pa2 ${this.state.errorOut.length === 0 ? 'dn' : 'dib'}`}>
                {this.state.errorOut}
              </div>
            </div>
          </div>
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <div className="customCombo fr">
                <button className="customComboButton" onClick={this.createNewTab}>New Tab</button>
                <input
                  className="customComboInput"
                  placeholder={`O${this.state.resultCount}`}
                  onChange={(e) => this.bindTabNameChange(e.target.value)}
                />
              </div>
              <div className="customTitle">
                  Results:
              </div>
            </div>
            <div className="customResultsBody">
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
            </div>
          </div>
        </div>
      </div>
    );
  }

}

export default withRouter(Objects);
