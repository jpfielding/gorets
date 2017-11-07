import React from 'react';

export default class KeyFormatter extends React.Component {
  static propTypes = {
    value: React.PropTypes.any,
    extra: React.PropTypes.any,
  };

  render() {
    const value = this.props.value;
    if (value === this.props.extra) {
      return (
        <div>
          <div className="customResultsButtonTitle" style={{ display: 'inline-block', marginRight: '5px' }}>
            Key
          </div>
          {value}
        </div>
      );
    }
    return (
      <div>
        {value}
      </div>
    );
  }
}
