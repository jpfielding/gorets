import React from 'react';
import Papa from 'papaparse';

export default class DownloadLink extends React.Component {

  static propTypes = {
    data: React.PropTypes.any,
    headers: React.PropTypes.any,
    tag: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = this.genStateBlob(props.data, props.headers);
  }

  componentWillReceiveProps(newProps) {
    if (this.props.data !== newProps.data || this.props.headers !== newProps.headers) {
      this.setState(this.genStateBlob(newProps.data, newProps.headers));
    }
  }

  genStateBlob = (data, headers) => {
    const blob = this.genBlob(data, headers);
    return {
      url: URL.createObjectURL(blob),
    };
  }

  genBlob = (data, headers) => {
    const columns = headers.map((item) => item.key);
    const newData = data.map((dataObject) => columns.map((key) => dataObject[key]));
    const final = [
      columns,
      ...newData,
    ];
    return new Blob([Papa.unparse(final)], { type: 'text/csv' });
  }

  render() {
    if (this.state.url == null) {
      return (
        <div>
          {'No Download Available'}
        </div>
      );
    }
    return (
      <div>
        <a href={this.state.url} download={`${this.props.tag || 'data'}.csv`}>Download</a>
      </div>
    );
  }

}
