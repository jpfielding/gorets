import React from 'react';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import PasswordForm from 'components/utils/PasswordForm';

class ConnectionForm extends React.Component {

  static propTypes = {
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
    };
  }

  connectionInputsChange(props) {
    const connectionForm = props;
    this.setState({ connectionForm });
  }

  render() {
    return (
      <div>
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
              type={!this.state.show ? 'password' : ''}
              className="border-box w-100 pa1 b--none outline-transparent"
            />
          </Field>
          <Field select="userAgent" label="User Agent">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="userAgentPw" label="User Agent Password" Input={PasswordForm}>
            <Input
              type={!this.state.show ? 'password' : ''}
              className="border-box w-100 pa1 b--none outline-transparent"
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
        <button
          className="customButton mt3"
          onClick={() => this.props.updateConnection(this.state.connectionForm.value, this.state.select.value)}
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
    );
  }

}

export default ConnectionForm;
