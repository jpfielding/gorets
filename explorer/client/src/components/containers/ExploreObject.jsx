import React from 'react';
import { withRouter } from 'react-router';

class ExploreObject extends React.Component {

  static propTypes = {
    k: React.PropTypes.any,
    value: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      key: props.k,
      value: null,
    };
  }

  processKeys(key, context) {
    const value = context[key];
    if (value === null) return null;
    if (typeof value !== 'object') {
      return (
        <div className="leaf" key={key} >{key}<span>{value}</span></div>
      );
    }
    return (
      <ExploreObject key={key} k={key} value={value} />
    );
  }

  render() {
    return (
      <div
        key={this.state.key}
        className="treeview"
      >
        <button
          onClick={() => {
            if (this.state.value == null) {
              this.setState({ value: this.props.value });
            } else {
              this.setState({ value: null });
            }
          }}
        >
          { this.state.value == null ? '\u25BC' : '\u25B2' }
        </button>
        {this.state.key}
        <div className="branch">
          {
            this.state.value == null ? null :
            Object.keys(this.state.value)
              .map((e) => this.processKeys(e, this.state.value))
          }
        </div>
      </div>
    );
  }
}

export default withRouter(ExploreObject);
