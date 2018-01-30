import React from 'react';
import SearchService from 'services/SearchService';
import StorageCache from 'util/StorageCache';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ReactDataGrid from 'react-data-grid';
import MetadataService from 'services/MetadataService';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import ContentHistory from 'components/containers/History';
import SearchFormatter from 'components/gridcells/SearchFormatter';
import RouteLink from 'components/elements/RouteLink';
import _ from 'underscore';

const Base64 = require('js-base64').Base64;

/*
  Responsible for the search tab of each class
  Sends wirlogs back up to be provided to the wirelogging tab
*/

class Search extends React.Component {

  static propTypes = {
    connection: React.PropTypes.any,
    metadata: React.PropTypes.any,
    location: React.PropTypes.any,
    init: React.PropTypes.any,
    shared: {
      connection: React.PropTypes.any,
      metadata: React.PropTypes.any,
      resource: React.PropTypes.any,
      class: React.PropTypes.any,
      fields: React.PropTypes.any,
      rows: React.PropTypes.any,
    },
    addTab: React.PropTypes.Func,
    pushWirelog: React.PropTypes.Func,
    idprefix: React.PropTypes.any,
  }

  static defaultProps = {
    connection: { id: null },
    metadata: MetadataService.empty,
  }

  constructor(props) {
    super(props);

    const target = (props.init && props.init.query ? props.init.query : SearchService.params);
    const searchForm = createValue({
      value: _.clone(target),
      onChange: this.searchInputsChange.bind(this),
    });

    this.state = {
      searchParams: _.clone(target),
      searchForm,
      searchHistory: [],
      searchResults: {},
      searchResultColumns: [],
      searchResultRows: [],
      searchCount: -1,
      selectedIndexes: [],
      searching: false,
      errorOut: '',
      resultCount: 1,
      tabName: '',
      searchHistoryName: '',
    };

    this.search = this.search.bind(this);

    this.submitSearchForm = this.submitSearchForm.bind(this);

    this.getRowAt = this.getRowAt.bind(this);

    this.setSearchHistory = this.setSearchHistory.bind(this);
    this.createNewTab = this.createNewTab.bind(this);

    this.bindTabNameChange = this.bindTabNameChange.bind(this);
    this.bindQueryNameChange = this.bindQueryNameChange.bind(this);
  }

  componentWillMount() {
    const sck = `${this.props.shared.source}-${this.props.shared.connection.id}-search-history`;
    const searchHistory = StorageCache.getFromCache(sck) || [];
    this.setSearchHistory(searchHistory);
  }

  componentDidMount() {
    if (this.props.init && this.props.init.launch === 'auto') {
      console.log('Auto Searching');
      this.search(this.state.searchParams);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.shared.class['METADATA-TABLE']
      && nextProps.shared !== this.props.shared) {
      const ClassName = nextProps.shared.class.ClassName;
      const resource = nextProps.shared.resource.ResourceID;
      const select = nextProps.shared.fields.map(i => i.row.SystemName).join(',');
      const ts = nextProps.shared.class['METADATA-TABLE'].Field.filter(f => f.StandardName === 'ModificationTimestamp');

      console.log('last modified fields:', ts);

      let query = this.state.searchForm.value.query;
      if (ts.length > 0) {
        const field = ts[0].SystemName.trim();
        const date = JSON.stringify(new Date());
        const day = date.substring(1, 12);
        query = `(${field}=${day}00:00:00+)`;
      }
      const value = this.state.searchForm.value;
      value.resource = resource;
      value.class = ClassName;
      value.select = select;
      value.query = query;

      const searchForm = createValue({
        value,
        onChange: this.searchInputsChange.bind(this),
      });
      this.setState({ searchForm });
    }
  }

  getRowAt(index) {
    if (index < 0) {
      return undefined;
    }
    return this.state.searchResultRows[index];
  }

  setSearchHistory(sh) {
    const searchHistory = _.clone(sh);
    this.setState({
      searchHistory,
    });
  }

  removeHistoryElement(params) {
    console.log('removing history element');
    const sck = `${this.props.shared.source}-${this.props.shared.connection.id}-search-history`;
    const searchHistory = StorageCache.getFromCache(sck) || [];
    if (some(searchHistory, params)) {
      searchHistory.splice(searchHistory.findIndex(i => _.isEqual(i, params)), 1);
      StorageCache.putInCache(sck, searchHistory, 720);
    }
    this.setState({
      searchHistory,
    });
  }

