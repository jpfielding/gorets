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
          <Field select="url" label="URLS">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="username" label="Username">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="password" label="Password" Input={PasswordForm}>
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="userAgent" label="User Agent">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="userAgentPw" label="User Agent Password" Input={PasswordForm}>
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="proxy" label="Proxy">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
          <Field select="version" label="Version">
            <Input className="border-box w-100 pa1 b--none outline-transparent" />
          </Field>
        </Fieldset>
        <button
          className="customButton mt3"
          onClick={() => this.props.updateConnection(this.state.connectionForm.value)}
        >
          Update Changes
        </button>
      </div>
    );
  }

}

export default ConnectionForm;
