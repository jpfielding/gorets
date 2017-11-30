import React from 'react';
import { withRouter } from 'react-router';
import MetadataService from 'services/MetadataService';
import ExploreObject from 'components/containers/ExploreObject';

class Explore extends React.Component {

  static propTypes = {
    shared: {
      connection: React.PropTypes.any,
      metadata: React.PropTypes.any,
      resource: React.PropTypes.any,
      class: React.PropTypes.any,
      fields: React.PropTypes.any,
      rows: React.PropTypes.any,
    },
  }

  static defaultProps = {
    metadata: MetadataService.empty,
  }

  constructor(props) {
    super(props);
    this.state = {
    };
  }

  processKeys(key, context) {
    const value = context[key];
    if (value === null) return null;
    if (typeof value !== 'object') {
      return (
        <div key={key} className="leaf" >{key}<span>{value}</span></div>
      );
    }
    return (
      <ExploreObject key={key} k={key} value={value} />
    );
  }

  render() {
    return (
      <div
        className="flex pa3"
        style={{
          maxWidth: '1500px',
          margin: 'auto',
        }}
      >
        <div className="customResultsSet w-100">
          <div className="customResultsTitle">
            <div className="customTitle">
              Explore:
            </div>
          </div>
          <div className="customResultsBody" >
            <div className="treeview" >
              {
                Object.keys(this.props.shared.metadata)
                  .map((e) => this.processKeys(e, this.props.shared.metadata))
              }
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(Explore);
