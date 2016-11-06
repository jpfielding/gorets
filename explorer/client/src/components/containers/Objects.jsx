import React from 'react';
import { withRouter } from 'react-router';
import ObjectsService from 'services/ObjectsService';

class Objects extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      searchParams: {
        id: null,
        resource: null,
        ids: null,
        types: null,
      },
      objects: {},
    };
  }

  componentWillMount() {
    this.setState({
      searchParams: this.props.location.query,
    });
    this.getObjects = this.getObjects.bind(this);
  }

  getObjects(type) {
    const { searchParams } = this.state;
    const contentId = searchParams.ids.split(',').map(
        // avoiding other lint issues
        i => [i, '*'].join(':')
    ).join(',');

    ObjectsService
      .getObjects({
        id: searchParams.id,
        resource: searchParams.resource,
        type,
        objectid: contentId,
      })
      .then((res) => res.json())
      .then(json => {
        console.log(json);
        this.setState({
          objects: json,
        });
      });
  }

  renderObjectInfo(obj) {
    return (
      <table>
        <tr><td>Content ID</td><td>{obj.ContentID}</td></tr>
        <tr><td>Content Type</td><td>{obj.ContentType}</td></tr>
        <tr><td>Object ID</td><td>{obj.ObjectID}</td></tr>
        <tr><td>UID</td><td>{obj.UID}</td></tr>
        <tr><td>Description</td><td>{obj.Description}</td></tr>
        <tr><td>SubDescription</td><td>{obj.SubDescription}</td></tr>
        <tr><td>Location</td><td>{obj.Location}</td></tr>
        <tr><td>Preferred</td><td>{obj.Preferred}</td></tr>
        <tr><td>ObjectData</td><td>{obj.ObjectData}</td></tr>
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
    const { searchParams, objects } = this.state;
    const hasResult = (objects.result && objects.result['Objects'].length > 0);
    return (
      <div className="pa2">
        <div>
          <span className="b">Connection: </span>
          {searchParams.id}
        </div>
        <div>
          <span className="b">Resource: </span>
          {searchParams.resource}
        </div>
        <div>
          <span className="b">Object IDs: </span>
          {searchParams.ids}
        </div>
        <div>
          <span className="b">Available Types: </span>
          {searchParams.types.split(',').map(type =>
            <button className="link" onClick={() => this.getObjects(type)}>
              {type}
            </button>
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
