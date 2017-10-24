import React from 'react';

export default class SearchHistory extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
    onClick: React.PropTypes.Func,
  }

  // constructor(props) {
  //   super(props);
  // }

  render() {
    if (this.props.params == null || this.props.params.id == null || this.props.params.resource == null) {
      return <div />;
    }
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
        <li>
          {'ID - '}<span className="moon-gray">{this.props.params.id}</span>
        </li>
        <li>
          {'Resource - '}<span className="moon-gray">{this.props.params.resource}</span>
        </li>
        <li>
          {'Class - '}<span className="moon-gray">{this.props.params.class}</span>
        </li>
        <li>
          {'Query - '}<span className="moon-gray">{this.props.params.query}</span>
        </li>
        <li>
          {'Select - '}<span className="moon-gray">{this.props.params.select}</span>
        </li>
      </ul>
    );
  }
}
