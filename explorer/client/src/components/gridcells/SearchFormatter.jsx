import React from 'react';

export default class SearchFormatter extends React.Component {
  static propTypes = {
    value: React.PropTypes.any,
    select: React.PropTypes.Func,
  };

  constructor(props) {
    super(props);
    this.state = {
      element: null,
    };
    this.copyToClipboard = this.copyToClipboard.bind(this);
  }

  copyToClipboard() {
    this.state.element.select();
    document.execCommand('Copy');
  }

  render() {
    const value = this.props.value;
    const element = (
      <div className="flex customSearchFormater">
        <input type="text" value={value} ref={(ref) => (this.state.element = ref)} />
        <button className="fr" onClick={() => this.copyToClipboard()} title="Copy to Clipboard" >
          {'c'}
        </button>
      </div>
    );
    this.state.element = element;
    return element;
  }
}
