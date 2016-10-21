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
  }

  componentWillMount() {
    this.search(this.props.location.query);
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
      });
    MetadataService
      .get(searchParams.id)
      .then(response => response.json())
      .then(json => {
        if (json.error !== null) {
          this.setState({ metadata: Search.emptyMetadata });
          return;
        }
        console.log(json.result.Metadata);
        this.setState({
          metadata: json.result.Metadata,
        });
      });
  }

  renderObjectMetadata() {
    const resources = this.state.metadata.System['METADATA-RESOURCE'].Resource || [];
    let selectedResource;
    resources.forEach(resource => {
      if (resource.ResourceID === this.state.searchParams.class) {
        selectedResource = resource;
      }
    });
    if (!selectedResource) {
      return null;
    }
    const keyField = selectedResource.KeyField;
    const metadataObjects = selectedResource['METADATA-OBJECT']['Object'];
    if (metadataObjects.length === 0) {
      return null;
    }
    const columns = [{
      key: 'ObjectType',
      name: keyField,
    }];
    const rowGetter = (i) => metadataObjects[i];
    return (
      <div>
        <ReactDataGrid
          columns={columns}
          rowGetter={rowGetter}
          rowHeight={35}
          rowsCount={metadataObjects.length}
          minHeight={(metadataObjects.length + 1) * 35}
        />
        {/* <pre className="f6 code">{JSON.stringify({ keyField, metadataObjects }, null, '  ')}</pre> */}
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
            <div className="b mv2">Object Metadata Types</div>
            {this.renderObjectMetadata()}
          </div>
          {/* <div>Search parameters:
            <pre className="f6 code">{JSON.stringify(this.state, null, '  ')}</pre>
          </div> */}
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
