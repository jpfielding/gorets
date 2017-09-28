import React from 'react';
import MetadataService from 'services/MetadataService';
import SearchService from 'services/SearchService';
import StorageCache from 'util/StorageCache';
import { withRouter } from 'react-router';
import some from 'lodash/some';
import ReactDataGrid from 'react-data-grid';

class Search extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
    router: React.PropTypes.any,
  }

  static emptyMetadata = {
    System: {
      'METADATA-RESOURCE': {
        Resource: [],
      },
    },
  };

  constructor(props) {
    super(props);
    this.state = {
      metadata: Search.emptyMetadata,
      searchResultColumns: [],
      searchResultRows: [],
      searchParams: {
        id: null,
        resource: null,
        class: null,
        select: null,
        query: null,
      },
      searchHistory: [],
      searchResults: {},
      selectedIndexes: [],
    };
    this.search = this.search.bind(this);
    this.onRowsSelected = this.onRowsSelected.bind(this);
    this.onRowsDeselected = this.onRowsDeselected.bind(this);
  }

  componentWillMount() {
    this.search(this.props.location.query);
  }

  onRowsSelected(rows) {
    this.setState({
      selectedIndexes: this.state.selectedIndexes.concat(rows.map(r => r.rowIdx)),
    });
  }

  onRowsDeselected(rows) {
    const rowIndexes = rows.map(r => r.rowIdx);
    this.setState({
      selectedIndexes: this.state.selectedIndexes.filter(i => rowIndexes.indexOf(i) === -1),
    });
  }


  getObjects() {
    const {
      searchResultRows,
      searchParams,
      selectedIndexes,
    } = this.state;
    const keyFieldCol = this.getKeyFieldColumn();
    const selectedRows = selectedIndexes.map(i => searchResultRows[i]);
    console.log('rows', selectedRows);
    const ids = selectedRows.map(r => r[keyFieldCol.key]);
    if (ids.length === 0) {
      console.log('no selected ids', selectedIndexes);
      return;
    }
    console.log('ids', ids);
    this.props.router.push({
      ...this.props.location,
      pathname: '/objects',
      query: {
        id: searchParams.id,
        resource: searchParams.resource,
        ids: ids.join(','),
        types: this.getObjectTypes().join(','),
      },
    });
  }

  getKeyFieldColumn() {
    const { searchResultColumns } = this.state;
    const keyField = this.getResource().KeyField;
    const keyFieldCols = searchResultColumns.filter(c => (c.name === keyField));
    if (keyFieldCols.length === 0) {
      return null;
    }
    return keyFieldCols[0];
  }

  getObjectTypes() {
    return this.getResource()['METADATA-OBJECT']['Object'].map(o => o.ObjectType) || [];
  }

  getResource() {
    // TODO errors?  those can happen?
    return this.state.metadata.System['METADATA-RESOURCE'].Resource.filter(
      r => (r.ResourceID === this.state.searchParams.resource)
    )[0];
  }

  applySearchState() {
    // Search Results table setup
    const { searchResults } = this.state;
    if (!searchResults.result) {
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

  // does the current state and selections support an object request
  canPullObjects() {
    if (this.state.selectedIndexes.length === 0) {
      return false;
    }
    if (this.getObjectTypes().length === 0) {
      return false;
    }
    if (this.getKeyFieldColumn() == null) {
      return false;
    }
    return true;
  }

  search(searchParams) {
    this.props.router.push({
      ...this.props.location,
      query: searchParams,
    });
    // search history cache key used for storage
    const sck = `${searchParams.id}-search-history`;
    const searchHistory = StorageCache.getFromCache(sck) || [];
    this.setState({
      searchParams,
      searchHistory,
    });
    console.log('source id:', searchParams.id);
    if (searchParams === undefined || searchParams.id === undefined) {
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
      });
    const mck = `${searchParams.id}-metadata`;
    const md = StorageCache.getFromCache(mck);
    if (md) {
      this.setState({ metadata: md });
    } else {
      MetadataService
        .get(searchParams.id)
        .then(response => response.json())
        .then(json => {
          if (json.error !== null) {
            this.setState({ metadata: Search.emptyMetadata });
            return;
          }
          this.setState({ metadata: json.result.Metadata });
          console.log('meta: ', json.result.Metadata);
          StorageCache.putInCache(mck, json.result.Metadata, 60);
          this.applySearchState();
        });
    }
  }

  renderSearchResultsTable() {
    const { searchResultRows, searchResultColumns } = this.state;
    if (searchResultRows.length === 0 || searchResultColumns.length === 0) {
      return null;
    }
    const rowGetter = (i) => searchResultRows[i];
    return (
      <ReactDataGrid
        columns={searchResultColumns}
        rowGetter={rowGetter}
        rowsCount={searchResultRows.length}
        rowSelection={{
          showCheckbox: true,
          enableShiftSelect: true,
          onRowsSelected: this.onRowsSelected,
          onRowsDeselected: this.onRowsDeselected,
          selectBy: {
            indexes: this.state.selectedIndexes,
          },
        }}
      />
    );
  }

  render() {
    return (
      <div>
        <div className="fl h-100-ns w-100 w-20-ns pa3 overflow-x-scroll nowrap">
          <div className="b">Current Search Params</div>
          <pre className="f6 code">{JSON.stringify(this.state.searchParams, null, '  ')}</pre>
          <div className="b">Search History</div>
          <ul className="pa0 ma0 no-list-style">
            {this.state.searchHistory.map(params =>
              <li>
                <pre
                  className="f6 code clickable"
                  onClick={() => this.search(params)}
                >
                  { JSON.stringify(params, null, '  ') }
                </pre>
              </li>
            )}
          </ul>
        </div>
        <div className="fl h-100 min-vh-100 w-100 w-80-ns pa3 bl-ns">
          <div>
            <div className="b mb2">Search Results:</div>
            {this.renderSearchResultsTable()}
          </div>
          <div className="b mb2">Operations:</div>
          <div>
            <button
              disabled={!this.canPullObjects()}
              className="link"
              onClick={() => this.getObjects()}
            >
              Objects: {this.getObjectTypes().join(', ')}
            </button>
          </div>
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
