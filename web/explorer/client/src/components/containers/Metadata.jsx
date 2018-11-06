import React from 'react';
import { withRouter } from 'react-router';
import ReactDataGrid from 'react-data-grid';
import _ from 'underscore';
import MetadataService from 'services/MetadataService';
import KeyFormatter from 'components/gridcells/KeyFormatter';

const ReactDataGridPlugins = require('react-data-grid-addons');

const Toolbar = ReactDataGridPlugins.Toolbar;

/*
  Responsible for the metadata tab of each class
  Sends row and class selection back up to have selections
    effect other tabs through the shared prop object
*/
class Metadata extends React.Component {

  static propTypes = {
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
    onRowsCleared: React.PropTypes.Func,
    onClassSelected: React.PropTypes.Func,
    idprefix: React.PropTypes.any,
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
    const possibleColumns = [];
    classRows.forEach(f => {
      Object.keys(f).forEach(key => {
        if (possibleColumns.includes(key)) {
          return;
        }
        possibleColumns.push(key);
      });
    });
    this.state = {
      filters: {},
      classRows: field,
      filteredRows: field,
      possibleColumns,
      displayContents: null,
      ref: null,
      sortDirection: 'NONE',
      sortColumn: '',
      selected: {},
    };
    this.handleGridSort = this.handleGridSort.bind(this);
    this.handleFilterChange = this.handleFilterChange.bind(this);

    this.onRowsSelected = this.onRowsSelected.bind(this);
    this.onRowsDeselected = this.onRowsDeselected.bind(this);
  }

  onRowsSelected(rows) {
    this.props.onRowsSelected(rows);
  }

  onRowsDeselected(rows) {
    this.props.onRowsDeselected(rows);
  }

  // Called when a class is selected
  onClassSelected(res, cls) {
    if (cls === this.props.shared.class) {
      return; // return if the class is the currently selected one.
    }
    this.props.onClassSelected(res, cls); // signal up changes
    const classRows = cls['METADATA-TABLE'].Field;
    const possibleColumns = [];
    classRows.forEach(field => {
      Object.keys(field).forEach(key => {
        if (possibleColumns.includes(key)) {
          return;
        }
        possibleColumns.push(key);
      });
    });
    this.setState({
      classRows,
      filteredRows: classRows,
      possibleColumns,
      sortColumn: '',
      sortDirection: 'NONE',
      filters: {},
    });
  }

  onClearFilters = () => {
    const filteredRows = this.sortRows(_.clone(this.state.classRows), this.state.sortColumn, this.state.sortDirection);
    this.setState({ filters: {}, filteredRows });
  }

  getObjectTypes() {
    const r = this.props.shared.resource;
    if (r == null || !r['METADATA-OBJECT'] || !r['METADATA-OBJECT']['Object']) {
      return [];
    }
    return r['METADATA-OBJECT']['Object'].map(o => o.ObjectType) || [];
  }

  setDisplay(displayContents) {
    this.setState({ displayContents });
  }

  handleGridSort(sortColumn, sortDirection) {
    const rows = this.sortRows(this.filterRows(this.state.classRows, this.state.filters), sortColumn, sortDirection);
    this.setState({ filteredRows: rows, sortColumn, sortDirection });
  }

  handleFilterChange(filter) {
    const newFilters = Object.assign({}, this.state.filters);
    if (filter.filterTerm) {
      newFilters[filter.column.key] = filter;
    } else {
      delete newFilters[filter.column.key];
    }
    const newRows = this.sortRows(
      this.filterRows(this.state.classRows, newFilters),
      this.state.sortColumn,
      this.state.sortDirection
    );
    if (newRows.length > 0) {
      this.setState({ filteredRows: newRows, filters: newFilters });
    }
  }

