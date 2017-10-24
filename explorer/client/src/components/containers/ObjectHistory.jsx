import React from 'react';

export default class ObjectHistory extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
    onClick: React.PropTypes.Func,
  }

  // constructor(props) {
  //   super(props);
  // }

  render() {
    if (this.props.params == null) {
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
          {'IDS - '}<span className="moon-gray">{this.props.params.ids}</span>
        </li>
        <li>
          {'Location - '}<span className="moon-gray">{this.props.params.location}</span>
        </li>
      </ul>
    );
  }
}
