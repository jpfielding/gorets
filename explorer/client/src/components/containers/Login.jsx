import React from 'react';
import { Fieldset, Field, createValue } from 'react-forms';
import LoginService from 'services/LoginService';
import PasswordForm from 'components/utils/PasswordForm';

export default class Login extends React.Component {

  constructor(props) {
    super(props);
    const formValue = createValue({
      value: null,
      onChange: this.onChange.bind(this),
    });
    this.state = { formValue };
  }

  onChange(formValue) {
    this.setState({ formValue });
  }

  onSubmit = () => {
    console.log('Submitting: ', this.state.formValue.value);
    LoginService
      .login(this.state.formValue.value)
      .then((response) => {
        console.log(response);
      });
  }

  render() {
    return (
      <div className="flex flex-row justify-center">
        <div className="w-60">
          <h1>Login</h1>
          <Fieldset formValue={this.state.formValue}>
            <Field select="url" label="Login URL" />
            <Field select="username" label="Username" />
            <Field select="password" label="Password" Input={PasswordForm} />
            <Field select="user-agent" label="User Agent" />
            <Field select="user-agent-pw" label="User Agent Password" Input={PasswordForm} />
            <Field select="rets-version" label="Protocol Version" />
            <Field select="id" label="Custom RETs Name" />
            <button onClick={this.onSubmit}>Submit</button>
          </Fieldset>
        </div>
      </div>
    );
  }
}
