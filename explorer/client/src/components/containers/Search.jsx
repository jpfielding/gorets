import React from 'react';
import MetadataService from 'services/MetadataService';

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
      searchResults: [],
    };
  }

  componentWillMount() {
    this.setState({
      searchParams: this.props.location.query,
    });
    MetadataService
      .search(this.props.location.query)
      .then(res => res.json())
      .then(json => this.setState({ searchResults: json }));
  }

  render() {
    return (
      <div>
        <h1>Search</h1>
        <div>Search parameters:
          <pre>{JSON.stringify(this.state, null, '  ')}</pre>
        </div>
      </div>
    );
  }

}
