import React from 'react';

export default class NotFound extends React.Component {

  static propTypes = {
    children: React.PropTypes.any,
  }

  render() {
    return (
      <h1>Page not found</h1>
    );
  }
}
