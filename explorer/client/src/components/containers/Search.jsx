import React from 'react';
import MetadataService from 'services/MetadataService';
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
      searchRows: [],
      searchColumns: [],
      searchParams: {
        id: null,
        resource: null,
        class: null,
        select: null,
        query: null,
      },
      searchHistory: StorageCache.getFromCache() || [],
      searchResults: {},
    };
    this.search = this.search.bind(this);
    this.onSearchResultCellSelected = this.onSearchResultCellSelected.bind(this);
  }

  componentWillMount() {
    this.search(this.props.location.query);
  }

  onSearchResultCellSelected(coordinates) {
    const { searchResults } = this.state;
    const rows = searchResults.result.rows;
    const selectedRow = rows[coordinates.rowIdx];
    const selectedVal = selectedRow[coordinates.idx];
    console.log(selectedVal);
  }

  setAvailableObjectsState() {
    const resources = this.state.metadata.System['METADATA-RESOURCE'].Resource || [];
    let selectedResource;
    resources.forEach(resource => {
      if (resource.ResourceID === this.state.searchParams.class) {
        selectedResource = resource;
      }
    });
    const { searchResults } = this.state;
    if (!selectedResource || !searchResults.result) {
      return;
    }
    const keyField = selectedResource.KeyField;
    if (!searchResults.result.columns.includes(keyField)) {
      return;
    }
    const metadataObjects = selectedResource['METADATA-OBJECT']['Object'];
    if (metadataObjects.length === 0) {
      return;
    }
    this.setState({
      searchRows: metadataObjects,
      searchColumns: [{
        key: 'ObjectType',
        name: keyField,
      }],
    });
  }

  search(searchParams) {
    this.props.router.push({
      ...this.props.location,
      query: searchParams,
    });
    this.setState({
      searchParams,
    });
    MetadataService
      .search(searchParams)
      .then(res => res.json())
      .then(json => {
        const searchHistory = StorageCache.getFromCache() || [];
        if (!some(searchHistory, searchParams)) {
          searchHistory.push(searchParams);
          StorageCache.putInCache(searchHistory, 60);
        }
        this.setState({
          searchResults: json,
          searchHistory,
        });
        this.setAvailableObjectsState();
      });
    MetadataService
      .get(searchParams.id)
      .then(response => response.json())
      .then(json => {
        if (json.error !== null) {
          this.setState({ metadata: Search.emptyMetadata });
          return;
        }
        this.setState({
          metadata: json.result.Metadata,
        });
        this.setAvailableObjectsState();
      });
  }

  renderObjectMetadata() {
    const { searchRows, searchColumns } = this.state;
    if (searchRows.length === 0 || searchColumns.length === 0) {
      return null;
    }
    const rowGetter = (i) => searchRows[i];
    return (
      <div>
        <div className="b mv2">Object Metadata Types</div>
        <ReactDataGrid
          columns={searchColumns}
          rowGetter={rowGetter}
          rowHeight={35}
          rowsCount={searchRows.length}
          minHeight={(searchRows.length + 1) * 35}
        />
      </div>
    );
  }

  renderSearchResultsTable() {
    const { searchResults } = this.state;
    if (!searchResults.result) {
      return null;
    }
    const rowGetter = (i) => searchResults.result.rows[i];
    const columns = searchResults.result.columns.map((column, index) => ({
      key: index,
      name: column,
      resizable: true,
    }));
    return (
      <ReactDataGrid
        columns={columns}
        rowGetter={rowGetter}
        rowsCount={searchResults.result.rows.length}
        onCellSelected={this.onSearchResultCellSelected}
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
            <div className="b mb2">Search Results</div>
            {this.renderSearchResultsTable()}
            {this.renderObjectMetadata()}
          </div>
          {/* <div>Search parameters:
            <pre className="f6 code">{JSON.stringify(this.state.searchColumns, null, '  ')}</pre>
          </div> */}
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
