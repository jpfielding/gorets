import React from 'react';

export default class TabSection extends React.Component {

  static propTypes={
    names: React.PropTypes.any,
    components: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.createTabBar = this.createTabBar.bind(this);
    this.state = {
      tab: 0,
      tabs: [],
    };
    this.createTabBar = this.createTabBar.bind(this);
    this.createTabTag = this.createTabTag.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.moveTab = this.moveTab.bind(this);
    this.handleKeys = this.handleKeys.bind(this);
  }

  createTabBar() {
    const tabs = this.props.names.map(this.createTabTag);
    return (
      <div className="customTabBar">
        {tabs}
      </div>
    );
  }

  moveTab(newTab) {
    let i = newTab;
    if (newTab < 0) {
      i = this.state.tabs.length - 1;
    } else if (newTab > this.state.tabs.length - 1) {
      i = 0;
    }
    this.setState({ tab: i });
    this.state.tabs[i].focus();
  }

  handleKeys(e, i) {
    if (e.keyCode === 37) {
      this.moveTab(i - 1);
    } else if (e.keyCode === 39) {
      this.moveTab(i + 1);
    }
  }

  createTabTag(e, i) {
    return (
      <button
        onClick={() => this.setState({ tab: i })}
        onKeyDown={(vent) => this.handleKeys(vent, i)}
        className={`${this.state.tab === i ? 'active' : ''}`}
        key={e}
        ref={(input) => { this.state.tabs[i] = input; }}
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
