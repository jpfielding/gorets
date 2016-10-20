import React from 'react';
import MetadataService from 'services/MetadataService';
import ReactDataGrid from 'react-data-grid';

export default class Explorer extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
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
      metadata: Explorer.emptyMetadata,
      selectedClass: null,
      defaultRows: [],
      selectedClassRows: [],
    };
    this.handleGridSort = this.handleGridSort.bind(this);
  }

  componentDidMount() {
    if (this.props.params.connection) {
      this.getMetadata(this.props.params.connection);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.params !== this.props.params && nextProps.params.connection) {
      this.getMetadata(nextProps.params.connection);
    } else {
      this.setState({
        metadata: Explorer.emptyMetadata,
      });
    }
  }

  getMetadata(connectionId) {
    this.setState({
      selectedClass: null,
      defaultRows: null,
      selectedClassRows: [],
    });
    MetadataService
      .get(connectionId)
      .then(response => response.json())
      .then(json => {
        this.setState({
          metadata: json,
        });
      });
  }

  metadataClassClick(selectedClass) {
    this.setState({
      selectedClass,
      defaultRows: selectedClass['METADATA-TABLE'].Field,
      selectedClassRows: selectedClass['METADATA-TABLE'].Field,
    });
  }

  handleGridSort(sortColumn, sortDirection) {
    const comparer = (a, b) => {
      if (sortDirection === 'ASC') {
        if (a[sortColumn] && b[sortColumn] && typeof a[sortColumn] === 'string') {
          return (a[sortColumn].toLowerCase() > b[sortColumn].toLowerCase()) ? 1 : -1;
        }
        return (a[sortColumn] > b[sortColumn]) ? 1 : -1;
      } else if (sortDirection === 'DESC') {
        if (a[sortColumn] && b[sortColumn] && typeof a[sortColumn] === 'string') {
          return (a[sortColumn].toLowerCase() < b[sortColumn].toLowerCase()) ? 1 : -1;
        }
        return (a[sortColumn] < b[sortColumn]) ? 1 : -1;
      }
      return null;
    };
    const rows = sortDirection === 'NONE'
      ? this.state.defaultRows
      : this.state.selectedClassRows.sort(comparer);
    this.setState({ selectedClassRows: rows });
  }

  render() {
    const { selectedClass } = this.state;
    let tableBody;
    if (selectedClass) {
      const fields = selectedClass['METADATA-TABLE'].Field;
      const availableFields = [];
      fields.forEach(field => {
        Object.keys(field).forEach(key => {
          if (availableFields.includes(key)) {
            return;
          }
          availableFields.push(key);
        });
      });
      const fieldSet = (fields && fields.length > 0)
        ? availableFields.map((name) => ({
          key: name,
          name,
          resizable: true,
          width: 200,
          sortable: true,
        }))
        : [];

      const rowGetter = (i) => fields[i];

      tableBody = (
        <div>
          <ReactDataGrid
            onGridSort={this.handleGridSort}
            columns={fieldSet}
            rowGetter={rowGetter}
            rowsCount={fields.length}
            minHeight={500}
          />
        </div>
      );
    } else {
      tableBody = null;
    }
    return (
      <div>
        <div className="fl w-100 w-30-ns no-list-style pa3">
          <ul className="pa0 ma0">
            {this.state.metadata.System['METADATA-RESOURCE'].Resource.map((resource) =>
              <li className="mb2">
                <div className="b">{resource.ResourceID}</div>
                <ul className="pa0 pl3 mv1">
                  {resource['METADATA-CLASS'].Class.map((mClass) =>
                    <li
                      onClick={() => this.metadataClassClick(mClass)}
                      className="clickable metadata-class-name"
                    >
                      {mClass.ClassName}
                    </li>
                  )}
                </ul>
              </li>
            )}
          </ul>
        </div>
        <div className="fl h-100 min-vh-100 w-100 w-70-ns pa3 bl-ns">
          {
            tableBody || <h1 className="f4">Please select a class to explore</h1>
          }
        </div>
      </div>
    );
  }
}
