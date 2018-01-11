import React from 'react';

export default class TabSection extends React.Component {

  static propTypes={
    components: React.PropTypes.any,
    onRemove: React.PropTypes.Func,
    enableRemove: React.PropTypes.any,
    removeOffset: React.PropTypes.any,

    className: React.PropTypes.any,
    initID: React.PropTypes.any,

    tag: React.PropTypes.any,
  }

  static defaultProps = {
    tag: '',
  }

  constructor(props) {
    super(props);

    this.state = {
      tab: (this.props.initID ? this.props.initID : ''),
      tabs: [],
    };

    if (this.props.components[0] && this.state.tab === '') {
      this.state.tab = this.props.components[0].id;
    }

    this.createTabBar = this.createTabBar.bind(this);
    this.createTabTag = this.createTabTag.bind(this);
    this.createTabs = this.createTabs.bind(this);
    this.moveTab = this.moveTab.bind(this);
    this.handleKeys = this.handleKeys.bind(this);
  }

  createTabBar() {
    const tabs = this.props.components.map(this.createTabTag);
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
    const e = this.props.components[i].id;
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
          onClick={() => { this.props.onRemove(e.id); }}
          id={this.props.tag === '' ? `${e.idprefix}-tab-close` : `${this.props.tag}-${e.idprefix}-tab-close`}
        >
          X
        </button>
      );
    }
    let tags = null;
    if (e.tags) {
      tags = e.tags.map((tag) => (
        <div style={{ backgroundColor: tag['color'] }} className="activeFullTag" >
          {tag.name}
        </div>
      ));
    }
    return (
      <div className="customTabSection">
        <button
          onClick={() => this.setState({ tab: e.id })}
          onKeyDown={(vent) => this.handleKeys(vent, i)}
          className={`customTabSelect ${this.state.tab === e.id ? 'active' : ''}`}
          key={e.id}
          id={this.props.tag === '' ? `${e.idprefix}-tab` : `${this.props.tag}-${e.idprefix}-tab`}
          ref={(input) => { this.state.tabs[i] = input; }}
        >
          {tags ? tags[0] : null}
          {e.name ? e.name : e.id}
        </button>
        {close}
      </div>
    );
  }

  createTabs(e) {
    const name = e.id;
    return (
      <div className={`customTab ${this.state.tab === name ? 'db' : 'dn'}`} key={name}>
        {e.page}
      </div>
    );
  }

  render() {
    return (
      <div className={`${this.props.className} customTabElementU`}>
        {this.createTabBar()}
        {this.props.components.map(this.createTabs)}
      </div>
    );
  }
}
