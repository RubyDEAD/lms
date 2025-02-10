
import '../src/login.css'

function login() {
    return (
        <div className="container">
          <h1 className="container-text-title">Welcome !</h1>
    
          <label className="container-text">User ID</label>
          <input type="text" className="container-input" placeholder="User ID" />
    
          <label className="container-text">Password</label>
          <input type="password" className="container-input" placeholder="Password" />
    
          <div className="auth-link">
            <a href="/register">Register</a>
            <a href="/forgot-password">Forgot Password?</a>
          </div>
          <button className='button-login'>Login</button>
        </div>
      );
    };

export default login;
