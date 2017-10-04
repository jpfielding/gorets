import React from 'react';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ObjectsService from 'services/ObjectsService';
import StorageCache from 'util/StorageCache';
import { Fieldset, Field, createValue, Input } from 'react-forms';

class Objects extends React.Component {

  static propTypes = {
    connection: React.PropTypes.any,
    metadata: React.PropTypes.any,
    router: React.PropTypes.any,
    location: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    const objectsForm = createValue({
      value: {
        resource: 'Property',
        type: 'Photo',
      },
      onChange: this.searchInputsChange.bind(this),
    });
    this.state = {
      objectsParams: ObjectsService.params,
      objectsForm,
      objectsHistory: [],
      objects: {},
    };
    this.getObjects = this.getObjects.bind(this);
  }

  // componentWillMount() {
  //   this.setState({
  //     getParams: this.props.location.query,
  //   });
  // }

  getObjects() {
    const ock = `${this.props.connection.id}-search-history`;
    const objectsHistory = StorageCache.getFromCache(ock) || [];
    const { objectsParams } = this.state;
    if (!objectsParams.ids) {
      return;
    }
    const contentId = objectsParams.ids.split(',').map(
        // avoiding other lint issues
        i => [i, '*'].join(':')
    ).join(',');

    ObjectsService
      .getObjects({
        id: this.props.connection.id,
        resource: objectsParams.resource,
        type: objectsParams.type,
        objectid: contentId,
      })
      .then((res) => res.json())
      .then(json => {
        if (!some(objectsHistory, objectsParams)) {
          objectsHistory.unshift(objectsParams);
          StorageCache.putInCache(ock, objectsHistory, 720);
        }
        console.log(json);
        this.setState({
          objects: json,
        });
      });
  }

  getObjectTypes() {
    if (!this.state.searchParams) {
      return [];
    }
    const r = this.getResource();
    if (r == null) {
      return [];
    }
    return r['METADATA-OBJECT']['Object'].map(o => o.ObjectType) || [];
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

  getResource() {
    if (!this.state.searchParams) {
      return [];
    }
    const rs = this.props.metadata.System['METADATA-RESOURCE'].Resource.filter(
      r => (r.ResourceID === this.state.searchParams.resource)
    );
    if (rs.length === 0) {
      return null;
    }
    return rs[0];
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
      <div className="pa2">
        <div className="fl h-100-ns w-100 w-20-ns pa3 overflow-x-scroll nowrap">
          <div className="b">Current Object Params</div>
          <pre className="f6 code">{JSON.stringify(this.state.objectsParams, null, '  ')}</pre>
          <div className="b">Objects History</div>
          <ul className="pa0 ma0 no-list-style">
            {this.state.objectsHistory.map(params =>
              <li>
                <pre
                  className="f6 code clickable"
                  onClick={() => {
                    // TODO convert getObjects to accept params directly
                    this.setState({
                      objectsParams: params,
                    });
                    this.getObjects();
                  }}
                >
                  { JSON.stringify(params, null, '  ') }
                </pre>
              </li>
              )}
          </ul>
        </div>
        <div className="fl h-100 min-vh-100 w-100 w-80-ns pa3 bl-ns">
          <Fieldset formValue={this.state.objectsForm}>
            <Field select="resource" label="Resource">
              <Input className="w-30" />
            </Field>
            <Field select="ids" label="IDs">
              <Input className="w-30" />
            </Field>
          </Fieldset>
        </div>
        <div>
          {this.getObjectTypes().map(type =>
            <button onClick={() => this.getObjects(type)}>{type}</button>
          )}
        </div>
        <div>
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
        <pre className="code f6">
          {JSON.stringify(this.state, null, '  ')}
        </pre>
      </div>
    );
  }

}

export default withRouter(Objects);