  filterRows(rows, filters) {
    const newRows = _.clone(rows);
    Object.keys(filters).forEach(filter => {
      const filterObj = filters[filter];
      if (filterObj.filterTerm) {
        for (let i = newRows.length - 1; i >= 0; i--) {
          const row = newRows[i];
          const val = row[filterObj.column.key] || '';
          const stringVal = String(val).toLowerCase();
          if (stringVal.indexOf(filterObj.filterTerm.toLowerCase()) === -1) {
            newRows.splice(i, 1);
          }
        }
      }
    });
    return newRows;
  }

  sortRows(rows, sortColumn, sortDirection) {
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
    return sortDirection === 'NONE'
      ? rows
      : rows.sort(comparer);
  }

  renderSelectedClassDescription(clazz) {
    return (
      <span title={clazz.Description}>
        {clazz.ClassName} - <span className="f6 moon-gray">{clazz.Description}</span>
      </span>
    );
  }

  render() {
    const { filteredRows, classRows } = this.state;
    const selectedResource = this.props.shared.resource;
    const selectedClass = this.props.shared.class;
    console.log('SHared', this.props.shared);
    const selected = this.state.selected;
    selected.class = selectedClass;
    selected.resource = selectedResource;
    let tableBody; // What will be rendered
    if (classRows) {
      const fieldSet = this.state.possibleColumns.map((name) => {
        if (name === 'SystemName') {
          return {
            key: name,
            name,
            resizable: true,
            width: 200,
            sortable: true,
            filterable: true,
            formatter: (
              <KeyFormatter
                metadataResource={selectedResource}
                metadataClass={selectedClass}
                selected={selected}
                displayContents={(e) => this.setDisplay(e)}
                container={this.state.ref}
              />
            ),
          };
        }
        return {
          key: name,
          name,
          resizable: true,
          width: 200,
          sortable: true,
          filterable: true,
        };
      });
      const rowGetter = (i) => filteredRows[i];
      tableBody = (
        <div key={selectedClass.ClassName}>
          {selectedResource
            ? (
              <span>
                <span className="b">{selectedResource.ResourceID} </span>
                {this.renderSelectedClassDescription(selectedClass)}
              </span>
            )
            : null
          }
          <div>
            {'Object Types - '}
            <span className="moon-gray">
              {this.getObjectTypes().join(', ')}
            </span>
          </div>
          <div>
            {this.state.displayContents}
          </div>
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
                keys: { rowKey: 'SystemName', values: this.props.shared.fields.map(r => r.row.SystemName) },
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
      <div
        ref={(ref) => { this.state.ref = ref; }}
        className="flex"
        style={{
          maxWidth: '1500px',
          margin: 'auto',
        }}
      >
        <div className="w-20-ns no-list-style pa3 pr0">
          <div className="customResultsCompactSet">
            <div className="customResultsTitle">
              <h1
                className="f5 nonclickable" title={system.SystemDescription}
                id={`${this.props.idprefix}-systemtitle`}
              >
                {system.SystemID}
              </h1>
              {system.SystemDescription}
            </div>
            <div className="customResultsBody overflow-x-scroll nowrap">
              <ul className="pa1 ma0 nonclickable">
                {system['METADATA-RESOURCE'].Resource.map((r) =>
                  <li className="mb2 nonclickable">
                    <div title={r.Description} className="b">{r.ResourceID}</div>
                    <ul className="pa0 pl3 mv1">
                      {r['METADATA-CLASS'].Class ? r['METADATA-CLASS'].Class.map((mClass) =>
                        <li
                          onClick={() => this.onClassSelected(r, mClass)}
                          className="clickable metadata-class-name"
                          id={`${this.props.idprefix}-${r.ResourceID}-${mClass.ClassName}`}
                        >
                          {this.renderSelectedClassDescription(mClass)}
                        </li>
                      ) : null }
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
                  <div id={`${this.props.idprefix}-body`}>
                    { tableBody }
                  </div>
                )
                : <h1 className="f4" id={`${this.props.idprefix}-default`}>Please select a class to explore</h1>
              }
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(Metadata);
