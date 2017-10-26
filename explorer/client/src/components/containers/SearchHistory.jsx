import React from 'react';

export default class SearchHistory extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
    onClick: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);
    this.state = {
      orderList: ['name', 'id', 'resource', 'class', 'select', 'count-type', 'query'],
    };
    this.sortKeys = this.sortKeys.bind(this);
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

  render() {
    if (this.props.params == null || this.props.params.id == null || this.props.params.resource == null) {
      return <div />;
    }
    const body = Object.keys(this.props.params).sort(this.sortKeys).map((key) => {
      if (key === 'name') {
        return (
          <li>
            {this.props.params[key]}
          </li>
        );
      }
      return (
        <li>
          {key}{' - '}<span className="moon-gray">{this.props.params[key]}</span>
        </li>
      );
    });
    return (
      <ul
        className="pa0 ma0 w-100 no-list-style clickable"
        onClick={this.props.onClick}
        style={{
          padding: '8px',
          margin: '8px 0px',
          wordBreak: 'break-all',
        }}
      >
        {body}
      </ul>
    );
  }
}
