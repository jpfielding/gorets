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

  constructor(props) {
    super(props);
    this.state = {
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
          {this.renderSearchResultsTable()}
          {/* <div>Search parameters:
            <pre className="f6 code">{JSON.stringify(this.state, null, '  ')}</pre>
          </div> */}
        </div>
      </div>
    );
  }

}

export default withRouter(Search);
