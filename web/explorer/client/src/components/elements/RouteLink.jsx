import React from 'react';

const Base64 = require('js-base64').Base64;

class RouteLink extends React.Component {

  static propTypes = {
    style: React.PropTypes.any,

    connection: React.PropTypes.any,
    init: React.PropTypes.any,
    type: React.PropTypes.any,
    args: React.PropTypes.any,

    idprefix: React.PropTypes.any,
  }

  /*
  init: {
    [tab: <tab inside of id to auto select>
    [query: <query to run once open>]]
  }

  tab = Metadata, Search, Object, Explore, or Wirelog
  */

  constructor(props) {
    super(props);
    this.state = {
      ref: null,
    };

    this.getURL = this.getURL.bind(this);
    this.copyToClipboard = this.copyToClipboard.bind(this);
  }

  getURL() {
    const base = `${location.protocol}//${window.location.host}/#/?`;
    let rtn = null;
    switch (this.props.type) {
      case 'basic': {
        // Drops in with a source open (only requires connection info)
        rtn = `${base}${this.encodeConnection()}&${this.encodeID()}`;
        break;
      }
      case 'basicTab': {
        // Drops in with a source open to a sepcific tab
        rtn = `${base}${this.encodeConnection()}&${this.encodeTab()}`;
        break;
      }
      case 'full': {
        // Drops in with a source open to a sepcific tab and fills in the query
        rtn = `${base}${this.encodeConnection()}&${this.encodeQuery()}`;
        break;
      }
      case 'fullAuto': {
        // Drops in with a source open to a sepcific tab and launches a query imidietly
        rtn = `${base}${this.encodeConnection()}&${this.encodeQuery({ launch: 'auto' })}`;
        break;
      }
      default: {
        rtn = base;
      }
    }
    return rtn;
  }

  encodeConnection() {
    return `s=${Base64.encode(JSON.stringify(this.props.connection))}`;
  }

  encodeID(extra) {
    return `i=${Base64.encode(JSON.stringify(
      {
        ...extra,
        id: this.props.connection.id,
        args: this.props.args,
      }
    ))}`;
  }

  encodeTab(extra) {
    return this.encodeID({
      ...extra,
      tab: this.props.init.tab,
    });
  }

  encodeQuery(extra) {
    return this.encodeTab({
      ...extra,
      query: this.props.init.query,
    });
  }


  copyToClipboard() {
    console.log('[Copied To Clipboard]', this.state.ref.value);
    this.state.ref.select();
    document.execCommand('Copy');
  }

  render() {
    return (
      <div className="customCombo" style={this.props.style}>
        <button
          className="customComboButton"
          onClick={() => this.copyToClipboard()}
          id={`${this.props.idprefix}-cpclipboard`}
        >Link</button>
        <input
          className="customComboInput"
          value={this.getURL()}
          ref={(ref) => (this.state.ref = ref)}
          id={`${this.props.idprefix}-value`}
        />
      </div>
    );
  }

}

export default RouteLink;
