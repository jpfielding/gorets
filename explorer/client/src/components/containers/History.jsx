import React from 'react';

export default class ContentHistory extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
    preventClose: React.PropTypes.any,

    clickHandle: React.PropTypes.Func,
    removeHistoryElement: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);
    this.state = {
      orderList: ['name', 'id', 'resource', 'class', 'select', 'counttype', 'query'],
      exclude: ['connection', 'name', 'counttype', 'type'],
      countTypes: ['Submit', 'with Count', 'only Count'],
    };
    this.sortKeys = this.sortKeys.bind(this);
    this.renderClose = this.renderClose.bind(this);
    this.renderName = this.renderName.bind(this);
    this.renderType = this.renderType.bind(this);
  }

  sortKeys(a, b) {
    const ai = this.state.orderList.indexOf(a);
    const bi = this.state.orderList.indexOf(b);
    if (bi === -1) {
      return -1;
    }
    if (ai === -1) {
      return 1;
    }
    return ai - bi;
  }

  renderClose() {
    if (this.props.preventClose) {
      return null;
    }
    return (
      <button
        className="customCloseButton"
        onClick={this.props.removeHistoryElement}
      >
        X
      </button>
    );
  }

  renderName() {
    if (this.props.params.name) {
      return (
        <div className="customResultsSimpleTitle">
          {this.props.params['name']}
        </div>
      );
    }
    return (
      <div className="customResultsSimpleTitle">
        N/A
      </div>
    );
  }

  renderType() {
    if (this.props.params['counttype'] >= 0) {
      return (
        <div className="customResultsButtonTitle fr">
          {this.state.countTypes[this.props.params['counttype']]}
        </div>
      );
    } else if (this.props.params.type) {
      return (
        <div className="customResultsButtonTitle fr">
          {this.props.params.type}
        </div>
      );
    }
    return null;
  }

  render() {
    if (this.props.params == null || !this.props.params['connection']) {
      return <div />;
    }
    const body = Object.keys(this.props.params).sort(this.sortKeys).map((key) => {
      if (this.state.exclude.indexOf(key) !== -1) {
        return null;
      }
      return (
        <div>
          {key}{' - '}<span className="moon-gray">{this.props.params[key]}</span>
        </div>
      );
    });
    return (
      <div className="customCloseElement clickable">
        <div
          className="customResultsCompactSet"
          onClick={this.props.clickHandle}
        >
          <div className="customResultsTitle">
            {this.renderName()}
            {this.renderType()}
          </div>
          <div className="customResultsBody">
            {body}
          </div>
        </div>
        {this.renderClose()}
      </div>
    );
  }
}
