import React from 'react';

class Login extends React.Component {
  render() {
    return ( <div className="login">Login</div> )
  }
}

Login.propTypes = {
  apiUrl = React.propTypes.string.isRequired
}

export default Login;
