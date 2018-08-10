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
      location: props.shared.metadata,
      treeSearch: '',
    };

    this.setSearch = this.setSearch.bind(this);

    this.bindSeachChange = this.bindSeachChange.bind(this);
    this.currentLocation = this.currentLocation.bind(this);
  }

  componentWillReceiveProps(newProps) {
    this.setState({ location: newProps.shared.metadata, treeSearch: '' });
  }

  setSearch(e) {
    const treeSearch = `${this.currentLocation()}/${e}`;
    this.bindSeachChange(treeSearch);
  }

  bindSeachChange(treeSearch) {
    const keyList = treeSearch.split('/');
    if (keyList.length > 0 && keyList[0] === '') {
      keyList.splice(0, 1);
    }
    if (keyList.length === 0) {
      this.setState({ treeSearch, location: this.props.shared.metadata });
      return;
    }
    let location = this.props.shared.metadata;
    keyList.forEach((e) => {
      if (location[e]) location = location[e];
    });
    this.setState({ treeSearch, location });
  }

  currentLocation() {
    const keyList = this.state.treeSearch.split('/');
    let rtn = '';
    if (keyList.length > 0 && keyList[0] === '') {
      keyList.splice(0, 1);
    }
    let location = this.props.shared.metadata;
    keyList.forEach((e) => {
      if (location[e]) {
        rtn += `/${e}`;
        location = location[e];
      }
    });
    return rtn;
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
      <ExploreObject key={key} k={key} value={value} setSearch={this.setSearch} />
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
            <input
              placeholder="/System/... "
              value={this.state.treeSearch}
              onChange={(e) => this.bindSeachChange(e.target.value)}
              className="treeSeachBar"
            />
            <button
              onClick={() => {
                const splt = this.currentLocation().split('/');
                if (splt.lenght !== 0) {
                  this.bindSeachChange(splt.slice(0, splt.length - 1).join('/'));
                }
              }}
              className="customButton"
            >
              {'< Back'}
            </button>
            <div className="treeview" >
              {
                Object.keys(this.state.location)
                  .map((e) => this.processKeys(e, this.state.location))
              }
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(Explore);
