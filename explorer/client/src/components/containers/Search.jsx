import React from 'react';
import SearchService from 'services/SearchService';
import StorageCache from 'util/StorageCache';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ReactDataGrid from 'react-data-grid';
import MetadataService from 'services/MetadataService';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import HistoryElement from 'components/containers/HistoryElement';

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
      selectedIndexes: [],
      searching: false,
      errorOut: '',
    };

    this.search = this.search.bind(this);
    this.onRowsSelected = this.onRowsSelected.bind(this);
    this.onRowsDeselected = this.onRowsDeselected.bind(this);
    this.submitSearchForm = this.submitSearchForm.bind(this);
    this.getRowAt = this.getRowAt.bind(this);
  }

  componentWillMount() {
      // search history cache key used for storage
    const sck = `${this.props.shared.connection.id}-search-history`;
    const searchHistory = StorageCache.getFromCache(sck) || [];
    let searchParams = Search.emptySearch;
    if (searchHistory.length > 0) {
      searchParams = searchHistory[0];
    }
    this.setState({
      searchParams,
      searchHistory,
    });
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

  searchInputsChange(searchForm) {
    this.setState({ searchForm });
  }

  applySearchState() {
    // Search Results table setup
    const { searchResults } = this.state;
    if (!searchResults.result || !searchResults.result.columns) {
      this.setState({ errorOut: 'No Results Found' });
      return;
    }
    console.log('setting search state');
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
  }

  submitSearchForm() {
    this.search({
      id: this.props.shared.connection.id,
      ...this.state.searchForm.value,
    });
  }

  search(searchParams) {
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
      searchHistory,
      searchForm,
    });
    this.setState({ searching: true, errorOut: '' });
    console.log('source id:', searchParams.id);
    if (searchParams === Search.emptySearch) {
      return;
    }
    console.log('cache key found', sck);
    SearchService
      .search(searchParams)
      .then(res => res.json())
      .then(json => {
        if (!some(searchHistory, searchParams)) {
          searchHistory.unshift(searchParams);
          StorageCache.putInCache(sck, searchHistory, 720);
        }
        console.log(json);
        this.setState({
          searchResults: json,
          searchHistory,
        });
        this.applySearchState();
        this.setState({ searching: false });
      });
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
        rowSelection={{
          showCheckbox: false,
          enableShiftSelect: true,
          onRowsSelected: this.onRowsSelected,
          onRowsDeselected: this.onRowsDeselected,
          selectBy: {
            indexes: this.props.shared.data.map(d => d.rowIdx),
          },
        }}
      />
    );
  }

  render() {
    return (
      <div className="min-vh-100 flex">
        <div className="fl w-100 w-20-ns pa3 overflow-x-scroll">
          <div className="b">Current Search Params</div>
          <HistoryElement
            className="clickable"
            onClick={() => this.search(this.state.searchParams)}
            params={this.state.searchParams}
          />
          <div className="b">Search History</div>
          <ul className="pa0 ma0 no-list-style">
            {this.state.searchHistory.slice(1).map(params =>
              <li>
                <HistoryElement
                  className="clickable"
                  onClick={() => this.search(params)}
                  params={params}
                />
              </li>
            )}
          </ul>
        </div>
        <div className="fl w-100 w-80-ns pa3 bl-ns">
          <div className="pb3">
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
              <button
                onClick={this.submitSearchForm}
                disabled={this.state.searching}
                className="ba black bg-transparent b--black pa1 mt2"
              >
                Submit
              </button>
            </Fieldset>
          </div>
          <div>
            <div className="b mb2">
                Search Results: {this.state.searchResults.error ? (`${this.state.searchResults.error}`) : ''}
            </div>
            <div className={`bg-dark-red white br1 pa2 ${this.state.errorOut.length === 0 ? 'dn' : 'dib'}`}>
              {this.state.errorOut}
            </div>
            {this.renderSearchResultsTable()}
          </div>
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
