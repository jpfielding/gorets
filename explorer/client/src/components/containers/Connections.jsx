import React from 'react';
import { Fieldset, Field, createValue } from 'react-forms';
import ConnectionService from 'services/ConnectionService';
import PasswordForm from 'components/utils/PasswordForm';

export default class Connections extends React.Component {

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
    ConnectionService
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
        <h1 className="f4 tc">Add a connection</h1>
        <div className="flex flex-row justify-center">
          <div
            style={{
              border: '1px solid black',
              padding: '20px',
            }}
          >
            <Fieldset formValue={this.state.formValue}>
              <Field select="url" label="Login URL" />
              <Field select="username" label="Username" />
              <Field select="password" label="Password" Input={PasswordForm} />
              <Field select="user-agent" label="User Agent" />
              <Field select="user-agent-pw" label="User Agent Password" Input={PasswordForm} />
              <Field select="rets-version" label="Protocol Version" />
              <Field select="id" label="Custom RETs Name" />
              <button onClick={this.onSubmit} className="ma2 ba black bg-transparent b--black">Submit</button>
            </Fieldset>
          </div>
        </div>
      </div>
    );
  }
}