  bindTabNameChange(tabName) {
    this.setState({ tabName });
  }

  bindQueryNameChange(searchHistoryName) {
    this.setState({ searchHistoryName });
  }

  createNewTab() {
    let tabName = this.state.tabName;
    if (tabName === '') {
      tabName = `R${this.state.resultCount}`;
      const resultCount = this.state.resultCount + 1;
      this.setState({ resultCount });
    }
    console.log(`creating new tab of name ${tabName}`);
    this.props.addTab(tabName, this.renderNewTab());
  }

  searchInputsChange(searchForm) {
    this.setState({ searchForm });
  }

  applySearchState() {
    // Search Results table setup
    const { searchResults } = this.state;
    if (!searchResults.result ||
      (searchResults.result.count === 0 &&
        (searchResults.result.rows === null ||
        searchResults.result.columns === null))) {
      if (searchResults.error) {
        this.setState({ errorOut: searchResults.error });
      } else {
        this.setState({ errorOut: 'No Results Found' });
      }
      return;
    }
    if (searchResults.error) {
      this.setState({ errorOut: 'No Results Found' });
      return;
    }

    let searchResultColumns = [];
    let searchResultRows = [];
    let searchCount = -1;
    let errorOut = '';

    console.log('applying search state');

    if (searchResults.result.columns && searchResults.result.rows) {
      searchResultColumns = searchResults.result.columns.map((column, index) => ({
        key: index,
        name: column,
        resizable: true,
        width: 150,
        formatter: <SearchFormatter />,
      }));
      searchResultRows = searchResults.result.rows;
    }
    if (searchResults.result.count) {
      if (searchResults.result.count < 0) {
        errorOut = `Bad Result Count ${searchResults.result.count}`;
        searchCount = -1;
      } else {
        searchCount = searchResults.result.count;
      }
    }
    this.setState({
      searchResultColumns,
      searchResultRows,
      searchCount,
      errorOut,
    });
  }

  submitSearchForm(countType) {
    const form = _.clone(this.state.searchForm.value);
    form['counttype'] = countType;
    form['name'] = this.state.searchHistoryName;
    form['submited'] = true;
    if (form.limit && typeof form.limit === 'string') {
      form.limit = parseInt(form.limit, 10);
    }
    this.search(form);
  }

  search(searchParams) {
    const searchHistoryName = searchParams.name;
    const searchForm = createValue({
      value: _.clone(searchParams),
      onChange: this.searchInputsChange.bind(this),
    });

    const sck = `${this.props.shared.source}-${this.props.shared.connection.id}-search-history`;
    const searchHistory = StorageCache.getFromCache(sck) || [];
    this.setState({
      searchParams,
      searchForm,
      searchHistoryName,
      searching: true,
      errorOut: '',
    });
    if (searchParams === Search.emptySearch) {
      return;
    }
    SearchService
      .search(this.props.shared.connection, searchParams)
      .then(res => {
        console.log(res);
        return res.json();
      }).then(json => {
        if (json.error && json.error !== null) {
          this.props.pushWirelog({ tag: 'Search', log: json.error, extra: { type: 'Error' } });
        } else {
          const log = Base64.decode(json.result.wirelog);
          console.log('adsfwefw', log);
          this.props.pushWirelog({ tag: 'Search', log });
        }
        if (!some(searchHistory, searchParams)) {
          searchHistory.unshift(searchParams);
          StorageCache.putInCache(sck, searchHistory, 720);
        }
        console.log('results:', json);
        this.setState({
          searchResults: json,
          searching: false,
        });
        this.applySearchState();
        this.setSearchHistory(searchHistory);
      });
  }

  renderNewTab() {
    return (
      <div className="ma3 customResultsSet">
        <div className="customResultsTitle">
          <div className="customTitle nonclickable">Results from Search query</div>
        </div>
        <div className="customResultsBody">
          <ContentHistory
            params={this.state.searchParams}
          />
          <div className="pa3">
            {this.renderSearchResultsTable()}
          </div>
        </div>
      </div>
    );
  }

  renderSearchResultsTable() {
    const { searchResultRows, searchResultColumns } = this.state;
    if (searchResultRows.length === 0 || searchResultColumns.length === 0) {
      return null;
    }
    return (
      <ReactDataGrid
        columns={searchResultColumns}
        rowGetter={this.getRowAt}
        rowsCount={searchResultRows.length}
      />
    );
  }

