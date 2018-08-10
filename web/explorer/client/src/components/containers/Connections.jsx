import React from 'react';
import { Fieldset, Field, createValue, Input } from 'react-forms';

export default class Connections extends React.Component {

  static propTypes = {
    addTab: React.PropTypes.Func,
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
    console.log('Launching tab with: ', this.state.formValue.value);
    this.props.addTab(this.state.formValue.value);
  }

  render() {
    return (
      <div
        className="customResultsSet center mt3"
        style={{
          maxWidth: '400px',
        }}
      >
        <div className="customResultsTitle">
          <div className="customTitle tc">
            Add a connection
          </div>
        </div>
        <div className="customResultsBody flex flex-row justify-center w-100">
          <Fieldset formValue={this.state.formValue} className="w-100">
            <Field select="id" label="ID (unique per config service)" >
              <Input className="border-box w-100 pa1 b--none outline-transparent" id="newcon-id" />
            </Field>
            <Field select="loginURL" label="Login URL">
              <Input className="border-box w-100 pa1 b--none outline-transparent" id="newcon-login" />
            </Field>
            <Field select="username" label="Username" >
              <Input className="border-box w-100 pa1 b--none outline-transparent" id="newcon-username" />
            </Field>
            <Field select="password" label="Password" >
              <Input className="border-box w-100 pa1 b--none outline-transparent masker" id="newcon-password" />
            </Field>
            <Field select="userAgent" label="User Agent" >
              <Input className="border-box w-100 pa1 b--none outline-transparent" id="newcon-useragent" />
            </Field>
            <Field select="userAgentPw" label="User Agent Password" >
              <Input
                className="border-box w-100 pa1 b--none outline-transparent masker"
                id="newcon-useragentpassword"
              />
            </Field>
            <Field select="retsVersion" label="Protocol Version" >
              <Input className="border-box w-100 pa1 b--none outline-transparent" id="newcon-version" />
            </Field>
            <Field select="proxy" label="Proxy (Socks5)" >
              <Input className="border-box w-100 pa1 b--none outline-transparent" id="newcon-proxy" />
            </Field>
          </Fieldset>
        </div>
        <div className="customResultsFoot">
          <button
            onClick={this.onSubmit}
            className="customButton db"
            style={{ margin: 'auto', width: '50%' }}
            id="newcon-submit"
          >
            Submit
          </button>
        </div>
      </div>
    );
  }
}
