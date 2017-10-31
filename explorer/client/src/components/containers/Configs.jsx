import React from 'react';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import ConfigService from 'services/ConfigService';
import PasswordForm from 'components/utils/PasswordForm';

export default class Configs extends React.Component {

  static propTypes = {
    updateCallback: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);
    const formValue = createValue({
      value: null,
      onChange: this.onChange.bind(this),
    });
    this.state = {
      formValue,
    };
    this.onSubmit = this.onSubmit.bind(this);
  }

  onChange(formValue) {
    this.setState({ formValue });
  }

  onSubmit = () => {
    console.log('Submitting: ', this.state.formValue.value);
    ConfigService
      .login(this.state.formValue.value)
      .then(response => {
        console.log(response);
        this.props.updateCallback();
      }).catch((err) => {
        console.error(err);
      });
  }

  render() {
    return (
      <div className="pa2">
        <h1 className="f4 tc nonclickable">Add a Config</h1>
        <div className="flex flex-row justify-center">
          <div className="b--solid w5 ba pa2 tc">
            <Fieldset formValue={this.state.formValue}>
              <Field select="url" label="Login URL">
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="username" label="Username" >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="password" label="Password" Input={PasswordForm} >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="userAgent" label="User Agent" >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="userAgentPw" label="User Agent Password" Input={PasswordForm} >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="version" label="Protocol Version" >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="proxy" label="Proxy (Socks5)" >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <Field select="id" label="Custom RETs Name" >
                <Input className="pa1 b--none outline-transparent" />
              </Field>
              <button
                onClick={this.onSubmit}
                className="ma2 ba black bg-transparent b--black outline-transparent rd-focus clickable"
              >
                Submit
              </button>
            </Fieldset>
          </div>
        </div>
      </div>
    );
  }
}
