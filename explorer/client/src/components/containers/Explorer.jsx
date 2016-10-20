import React from 'react';
import MetadataService from 'services/MetadataService';
import ReactDataGrid from 'react-data-grid';

const ReactDataGridPlugins = require('react-data-grid/addons');

const Toolbar = ReactDataGridPlugins.Toolbar;

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
      filters: {},
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

  onClearFilters = () => {
    // all filters removed
    this.setState({ filters: {} });
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
      const aVal = a[sortColumn] ? String(a[sortColumn]).toLowerCase() : '';
      const bVal = b[sortColumn] ? String(b[sortColumn]).toLowerCase() : '';
      if (sortDirection === 'ASC') {
        return (aVal > bVal) ? 1 : -1;
      } else if (sortDirection === 'DESC') {
        return (aVal < bVal) ? 1 : -1;
      }
      return null;
    };
    const rows = sortDirection === 'NONE'
      ? this.state.selectedClassRows
      : this.state.selectedClassRows.sort(comparer);
    this.setState({ selectedClassRows: rows });
  }

  handleFilterChange = (filter) => {
    const newFilters = Object.assign({}, this.state.filters);
    if (filter.filterTerm) {
      newFilters[filter.column.key] = filter;
    } else {
      delete newFilters[filter.column.key];
    }
    this.setState({ filters: newFilters });
    const rows = this.state.defaultRows;
    const newRows = [...rows];
    Object.keys(newFilters).forEach(newFilter => {
      const filterObj = newFilters[newFilter];
      if (filterObj.filterTerm) {
        for (let i = rows.length - 1; i >= 0; i--) {
          const row = rows[i];
          const val = row[filterObj.column.key] || '';
          const stringVal = String(val);
          if (stringVal.indexOf(filterObj.filterTerm) === -1) {
            newRows.splice(i, 1);
          }
        }
      }
    });
    this.setState({ selectedClassRows: newRows });
  }

  render() {
    const { selectedClassRows } = this.state;
    let tableBody;
    if (selectedClassRows) {
      const fields = selectedClassRows;
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
          filterable: true,
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
            toolbar={<Toolbar enableFilter />}
            onAddFilter={this.handleFilterChange}
            onClearFilters={this.onClearFilters}
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
