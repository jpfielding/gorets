import React from 'react';
import { withRouter } from 'react-router';
import WireEvent from 'components/containers/WireEvent';

class Wireloger extends React.Component {

  static propTypes = {
    wirelog: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      extention: 0,
      heights: ['0vh', '20vh', '81vh'],
    };
  }

  render() {
    return (
      <div className="wireloger">
        <div className="pulltab" >
          <div className="expand">
            <button
              className={this.state.extention < 2 ? '' : 'dn'}
              onClick={() => {
                const extention = this.state.extention + 1;
                this.setState({ extention });
              }}
            >
              {'\u25B2'}
            </button>
            <button
              className={this.state.extention > 0 ? '' : 'dn'}
              onClick={() => {
                const extention = this.state.extention - 1;
                this.setState({ extention });
              }}
            >
              {'\u25BC'}
            </button>
          </div>
        </div>
        <div
          className={`body ${this.state.extention > 0 ? '' : 'fdn'}`}
          style={{ height: this.state.heights[this.state.extention] }}
        >
          <h1 className="title">Wirelog</h1>
          <div className="scrollable">
            {this.props.wirelog.map((e, i) => (
              <WireEvent key={this.props.wirelog.length - i} event={e} />
              ))}
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(Wireloger);
