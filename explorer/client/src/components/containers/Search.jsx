import React from 'react';
import MetadataService from 'services/MetadataService';
import StorageCache from 'util/StorageCache';
import some from 'lodash/some';

export default class Search extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
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
      searchResults: [],
    };
  }

  componentWillMount() {
    const searchParams = this.props.location.query;
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

  render() {
    return (
      <div>
        <div className="fl h-100-ns w-100 w-20-ns pa3 overflow-x-scroll nowrap">
          <div className="b">Current Search Params</div>
          <pre className="f6 code">{JSON.stringify(this.state.searchParams, null, '  ')}</pre>
          <div className="b">Search History</div>
          <ul className="pa0 ma0 no-list-style">
            {this.state.searchHistory.map(search =>
              <li>
                <pre className="f6 code">{JSON.stringify(search, null, '  ')}</pre>
              </li>
            )}
          </ul>
        </div>
        <div className="fl h-100 min-vh-100 w-100 w-80-ns pa3 bl-ns">
          <div>Search parameters:
            <pre className="f6 code">{JSON.stringify(this.state, null, '  ')}</pre>
          </div>
        </div>
      </div>
    );
  }

}
