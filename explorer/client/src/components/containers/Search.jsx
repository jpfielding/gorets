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
      metadataRows: [],
      metadataColumns: [],
      searchResultColumns: [],
      searchResultRows: [],
      searchParams: {
        id: null,
        resource: null,
        class: null,
        select: null,
        query: null,
      },
      searchHistory: StorageCache.getFromCache() || [],
      searchResults: {},
      selectedIndexes: [],
    };
    this.search = this.search.bind(this);
    this.onSearchResultCellSelected = this.onSearchResultCellSelected.bind(this);
    this.onRowsSelected = this.onRowsSelected.bind(this);
    this.onRowsDeselected = this.onRowsDeselected.bind(this);
  }

  componentWillMount() {
    this.search(this.props.location.query);
  }

  onRowsSelected(rows) {
    console.log(rows);
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


  onSearchResultCellSelected(coordinates) {
    const {
      searchResultRows,
      searchResultColumns,
      metadataRows,
      metadataColumns,
      searchParams,
    } = this.state;
    const selectedRow = searchResultRows[coordinates.rowIdx];
    const selectedVal = selectedRow[coordinates.idx];
    const selectedColumn = searchResultColumns[coordinates.idx].name;
    const selectableMetdataColumns = metadataColumns.map(x => x.name);
    if (selectableMetdataColumns.includes(selectedColumn)) {
      const availableObjectTypes = metadataRows.map(x => x.ObjectType);
      this.props.router.push({
        ...this.props.location,
        pathname: '/objects',
        query: {
          id: searchParams.id,
          resource: searchParams.resource,
          keyFieldValue: selectedVal,
          types: availableObjectTypes.join(','),
        },
      });
    }
  }

  setAvailableObjectsState() {
    // Search Results table setup
    const { searchResults } = this.state;
    if (!searchResults.result) {
      return;
    }
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
    // Object metadata table setup
    const resources = this.state.metadata.System['METADATA-RESOURCE'].Resource || [];
    let selectedResource;
    resources.forEach(resource => {
      if (resource.ResourceID === this.state.searchParams.resource) {
        selectedResource = resource;
      }
    });
    if (!selectedResource) {
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
      metadataRows: metadataObjects,
      metadataColumns: [{
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
        // StorageCache.clearAll();
        if (!some(searchHistory, searchParams)) {
          searchHistory.unshift(searchParams);
          StorageCache.putInCache(searchHistory, 60);
        }
        console.log(json);
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
        console.log('meta: ', json.result.Metadata);
        this.setAvailableObjectsState();
      });
  }

  renderObjectMetadata() {
    const { metadataRows, metadataColumns } = this.state;
    if (metadataRows.length === 0 || metadataColumns.length === 0) {
      return null;
    }
    const rowGetter = (i) => metadataRows[i];
    return (
      <div>
        <div className="b mv2">Object Metadata Types</div>
        <ReactDataGrid
          columns={metadataColumns}
          rowGetter={rowGetter}
          rowHeight={35}
          rowsCount={metadataRows.length}
          minHeight={(metadataRows.length + 1) * 35}
        />
      </div>
    );
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
            <div className="b mb2">Search Results</div>
            {this.renderSearchResultsTable()}
            {this.renderObjectMetadata()}
          </div>
          <div>Search parameters:
            <pre className="f6 code">{JSON.stringify(this.state.metadataColumns, null, '  ')}</pre>
          </div>
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
