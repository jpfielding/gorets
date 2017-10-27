import React from 'react';
import { withRouter } from 'react-router';
import ReactDataGrid from 'react-data-grid';
import MetadataService from 'services/MetadataService';

const ReactDataGridPlugins = require('react-data-grid/addons');

const Toolbar = ReactDataGridPlugins.Toolbar;

class Metadata extends React.Component {

  static propTypes = {
    params: React.PropTypes.any,
    location: React.PropTypes.any,
    router: React.PropTypes.any,
    shared: {
      connection: React.PropTypes.any,
      metadata: React.PropTypes.any,
      resource: React.PropTypes.any,
      class: React.PropTypes.any,
      fields: React.PropTypes.any,
      rows: React.PropTypes.any,
    },
    onRowsSelected: React.PropTypes.Func,
    onRowsDeselected: React.PropTypes.Func,
    onClassSelected: React.PropTypes.Func,
  }

  static defaultProps = {
    connection: { id: 'n/a' },
    metadata: MetadataService.empty,
  }

  constructor(props) {
    super(props);

    const mtable = props.shared.class['METADATA-TABLE'];
    const field = (mtable) ? mtable.Field : [];
    const classRows = (mtable) ? mtable.Field : [];
    const selectedFieldSet = [];
    classRows.forEach(f => {
      Object.keys(f).forEach(key => {
        if (selectedFieldSet.includes(key)) {
          return;
        }
        selectedFieldSet.push(key);
      });
    });
    this.state = {
      filters: {},
      classRows: field,
      filteredRows: field,
      selectedFieldSet,
    };
    this.handleGridSort = this.handleGridSort.bind(this);
    this.onRowsSelected = this.onRowsSelected.bind(this);
    this.onRowsDeselected = this.onRowsDeselected.bind(this);
  }

  onRowsSelected(rows) {
    this.props.onRowsSelected(rows);
  }

  onRowsDeselected(rows) {
    this.props.onRowsDeselected(rows);
  }

  onClassSelected(res, cls) {
    if (cls === this.props.shared.class) {
      return;
    }
    this.props.onClassSelected(res, cls);
    const classRows = cls['METADATA-TABLE'].Field;
    const filteredRows = cls['METADATA-TABLE'].Field;
    const selectedFieldSet = [];
    classRows.forEach(field => {
      Object.keys(field).forEach(key => {
        if (selectedFieldSet.includes(key)) {
          return;
        }
        selectedFieldSet.push(key);
      });
    });
    this.setState({
      classRows,
      filteredRows,
      selectedFieldSet,
    });
  }

  onClearFilters = () => {
    this.setState({ filters: {} });
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
      ? this.state.filteredRows
      : this.state.filteredRows.sort(comparer);
    this.setState({ filteredRows: rows });
  }

  // TODO row selection and filters arent playing well together
  handleFilterChange = (filter) => {
    const newFilters = Object.assign({}, this.state.filters);
    if (filter.filterTerm) {
      newFilters[filter.column.key] = filter;
    } else {
      delete newFilters[filter.column.key];
    }
    this.setState({ filters: newFilters });
    const rows = this.state.classRows;
    const newRows = [...rows];
    Object.keys(newFilters).forEach(newFilter => {
      const filterObj = newFilters[newFilter];
      if (filterObj.filterTerm) {
        console.log(filterObj.filterTerm);
        console.log('rowLen: ', rows.length);
        for (let i = rows.length - 1; i >= 0; i--) {
          const row = rows[i];
          const val = row[filterObj.column.key] || '';
          const stringVal = String(val).toLowerCase();
          if (stringVal.indexOf(filterObj.filterTerm.toLowerCase()) === -1) {
            newRows.splice(i, 1);
          }
        }
      }
    });
    if (newRows.length > 0) {
      this.setState({ filteredRows: newRows });
    }
  }

  renderSelectedClassDescription(clazz) {
    return (
      <span title={clazz.Description}>
        {clazz.ClassName} - <span className="f6 moon-gray">{clazz.Description}</span>
      </span>
    );
  }

  render() {
    const { filteredRows } = this.state;
    let tableBody;
    if (filteredRows) {
      const availableFields = this.state.selectedFieldSet;
      const fieldSet = (filteredRows && filteredRows.length > 0)
        ? availableFields.map((name) => ({
          key: name,
          name,
          resizable: true,
          width: 200,
          sortable: true,
          filterable: true,
        }))
        : [];
      const rowGetter = (i) => filteredRows[i];
      const selectedResource = this.props.shared.resource;
      const selectedClass = this.props.shared.class;
      tableBody = (
        <div>
          {selectedResource
            ? (
              <span>
                <span className="b">{selectedResource.ResourceID} </span>
                {this.renderSelectedClassDescription(selectedClass)}
              </span>
            )
            : null
          }
          <ReactDataGrid
            onGridSort={this.handleGridSort}
            columns={fieldSet}
            rowGetter={rowGetter}
            rowsCount={filteredRows.length}
            toolbar={<Toolbar enableFilter />}
            onAddFilter={this.handleFilterChange}
            onClearFilters={this.onClearFilters}
            minHeight={500}
            rowSelection={{
              showCheckbox: true,
              enableShiftSelect: true,
              onRowsSelected: this.onRowsSelected,
              onRowsDeselected: this.onRowsDeselected,
              selectBy: {
                // TODO this does not deal with filter/sort
                indexes: this.props.shared.fields.map(r => r.rowIdx),
              },
            }}
          />
        </div>
      );
    } else {
      tableBody = null;
    }
    const system = this.props.shared.metadata.System;
    return (
      <div className="flex">
        <div className="w-20-ns no-list-style pa3 pr0">
          <div className="customResultsCompactSet">
            <div className="customResultsTitle">
              <h1 className="f5 nonclickable" title={system.SystemDescription}>
                {system.SystemID}
              </h1>
            </div>
            <div className="customResultsBody overflow-x-scroll nowrap">
              <ul className="pa1 ma0 nonclickable">
                {system['METADATA-RESOURCE'].Resource.map((r) =>
                  <li className="mb2 nonclickable">
                    <div title={r.Description} className="b">{r.ResourceID}</div>
                    <ul className="pa0 pl3 mv1">
                      {r['METADATA-CLASS'].Class.map((mClass) =>
                        <li
                          onClick={() => this.onClassSelected(r, mClass)}
                          className="clickable metadata-class-name"
                        >
                          {this.renderSelectedClassDescription(mClass)}
                        </li>
                      )}
                    </ul>
                  </li>
                )}
              </ul>
            </div>
          </div>
        </div>
        <div className="fl w-100 w-80-ns pa3">
          <div className="customResultsSet">
            <div className="customResultsTitle">
              <div className="customTitle">
                  Metadata:
              </div>
            </div>
            <div className="customResultsBody">
              { this.state.classRows.length > 0
                ? (
                  <div>
                    { tableBody }
                  </div>
                )
                : <h1 className="f4">Please select a class to explore</h1>
              }
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(Metadata);
