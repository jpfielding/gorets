import React from 'react';

export default class NewUrl extends React.Component {

  static propTypes = {
    submit: React.PropTypes.Func,
    close: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);
    this.state = {
      name: '',
      url: '',
    };
    this.bindURL = this.bindURL.bind(this);
    this.bindName = this.bindName.bind(this);
  }

  bindURL(url) {
    this.setState({ url });
  }

  bindName(name) {
    this.setState({ name });
  }

  render() {
    return (
      <div
        style={{ maxWidth: '800px', margin: '100px auto' }}
        className="customResultsSet shadowed"
      >
        <div className="customResultsTitle">
          <div className="customTitle di">New URL</div>
          <button
            className="customButton fr" onClick={() => {
              this.props.close();
            }}
          > Close </button>
        </div>
        <div className="customResultsBody bg-mainbg" style={{ display: 'flex' }}>
          <input
            placeholder="Name"
            className="customInput"
            onChange={(e) => this.bindName(e.target.value)}
          />
          <input
            placeholder="URL"
            style={{ flex: '1' }}
            className="customInput"
            onChange={(e) => this.bindURL(e.target.value)}
          />
        </div>
        <div className="customResultsFoot">
          <button
            className="customButton"
            onClick={() => {
              this.props.submit(this.state);
              this.props.close();
            }}
          >Submit</button>
        </div>
      </div>
    );
  }
}
