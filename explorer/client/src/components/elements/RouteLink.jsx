import React from 'react';

const Base64 = require('js-base64').Base64;

class RouteLink extends React.Component {

  static propTypes = {
    style: React.PropTypes.any,

    connection: React.PropTypes.any,
    init: React.PropTypes.any,
    type: React.PropTypes.any,
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
        rtn = `${base}s=${Base64.encode(JSON.stringify(this.props.connection))}&`;
        rtn = `${rtn}i=${Base64.encode(`{"id":"${this.props.connection.id}"}`)}`;
        break;
      }
      case 'basicTab': {
        // Drops in with a source open to a sepcific tab
        rtn = `${base}s=${Base64.encode(JSON.stringify(this.props.connection))}&`;
        rtn = `${rtn}i=${Base64.encode(`{"id":"${this.props.connection.id}","tab":"${this.props.init.tab}"}`)}`;
        break;
      }
      case 'full': {
        // Drops in with a source open to a sepcific tab and fills in the query
        const id = this.props.connection.id;
        const tab = this.props.init.tab;
        const query = JSON.stringify(this.props.init.query);
        rtn = `${base}s=${Base64.encode(JSON.stringify(this.props.connection))}&`;
        rtn = `${rtn}i=${Base64.encode(`{"id":"${id}","tab":"${tab}","query":${query}}`)}`;
        break;
      }
      case 'fullAuto': {
        // Drops in with a source open to a sepcific tab and launches a query imidietly
        const id = this.props.connection.id;
        const tab = this.props.init.tab;
        const query = JSON.stringify(this.props.init.query);
        rtn = `${base}s=${Base64.encode(JSON.stringify(this.props.connection))}&`;
        rtn = `${rtn}i=${Base64.encode(`{"id":"${id}","tab":"${tab}","query":${query},"launch":"auto"}`)}`;
        break;
      }
      default: {
        rtn = base;
      }
    }
    return rtn;
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
        >Link</button>
        <input
          className="customComboInput"
          value={this.getURL()}
          ref={(ref) => (this.state.ref = ref)}
        />
      </div>
    );
  }

}

export default RouteLink;
