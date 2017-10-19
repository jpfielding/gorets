import React from 'react';

export default class TabSection extends React.Component {

  static propTypes={
    names: React.PropTypes.any,
    components: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.state = {
      tab: 0,
    };
    this.createTabBar = this.createTabBar.bind(this);
    this.createTabTag = this.createTabTag.bind(this);
    this.createTabs = this.createTabs.bind(this);
  }

  createTabBar() {
    return (
      <div className="customTabBar">
        {this.props.names.map(this.createTabTag)}
      </div>
    );
  }

  createTabTag(e, i) {
    return (
      <button
        onClick={() => this.setState({ tab: i })}
        className={`${this.state.tab === i ? 'active' : ''}`}
        key={e}
      >
        {e}
      </button>
    );
  }

  createTabs(e, i) {
    return (
      <div className={`customTab ${this.state.tab === i ? 'db' : 'dn'}`} key={i}>
        {e}
      </div>
    );
  }

  render() {
    return (
      <div className="customTabElement">
        {this.createTabBar()}
        {this.props.components.map(this.createTabs)}
      </div>
    );
  }
}
