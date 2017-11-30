import React from 'react';

class DragableWindow extends React.Component {

  static propTypes = {
    content: React.PropTypes.any,
    title: React.PropTypes.any,
    removeSelf: React.PropTypes.Func,
    container: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      ref: null,
    };

    this.startMove = this.startMove.bind(this);
    this.stopMove = this.stopMove.bind(this);
  }

  startMove(evt) {
    const posX = evt.clientX;
    const posY = evt.clientY;
    const divTop = this.state.ref.style.top.replace('px', '');
    const divLeft = this.state.ref.style.left.replace('px', '');
    const diffX = posX - divLeft;
    const diffY = posY - divTop;
    console.log(posX, posY, divTop, divLeft, diffX, diffY);
    document.onmousemove = (e) => {
      const setX = e.clientX - diffX;
      let setY = e.clientY - diffY;
      if (setY < 175) setY = 175;
      this.state.ref.style.left = `${setX}px`;
      this.state.ref.style.top = `${setY}px`;
    };
  }

  stopMove() {
    document.onmousemove = () => {};
  }

  render() {
    return (
      <div
        ref={(ref) => { this.state.ref = ref; }}
        className="dragableWindow"
        style={{
          left: '200px',
          top: '200px',
        }}
      >
        <div
          className="title"
          onMouseDown={this.startMove}
          onMouseUp={this.stopMove}
        >
          {this.props.title}
          <div style={{ flex: '1' }} />
          {this.props.removeSelf ? <button
            className="customButton"
            onClick={this.props.removeSelf}
          > X </button> : null}
        </div>
        <div className="body">
          {this.props.content}
        </div>
      </div>
    );
  }

}

export default DragableWindow;
