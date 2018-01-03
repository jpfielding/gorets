import React from 'react';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import PasswordForm from 'components/utils/PasswordForm';
import RouteLink from 'components/elements/RouteLink';

class ConnectionForm extends React.Component {

  static propTypes = {
    location: React.PropTypes.any,
    connection: React.PropTypes.any,
    updateConnection: React.PropTypes.Func,
    idprefix: React.PropTypes.any,
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
        <RouteLink
          connection={this.state.connectionForm.value}
          type={'basic'} style={{ marginBottom: '10px' }}
          idprefix={`${this.props.idprefix}-link`}
        />
        <Fieldset formValue={this.state.connectionForm}>
          <Field select="id" label="ID">
            <Input
              className="border-box w-100 pa1 b--none outline-transparent"
              id={`${this.props.idprefix}-id`}
            />
          </Field>
          <Field select="loginURL" label="Login URL">
            <Input
              className="border-box w-100 pa1 b--none outline-transparent"
              id={`${this.props.idprefix}-loginurl`}
            />
          </Field>
          <Field select="username" label="Username">
            <Input
              className="border-box w-100 pa1 b--none outline-transparent"
              id={`${this.props.idprefix}-username`}
            />
          </Field>
          <Field select="password" label="Password" Input={PasswordForm}>
            <Input
              className={`border-box w-100 pa1 b--none outline-transparent ${!this.state.show ? 'masker' : ''}`}
              id={`${this.props.idprefix}-password`}
            />
          </Field>
          <Field select="userAgent" label="User Agent">
            <Input
              className="border-box w-100 pa1 b--none outline-transparent"
              id={`${this.props.idprefix}-useragent`}
            />
          </Field>
          <Field select="userAgentPw" label="User Agent Password" Input={PasswordForm}>
            <Input
              className={`border-box w-100 pa1 b--none outline-transparent ${!this.state.show ? 'masker' : ''}`}
              id={`${this.props.idprefix}-useragentpassword`}
            />
          </Field>
          <Field select="proxy" label="Proxy">
            <Input
              className="border-box w-100 pa1 b--none outline-transparent"
              id={`${this.props.idprefix}-proxy`}
            />
          </Field>
          <Field select="retsVersion" label="Version">
            <Input
              className="border-box w-100 pa1 b--none outline-transparent"
              id={`${this.props.idprefix}-version`}
            />
          </Field>
        </Fieldset>
        <select className="customSelect" ref={(e) => { this.state.select = e; }} id={`${this.props.idprefix}-type`}>
          <option value="COMPACT" id={`${this.props.idprefix}-type-compact`}>COMPACT</option>
          <option
            value="COMPACT-INCREMENTAL"
            id={`${this.props.idprefix}-type-incremantal`}
          >COMPACT-INCREMENTAL</option>
          <option value="STANDARD-XML" id={`${this.props.idprefix}-type-xml`}>STANDARD-XML</option>
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
            id={`${this.props.idprefix}-update`}
          >
            Update Changes
          </button>
          <button
            className="customButton mt3 fr"
            onClick={() => {
              const show = !this.state.show;
              this.setState({ show });
            }}
            id={`${this.props.idprefix}-pwdtoggle`}
          >
            {this.state.show ? 'Hide Passwords' : 'Show Passwords'}
          </button>
        </div>
      </div>
    );
  }

}

export default ConnectionForm;
