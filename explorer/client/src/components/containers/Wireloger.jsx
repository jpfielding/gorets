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
      <div
        className="flex"
        style={{
          maxWidth: '1500px',
          margin: 'auto',
        }}
      >
        <div className="wireloger pa3">
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
