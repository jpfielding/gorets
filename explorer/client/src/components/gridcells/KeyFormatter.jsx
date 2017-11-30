import React from 'react';

export default class KeyFormatter extends React.Component {
  static propTypes = {
    value: React.PropTypes.any,
    metadataResource: React.PropTypes.any,
    metadataClass: React.PropTypes.any,
    displayContents: React.PropTypes.Func,
  };

  constructor(props) {
    super(props);
    this.state = {
      current: {},
      values: {},
    };
  }

  isLookup(value) {
    const selectedClass = this.props.metadataClass['METADATA-TABLE'].Field;
    const current = selectedClass.filter((e) => (value === e.SystemName));
    if (current.length !== 0 && current[0].LookupName) {
      const lookup = this.props.metadataResource['METADATA-LOOKUP'].Lookup;
      const values = lookup.filter((e) => (current[0].LookupName === e.LookupName));
      if (values.length !== 0) {
        this.state.values = values[0];
        this.state.current = current[0];
        return true;
      }
    }
    return false;
  }

  renderLookup(values) {
    return (
      <table>
        <tr>
          <th>
            Value
          </th>
          <th>
            Short Value
          </th>
          <th>
            Long Value
          </th>
        </tr>
        {values.map((e) => (
          <tr>
            <td>
              {e.Value}
            </td>
            <td>
              {e.ShortValue}
            </td>
            <td>
              {e.LongValue}
            </td>
          </tr>
          ))}
      </table>
    );
  }

  render() {
    const value = this.props.value;
    if (value === this.props.metadataResource.KeyField) {
      return (
        <div>
          <div className="customResultsButtonTitle" style={{ display: 'inline-block', marginRight: '5px' }}>
            Key
          </div>
          {value}
        </div>
      );
    } else if (this.isLookup(value)) {
      return (
        <div className="flex customSearchFormater">
          <div style={{ flex: '1' }}>
            {value}
          </div>
          <button
            className="fr"
            style={{ display: 'inline-block', marginRight: '5px' }}
            onClick={() => {
              this.props.displayContents(
                <div
                  style={{
                    position: 'absolute',
                    zIndex: '500',
                  }}
                  className="customResultsSet bg-mainbg"
                >
                  <div className="customResultsTitle" >
                    {this.state.current.LookupName}
                    <button
                      onClick={() => this.props.displayContents(null)}
                      className="fr customButton"
                    >
                      X
                    </button>
                  </div>
                  <div className="customResultsBody" >
                    {this.renderLookup(this.state.values['METADATA-LOOKUP_TYPE'].LookupType)}
                  </div>
                </div>
              );
            }}
          >
              L
          </button>
        </div>
      );
    }
    return (
      <div>
        {value}
      </div>
    );
  }

}
