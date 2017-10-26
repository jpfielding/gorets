import React from 'react';
import SearchService from 'services/SearchService';
import StorageCache from 'util/StorageCache';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ReactDataGrid from 'react-data-grid';
import MetadataService from 'services/MetadataService';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import SearchHistory from 'components/containers/SearchHistory';

class Search extends React.Component {

  static propTypes = {
    connection: React.PropTypes.any,
    metadata: React.PropTypes.any,
    location: React.PropTypes.any,
    router: React.PropTypes.any,
    shared: {
      connection: React.PropTypes.any,
      metadata: React.PropTypes.any,
      resource: React.PropTypes.any,
      class: React.PropTypes.any,
      fields: React.PropTypes.any,
      rows: React.PropTypes.any,
    },
    onRowsSelected: React.PropTypes.Func,
    onRowsDeselected: React.PropTypes.Func,
    addTab: React.PropTypes.Func,
  }

  static defaultProps = {
    connection: { id: null },
    metadata: MetadataService.empty,
  }

  constructor(props) {
    super(props);

    const searchForm = createValue({
      value: {
        resource: props.shared.resource.ResourceID,
        class: props.shared.class.ClassName,
        query: null,
      },
      onChange: this.searchInputsChange.bind(this),
    });

    this.state = {
      searchParams: SearchService.params,
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
    this.onRowsSelected = this.onRowsSelected.bind(this);
    this.onRowsDeselected = this.onRowsDeselected.bind(this);
    this.submitSearchForm = this.submitSearchForm.bind(this);
    this.submitSearchFormWithCount = this.submitSearchFormWithCount.bind(this);
    this.submitSearchFormOnlyCount = this.submitSearchFormOnlyCount.bind(this);
    this.getRowAt = this.getRowAt.bind(this);
    this.setSearchHistory = this.setSearchHistory.bind(this);
    this.createNewTab = this.createNewTab.bind(this);
    this.tabNameChange = this.tabNameChange.bind(this);
  }

  componentWillMount() {
    this.setSearchHistory();
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.shared.class['METADATA-TABLE']) {
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
      const searchForm = createValue({
        value: { resource, class: ClassName, select, query },
        onChange: this.searchInputsChange.bind(this),
      });
      this.setState({ searchForm });
    }
  }

  onRowsSelected(rows) {
    this.props.onRowsSelected(rows);
  }

  onRowsDeselected(rows) {
    this.props.onRowsDeselected(rows);
  }

  getRowAt(index) {
    if (index < 0) {
      return undefined;
    }
    return this.state.searchResultRows[index];
  }

  setSearchHistory() {
    const sck = `${this.props.shared.connection.id}-search-history`;
    let searchHistory = StorageCache.getFromCache(sck) || [];
    if (searchHistory && searchHistory.length > 0) {
      searchHistory = searchHistory.filter((i) => (i.query));
    }
    this.setState({
      searchHistory,
    });
  }

  tabNameChange(tabName) {
    this.setState({ tabName });
  }

  queryNameChange(searchHistoryName) {
    this.setState({ searchHistoryName });
  }

  createNewTab() {
    let tabName = this.state.tabName;
    if (tabName === '') {
      tabName = `R${this.state.resultCount}`;
    }
    this.props.addTab(tabName, this.renderNewTab());
    const resultCount = this.state.resultCount + 1;
    this.setState({ resultCount });
  }

  searchInputsChange(searchForm) {
    this.setState({ searchForm });
  }

  applySearchState() {
    // Search Results table setup
    const { searchResults } = this.state;
    if (!searchResults.result) {
      this.setState({ errorOut: 'No Results Found' });
      return;
    }
    console.log('[SEARCH] Applying search state');
    if (searchResults.result.columns && searchResults.result.rows) {
      const searchResultColumns = searchResults.result.columns.map((column, index) => ({
        key: index,
        name: column,
        resizable: true,
        width: 150,
      }));
      const searchResultRows = searchResults.result.rows;
      this.setState({
        searchResultColumns,
        searchResultRows,
      });
    } else {
      this.setState({
        searchResultColumns: [],
        searchResultRows: [],
      });
    }
    if (searchResults.result.count) {
      if (searchResults.result.count < 0) {
        this.setState({
          errorOut: `Bad Result Count ${searchResults.result.count}`,
          searchCount: -1,
        });
      } else {
        this.setState({ searchCount: searchResults.result.count });
      }
    } else {
      this.setState({ searchCount: -1 });
    }
  }

