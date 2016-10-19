import React from 'react';
import { withFormValue } from 'react-forms';

/* eslint-disable jsx-a11y/label-has-for */
class Field extends React.Component {

  static propTypes = {
    formValue: React.PropTypes.any,
  }

  onChange = (e) => this.props.formValue.update(e.target.value)

  render() {
    const { formValue } = this.props;
    return (
      <input type="password" value={formValue.value} onChange={this.onChange} />
    );
  }
}

export default withFormValue(Field);
