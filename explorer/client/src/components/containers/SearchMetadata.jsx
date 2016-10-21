import React from 'react';
import { withRouter } from 'react-router';

class SearchMetadata extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      searchParams: {
        id: null,
        resource: null,
        keyFieldValue: null,
        types: null,
      },
    };
  }

  render() {
    return (
      <div>
        <h1>Ayy SearchMetadata</h1>
        <pre className="code f6">
          {JSON.stringify(this.state.searchParams, null, '  ')}
        </pre>
      </div>
    );
  }

}

export default withRouter(SearchMetadata);
