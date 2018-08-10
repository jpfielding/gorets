import React from 'react';
import { withRouter } from 'react-router';

class WireEvent extends React.Component {

  static propTypes = {
    event: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      type: '',
    };
    if (props.event.extra) {
      if (props.event.extra.type) this.state.type = props.event.extra.type;
    }
  }

  render() {
    return (
      <div className={`${this.state.type}event`}>
        <h1>
          {this.props.event.tag}
          <button
            className="customButton"
            onClick={() => {
              const open = !this.state.open;
              this.setState({ open });
            }}
          >
            {this.state.open ? '\u25B2' : '\u25BC'}
          </button>
        </h1>
        <pre style={this.state.open ? {} : { maxHeight: '90px' }}>
          {this.state.open ? this.props.event.log : this.props.event.log.slice(0, 500)}
        </pre>
      </div>
    );
  }
}

export default withRouter(WireEvent);