  renderHistoryBar() {
    return (
      <div className="fl w-100 w-20-ns pa3 pr0">
        <div className="customResultsCompactSet">
          <div className="customResultsTitle">
            <div className="b nonclickable">Current Search Params</div>
          </div>
          <div className="customResultsBody">
            <ContentHistory
              clickHandle={() => this.search(this.state.searchParams)}
              params={this.state.searchParams}
              preventClose
              idprefix={`${this.props.idprefix}-history-current`}
            />
          </div>
        </div>
        <div className="customResultsCompactSet">
          <div className="customResultsTitle">
            <div className="b nonclickable">Search History</div>
          </div>
          <div className="customResultsBody">
            <ul className="pa0 ma0 no-list-style customListSpacing">
              {this.state.searchHistory.map((params, i) =>
                <li>
                  <ContentHistory
                    clickHandle={() => this.search(params)}
                    removeHistoryElement={() => this.removeHistoryElement(params)}
                    params={params}
                    idprefix={`${this.props.idprefix}-history-${i}`}
                  />
                </li>
              )}
            </ul>
          </div>
        </div>
      </div>
    );
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
        {this.renderHistoryBar()}
        <div className="fl w-100 w-80-ns pa3">
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <div className="customTitle">
                  Query:
              </div>
              <input
                className="customInputBar pt1"
                placeholder="Query Name"
                onChange={(e) => this.bindQueryNameChange(e.target.value)}
                value={this.state.searchHistoryName}
                id={`${this.props.idprefix}-query-name`}
              />
              <RouteLink
                type={'full'}
                connection={this.props.shared.connection}
                init={{ tab: 'Search', query: this.state.searchForm.value }}
                style={{ float: 'right' }}
                idprefix={`${this.props.idprefix}-query-link`}
              />
            </div>
            <div className="customResultsBody">
              <Fieldset formValue={this.state.searchForm}>
                <Field select="resource" label="Resource">
                  <Input
                    className="w-30 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-resource`}
                  />
                </Field>
                <Field select="class" label="Class">
                  <Input
                    className="w-30 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-class`}
                  />
                </Field>
                <Field select="select" label="Columns">
                  <Input
                    className="w-80 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-select`}
                  />
                </Field>
                <Field select="limit" label="Limit">
                  <Input
                    className="w-80 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-limit`}
                  />
                </Field>
                <Field select="query" label="Query">
                  <Input
                    className="w-80 pa1 b--none outline-transparent"
                    id={`${this.props.idprefix}-query-query`}
                  />
                </Field>
              </Fieldset>
            </div>
            <div className="customResultsFoot">
              <div className="customButtonSet">
                <button
                  onClick={() => this.submitSearchForm(0)}
                  disabled={this.state.searching}
                  className="da-effect"
                  id={`${this.props.idprefix}-submit`}
                >
                  Submit
                </button>
                <button
                  onClick={() => this.submitSearchForm(1)}
                  disabled={this.state.searching}
                  className="da-effect"
                  id={`${this.props.idprefix}-submit-withcount`}
                >
                  with Count
                </button>
                <button
                  onClick={() => this.submitSearchForm(2)}
                  disabled={this.state.searching}
                  className="da-effect"
                  id={`${this.props.idprefix}-submit-countonly`}
                >
                  only Count
                </button>
              </div>
            </div>
          </div>
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <RouteLink
                type={'fullAuto'}
                connection={this.props.shared.connection}
                init={{ tab: 'Search', query: this.state.searchForm.value }}
                style={{ float: 'right' }}
                idprefix={`${this.props.idprefix}-results-link`}
              />
              <div className="customTitle">
                  Search Results: {this.state.searchResults.error ? (`${this.state.searchResults.error}`) : ''}
              </div>
            </div>
            <div className="customResultsBody">
              <div className={`bg-dark-red white br1 pa2 ${this.state.errorOut.length === 0 ? 'dn' : 'dib'}`}>
                {this.state.errorOut}
              </div>
              <div className={`${this.state.searchCount >= 0 ? 'dib' : 'dn'}`}>
                Count: {this.state.searchCount}
              </div>
              <div>
                {this.renderSearchResultsTable()}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
