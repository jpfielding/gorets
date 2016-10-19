import React from 'react';
import MetadataService from 'services/MetadataService';

export default class Explorer extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      metadata: {
        System: {
          'METADATA-RESOURCE': {
            Resource: [],
          },
        },
      },
      selectedClass: null,
    };
  }

  componentDidMount() {
    if (this.props.params.connection) {
      this.getMetadata(this.props.params.connection);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.params !== this.props.params && nextProps.params.connection) {
      this.getMetadata(nextProps.params.connection);
    }
  }

  getMetadata(connectionId) {
    MetadataService
      .get(connectionId)
      .then(response => response.json())
      .then(json => {
        console.log(json);
        this.setState({
          metadata: json,
        });
      });
  }

  metadataClassClick(selectedClass) {
    this.setState({
      selectedClass,
    });
  }

  render() {
    return (
      <div>
        <div className="fl w-100 w-30-ns no-list-style pa3">
          <ul className="pa0 ma0">
            {this.state.metadata.System['METADATA-RESOURCE'].Resource.map((resource) =>
              <li className="mb2">
                <div className="b">{resource.ResourceID}</div>
                <ul className="pa0 pl3 mv1">
                  {resource['METADATA-CLASS'].Class.map((mClass) =>
                    <li
                      onClick={() => this.metadataClassClick(mClass)}
                      className="clickable metadata-class-name"
                    >
                      {mClass.ClassName}
                    </li>
                  )}
                </ul>
              </li>
            )}
          </ul>
        </div>
        <div className="fl h-100 min-vh-100 w-100 w-70-ns pa3 bl-ns">
          <div className="f4">Table Continer</div>
          <pre className="overflow-x-scroll">{JSON.stringify(this.state.selectedClass, null, '  ')}</pre>
        </div>
      </div>
    );
  }
}
