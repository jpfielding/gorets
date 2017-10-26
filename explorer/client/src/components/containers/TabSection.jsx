import React from 'react';

export default class TabSection extends React.Component {

  static propTypes={
    names: React.PropTypes.any,
    components: React.PropTypes.any,
    onRemove: React.PropTypes.Func,
    enableRemove: React.PropTypes.any,
    removeOffset: React.PropTypes.any,
  }

  constructor(props) {
    super(props);
    this.createTabBar = this.createTabBar.bind(this);
    this.state = {
      tab: '',
      tabs: [],
    };
    this.createTabBar = this.createTabBar.bind(this);
    this.createTabTag = this.createTabTag.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.moveTab = this.moveTab.bind(this);
    this.handleKeys = this.handleKeys.bind(this);
  }

  componentWillMount() {
    this.setState({ tab: this.props.names[0] });
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
    this.state.tabs[i].focus();
    const e = this.props.names[i];
    this.setState({ tab: e });
  }

  handleKeys(e, i) {
    if (e.keyCode === 37) {
      this.moveTab(i - 1);
    } else if (e.keyCode === 39) {
      this.moveTab(i + 1);
    }
  }

  createTabTag(e, i) {
    let close = null;
    if (this.props.enableRemove && i >= this.props.removeOffset) {
      close = (
        <button
          className="customClose"
          onClick={() => { this.props.onRemove(e); }}
        >
          X
        </button>
      );
    }
    return (
      <div className="customTabSection">
        <button
          onClick={() => this.setState({ tab: e })}
          onKeyDown={(vent) => this.handleKeys(vent, i)}
          className={`customTabSelect ${this.state.tab === e ? 'active' : ''}`}
          key={e}
          ref={(input) => { this.state.tabs[i] = input; }}
        >
          {e}
        </button>
        {close}
      </div>
    );
  }

  createTabs(e, i) {
    const name = this.props.names[i];
    return (
      <div className={`customTab ${this.state.tab === name ? 'db' : 'dn'}`} key={name}>
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
