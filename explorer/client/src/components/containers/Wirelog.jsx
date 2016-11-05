import React from 'react';
import Websocket from 'react-websocket';

class Wirelog extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      wirelogParams: {
        id: null,
        offset: 0,
      },
      log: '',
    };
    this.handleData.bind(this);
  }

  handleData(data) {
    const result = JSON.parse(data);
    this.setState({
      log: this.state.log + result.chunk,
      currentSize: result.Size,
      wirelogParams: {
        offset: this.state.offset + result.Length,
      },
    });
  }

  render() {
    return (
      <div>
        Wirelog: <div >{this.state.count}</div>
        <pre className="f6 code">{this.state.log}</pre>
        <Websocket
          url="ws://localhost:8888/wirelog/"
          onMessage={this.handleData}
        />
      </div>
    );
  }
}

export default Wirelog;
