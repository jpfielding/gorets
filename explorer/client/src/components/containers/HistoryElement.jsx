import React from 'react';

export default class HistoryElement extends React.Component {

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
        className="pa0 ma0 w-100 no-list-style"
        onClick={this.props.onClick}
        style={{
          backgroundColor: 'white',
          padding: '5px',
          margin: '8px 0px',
          wordBreak: 'break-all',
        }}
      >
        <li><h1 className="mb0 mt0 f5 dib">ID:</h1>{this.props.params.id}</li>
        <li><h1 className="mb0 mt3 f5 dib">Resource:</h1>{this.props.params.resource}</li>
        <li><h1 className="mb0 mt3 f5 dib">Class:</h1>{this.props.params.class}</li>
        <li><h1 className="mb0 mt3 f5 dib">Query:</h1>{this.props.params.query}</li>
        <li><h1 className="mb0 mt3 f5 dib">Select:</h1>{this.props.params.select}</li>
      </ul>
    );
  }
}
