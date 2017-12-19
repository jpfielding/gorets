import React from 'react';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import PasswordForm from 'components/utils/PasswordForm';
import RouteLink from 'components/elements/RouteLink';

class ConnectionForm extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
    connection: React.PropTypes.any,
    updateConnection: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);

    const connectionForm = createValue({
      value: props.connection,
      onChange: this.connectionInputsChange.bind(this),
    });

    this.state = {
      connectionForm,
      show: false,
      element: null,
    };
  }

  connectionInputsChange(props) {
    const connectionForm = props;
    this.setState({ connectionForm });
  }

  render() {
    return (
      <div>
        <RouteLink connection={this.state.connectionForm.value} type={'basic'} style={{ marginBottom: '10px' }} />
        <Fieldset formValue={this.state.connectionForm}>
          <Field select="id" label="ID">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="loginURL" label="Login URL">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="username" label="Username">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="password" label="Password" Input={PasswordForm}>
            <Input
              className={`border-box w-100 pa1 b--none outline-transparent ${!this.state.show ? 'masker' : ''}`}
            />
          </Field>
          <Field select="userAgent" label="User Agent">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="userAgentPw" label="User Agent Password" Input={PasswordForm}>
            <Input
              className={`border-box w-100 pa1 b--none outline-transparent ${!this.state.show ? 'masker' : ''}`}
            />
          </Field>
          <Field select="proxy" label="Proxy">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="retsVersion" label="Version">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
        </Fieldset>
        <select className="customSelect" ref={(e) => { this.state.select = e; }}>
          <option value="COMPACT">COMPACT</option>
          <option value="COMPACT-INCREMENTAL">COMPACT-INCREMENTAL</option>
          <option value="STANDARD-XML">STANDARD-XML</option>
        </select>
        <div>
          <button
            className="customButton mt3"
            onClick={() => {
              const args = {
                extraction: this.state.select.value,
                oldest: 0,
              };
              this.props.updateConnection(this.state.connectionForm.value, args);
            }}
          >
            Update Changes
          </button>
          <button
            className="customButton mt3 fr"
            onClick={() => {
              const show = !this.state.show;
              this.setState({ show });
            }}
          >
            {this.state.show ? 'Hide Passwords' : 'Show Passwords'}
          </button>
        </div>
      </div>
    );
  }

}

export default ConnectionForm;
