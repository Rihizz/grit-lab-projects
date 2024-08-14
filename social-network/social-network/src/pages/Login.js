import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";
import "../css/Login.css";

function Login() {
  const navigate = useNavigate();
  const { setLoggedIn, checkLoginStatus } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const [registerConfirmPassword, setRegisterConfirmPassword] = useState("");
  const [registerNickname, setRegisterNickname] = useState("");
  const [registerFirstName, setRegisterFirstName] = useState("");
  const [registerLastName, setRegisterLastName] = useState("");
  const [registerBirthday, setRegisterBirthday] = useState("");
  const [registerAboutMe, setRegisterAboutMe] = useState("");
  const [registerProfilePicture, setRegisterProfilePicture] = useState("");
  const [registrationSuccess, setRegistrationSuccess] = useState(false);
  const [registerError, setRegisterError] = useState(false);
  const [invalidDateError, setInvalidDateError] = useState("");
  const [regprivacy, setRegPrivacy] = useState("0");
  function navigateToHomePage() {
    navigate("/");
  }
  async function handleLoginSubmit(event) {
    event.preventDefault();
    // Call backend API to log in user with email and password
    const response = await fetch("http://localhost:6969/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });
    if (response.ok) {
      // Handle successful login
      // Call checkLoginStatus function from AuthContext to fetch the login status
      await checkLoginStatus(); // Make sure to import it from your AuthContext
      setLoggedIn(true);
      navigateToHomePage();
    } else {
      // Handle unsuccessful login
      alert("Invalid email or password. Please try again.");
    }
  }
  async function handleRegisterSubmit(event) {
    event.preventDefault();
    setRegisterError(false);
    setRegistrationSuccess(false);
    if (registerPassword !== registerConfirmPassword) {
      // Handle case where password and confirm password do not match
      alert("Passwords do not match. Please try again.");
      return;
    }
    // Create a FormData object and append all the form fields
    const formData = new FormData();
    formData.append("email", registerEmail);
    formData.append("password", registerPassword);
    formData.append("nickname", registerNickname);
    formData.append("firstName", registerFirstName);
    formData.append("lastName", registerLastName);
    formData.append("birthday", registerBirthday);
    formData.append("aboutMe", registerAboutMe);
    formData.append("profilePicture", registerProfilePicture);
    formData.append("privacy", regprivacy);
    const response = await fetch("http://localhost:6969/api/register", {
      method: "POST",
      headers: {},
      body: formData,
      credentials: "include",
    });
    if (response.ok) {
      // Handle successful registration
      setRegistrationSuccess(true);
      setEmail(registerEmail);
      setPassword(registerPassword);
      setRegisterAboutMe("");
      setRegisterBirthday("");
      setRegisterConfirmPassword("");
      setRegisterEmail("");
      setRegisterFirstName("");
      setRegisterLastName("");
      setRegisterNickname("");
      setRegisterPassword("");
      setRegisterProfilePicture(null);
      setInvalidDateError(false);
      setRegPrivacy("0");
    } else if (response.status === 400) {
      const text = await response.text();
      const cleaned = text.slice(0, -1);
      setInvalidDateError(cleaned);
    } else {
      // Handle unsuccessful registration
      setRegisterError(true);
    }
  }
  return (
    <div>
      <div className="login">
        <h2>Login</h2>
        <form onSubmit={handleLoginSubmit}>
          <div className="logemail">
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              autoComplete="username"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
            />
          </div>
          <div className="logpassword">
            <label htmlFor="password">Password:</label>
            <input
              type="password"
              id="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
          </div>
          <button type="submit">Login</button>
        </form>
      </div>
      <div className="register">
        <h2>Register</h2>
        {invalidDateError !== "" && <p id="invalidDate">{invalidDateError}</p>}
        <form onSubmit={handleRegisterSubmit}>
          <div className="register-left">
            <div className="regemail">
              <label htmlFor="registerEmail">Email:</label>
              <input
                type="email"
                id="registerEmail"
                autoComplete="username"
                value={registerEmail}
                onChange={(event) => setRegisterEmail(event.target.value)}
                required
                // pattern="^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
                title="Enter a valid email"
              />
            </div>
            <div className="regpassword">
              <label htmlFor="registerPassword">Password:</label>
              <input
                type="password"
                id="registerPassword"
                autoComplete="new-password"
                value={registerPassword}
                onChange={(event) => setRegisterPassword(event.target.value)}
                required
                // pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$"
                title="Password must contain at least 8 characters, 1 uppercase, 1 lowercase, 1 number, 1 special character"
                maxLength="50"
              />
            </div>
            <div className="regconfirmpass">
              <label htmlFor="registerConfirmPassword">Confirm Password:</label>
              <input
                type="password"
                id="registerConfirmPassword"
                autoComplete="new-password"
                value={registerConfirmPassword}
                onChange={(event) =>
                  setRegisterConfirmPassword(event.target.value)
                }
                required
                maxLength="50"
              />
            </div>
            <div className="regnickname">
              <label htmlFor="registerNickname">Nickname:</label>
              <input
                type="text"
                id="registerNickname"
                value={registerNickname}
                onChange={(event) => setRegisterNickname(event.target.value)}
                //    pattern="^[a-zA-Z\u00C0-\u017F]{3,20}$"
                title="Nickname must contain at least 3 characters, max 20 characters, no special characters"
                maxLength="20"
              />
            </div>
          </div>
          <div className="register-right">
            <div className="regfirstname">
              <label htmlFor="registerFirstName">First Name:</label>
              <input
                type="text"
                id="registerFirstName"
                value={registerFirstName}
                onChange={(event) => setRegisterFirstName(event.target.value)}
                required
                //     pattern="^[a-zA-Z\u00C0-\u017F]{3,20}$"
                title="First name must contain at least 3 characters, max 20 characters, no special characters"
                maxLength="20"
              />
            </div>
            <div className="reglastname">
              <label htmlFor="registerLastName">Last Name:</label>
              <input
                type="text"
                id="registerLastName"
                value={registerLastName}
                onChange={(event) => setRegisterLastName(event.target.value)}
                required
                //      pattern="^[a-zA-Z\u00C0-\u017F]{3,20}$"
                title="Last name must contain at least 3 characters, max 20 characters, no special characters"
                maxLength="20"
              />
            </div>
            <div className="regbirthday">
              <label htmlFor="registerBirthday">Birthday (DD/MM/YYYY):</label>
              <input
                type="text"
                id="registerBirthday"
                value={registerBirthday}
                onChange={(event) => setRegisterBirthday(event.target.value)}
                required
                //   pattern="^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$"
                title="Date must be in the format DD/MM/YYYY for example 01/01/2000"
              />
            </div>
            <div className="regaboutme">
              <label htmlFor="registerAboutMe">About Me:</label>
              <input
                type="text"
                id="registerAboutMe"
                value={registerAboutMe}
                onChange={(event) => setRegisterAboutMe(event.target.value)}
                maxLength="500"
                minLength="1"
                title="Write at least something, but max 500 characters"

                //pattern for 1-500 characters
                //  pattern="^.{1,500}$"
              />
            </div>
            <div className="regprofilepic">
              <label htmlFor="registerProfilePicture">Profile Picture:</label>
              <input
                type="file"
                accept="image/*"
                id="registerProfilePicture"
                onChange={(event) =>
                  setRegisterProfilePicture(event.target.files[0])
                }
              />
            </div>
            <div className="regprivacy">
              <label htmlFor="registerPrivacy">Privacy:</label>
              <select
                id="registerPrivacy"
                value={regprivacy}
                onChange={(event) => setRegPrivacy(event.target.value)}
              >
                <option value="0">Public</option>
                <option value="1">Private</option>
              </select>
            </div>
          </div>
          <button type="submit">Register</button>
          {registrationSuccess && <p>Registration successful!</p>}
          {registerError && (
            <p>Registration failed. Please try again with a different email.</p>
          )}
        </form>
      </div>
    </div>
  );
}

export default Login;