  submitSearchForm() {
    this.search({
      id: this.props.shared.connection.id,
      ...this.state.searchForm.value,
    });
  }

  submitSearchFormWithCount() {
    const form = Object.assign({}, this.state.searchForm.value);
    form['id'] = this.props.shared.connection.id;
    form['count-type'] = 1;
    this.search(form);
  }

  submitSearchFormOnlyCount() {
    const form = Object.assign({}, this.state.searchForm.value);
    form['id'] = this.props.shared.connection.id;
    form['count-type'] = 2;
    this.search(form);
  }

  search(sp) {
    const searchParams = sp;
    let searchHistoryName = this.state.searchHistoryName;
    if (searchParams.name) {
      searchHistoryName = searchParams.name;
    }
    if (this.state.searchHistoryName !== '') {
      searchParams.name = this.state.searchHistoryName;
    }

    const searchForm = this.state.searchForm;
    searchForm.value.resource = searchParams.resource;
    searchForm.value.class = searchParams.class;
    searchForm.value.query = searchParams.query;
    searchForm.value.select = searchParams.select;
    // search history cache key used for storage
    const sck = `${this.props.shared.connection.id}-search-history`;
    const searchHistory = StorageCache.getFromCache(sck) || [];
    this.setState({
      searchParams,
      searchForm,
      searchHistoryName,
    });
    this.setState({ searching: true, errorOut: '' });
    if (searchParams === Search.emptySearch) {
      return;
    }
    console.log('[SEARCH] Cache key found', sck);
    SearchService
      .search(searchParams)
      .then(res => res.json())
      .then(json => {
        if (!some(searchHistory, searchParams)) {
          searchHistory.unshift(searchParams);
          StorageCache.putInCache(sck, searchHistory, 720);
        }
        console.log('[SEARCH] Results:', json);
        this.setState({
          searchResults: json,
        });
        this.applySearchState();
        this.setState({ searching: false });
        this.setSearchHistory();
      });
  }

  renderNewTab() {
    return (
      <div className="pa3">
        <div className="b nonclickable">Results from Search query</div>
        <SearchHistory
          params={this.state.searchParams}
        />
        {this.renderSearchResultsTable()}
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
      <div className="fl w-100 w-20-ns pa3">
        <div className="customResultsCompactSet">
          <div className="customResultsTitle">
            <div className="b nonclickable">Current Search Params</div>
          </div>
          <div className="customResultsBody">
            <SearchHistory
              onClick={() => this.search(this.state.searchParams)}
              params={this.state.searchParams}
            />
          </div>
        </div>
        <div className="customResultsCompactSet">
          <div className="customResultsTitle">
            <div className="b nonclickable">Search History</div>
          </div>
          <div className="customResultsBody">
            <ul className="pa0 ma0 no-list-style">
              {this.state.searchHistory.map(params =>
                <li>
                  <SearchHistory
                    className="clickable"
                    onClick={() => this.search(params)}
                    params={params}
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
      <div className="min-vh-100 flex">
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
                onChange={(e) => this.queryNameChange(e.target.value)}
                value={this.state.searchHistoryName}
              />
            </div>
            <div className="customResultsBody">
              <Fieldset formValue={this.state.searchForm}>
                <Field select="resource" label="Resource">
                  <Input className="w-30 pa1 b--none outline-transparent" />
                </Field>
                <Field select="class" label="Class">
                  <Input className="w-30 pa1 b--none outline-transparent" />
                </Field>
                <Field select="select" label="Columns">
                  <Input className="w-80 pa1 b--none outline-transparent" />
                </Field>
                <Field select="query" label="Query">
                  <Input className="w-80 pa1 b--none outline-transparent" />
                </Field>
              </Fieldset>
            </div>
            <div className="customResultsFoot">
              <div className="customButtonSet">
                <button
                  onClick={this.submitSearchForm}
                  disabled={this.state.searching}
                  className="da-effect"
                >
                  Submit
                </button>
                <button
                  onClick={this.submitSearchFormWithCount}
                  disabled={this.state.searching}
                  className="da-effect"
                >
                  with Count
                </button>
                <button
                  onClick={this.submitSearchFormOnlyCount}
                  disabled={this.state.searching}
                  className="da-effect"
                >
                  only Count
                </button>
              </div>
            </div>
          </div>
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <div className="customCombo fr">
                <button className="customComboButton" onClick={this.createNewTab}>New Tab</button>
                <input
                  className="customComboInput"
                  placeholder={`R${this.state.resultCount}`}
                  onChange={(e) => this.tabNameChange(e.target.value)}
                />
              </div>
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
