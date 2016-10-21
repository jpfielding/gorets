import React from 'react';
import { withRouter } from 'react-router';
import MetadataService from 'services/MetadataService';

class SearchMetadata extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      searchParams: {
        id: null,
        resource: null,
        keyFieldValue: null,
        types: null,
      },
      objectMetadata: {},
    };
  }

  componentWillMount() {
    this.setState({
      searchParams: this.props.location.query,
    });
    this.getObjectMetadata = this.getObjectMetadata.bind(this);
  }

  getObjectMetadata(type) {
    console.log(type);
    const { searchParams } = this.state;
    MetadataService
      .getObjectMetadata({
        id: searchParams.id,
        resource: searchParams.resource,
        type,
        objectid: `${searchParams.keyFieldValue}:*`,
      })
      .then((res) => res.json())
      .then(json => {
        console.log(json);
        this.setState({
          objectMetadata: json,
        });
      });
  }

  render() {
    const { searchParams } = this.state;
    return (
      <div>
        <div>
          <span className="b">Connection: </span>
          {searchParams.id}
        </div>
        <div>
          <span className="b">Resource: </span>
          {searchParams.resource}
        </div>
        <div>
          <span className="b">KeyFieldValue: </span>
          {searchParams.keyFieldValue}
        </div>
        <div>
          <span className="b">Available Types: </span>
          {searchParams.types.split(',').map(type =>
            <button className="link" onClick={() => this.getObjectMetadata(type)}>
              {type}
            </button>
          )}
        </div>
        <pre className="code f6">
          {JSON.stringify(this.state, null, '  ')}
        </pre>
      </div>
    );
  }

}

export default withRouter(SearchMetadata);
