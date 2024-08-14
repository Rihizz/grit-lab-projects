import React, { useState, useEffect } from "react";
import "../styles/ProfileCard.css";

function EditProfileModal(props) {
  const {
    show,
    handleClose,
    handleSave,
    userId,
    currentUserData,
    errorMessage,
  } = props;
  const [nickname, setNickname] = useState(currentUserData.nickname);
  const [email, setEmail] = useState(currentUserData.email);
  const [aboutMe, setAboutMe] = useState(currentUserData.aboutMe);
  const [validationMessage, setValidationMessage] = useState("");
  const [avatar, setAvatar] = useState(
    currentUserData.avatar
      ? `data:image/jpeg;base64,${currentUserData.avatar}`
      : null
  );
  const [avatarFile, setAvatarFile] = useState(null); // New state for storing the File object of the avatar
  const [privacy, setPrivacy] = useState(currentUserData.privacy || "0");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [privacyChanged, setPrivacyChanged] = useState(false); //stat to track privacy change

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    setAvatarFile(file);

    // Read the file and convert it to a base64 string
    const reader = new FileReader();
    reader.onloadend = function () {
      setAvatar(reader.result);
    };
    reader.readAsDataURL(file);
  };

  const handleSubmit = () => {
    // Check if the password or confirmPassword field is not empty
    if (newPassword || confirmPassword) {
      if (newPassword !== confirmPassword) {
        setValidationMessage("The passwords do not match.");
        return;
      }
    }
    //check that nickname and about me do not consist of only spaces
    if (nickname) {
      if (nickname.trim().length === 0) {
        setValidationMessage(
          "Nickname must contain at least one non-space character."
        );
        return;
      }
    }
    if (aboutMe) {
      if (aboutMe.trim().length === 0) {
        setValidationMessage("About Me can't be just whitespace.");
        return;
      }
    }
    //check that email is not space @ space, only using trim wont work because there would still be the @ sign
    if (email) {
      if (email.trim().length === 0 || email.trim().length < 3) {
        setValidationMessage("Email not valid.");
        return;
      }
    }

    // setValidationMessage(""); // Clear the validation message if there are no issues

    // console.log(privacy);
    handleSave({
      userId,
      nickname,
      email,
      aboutMe,
      avatar: avatarFile, // Pass the File object of the avatar, not the base64 string
      newPassword,
      confirmPassword,
      privacy,
    });
  };

  return (
    <div
      className={`modal ${show ? "show" : ""}`}
      style={{ display: show ? "block" : "none" }}
    >
      <div className="profile-edit-box">
        <div className="modal-header">
          <h5 className="title">Edit Profile</h5>
        </div>
        <div className="content-container">
          <div className="modal-body">
            <div className="newavatar">
              <label htmlFor="avatar" className="form-avatar">
                New Avatar
              </label>
            </div>

            <input
              type="file"
              accept="image/*"
              className="avatar-box"
              id="avatar"
              onChange={handleFileChange}
            />

            <div className="newnickname">
              <label htmlFor="nickname" className="form-nickname">
                New Nickname
              </label>
            </div>

            <input
              type="text"
              className="nickname-box"
              id="nickname"
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              required
              maxLength="20"
              
            />

            <div className="newemail">
              <label htmlFor="email" className="form-email">
                New Email
              </label>
            </div>

            <input
              type="email"
              className="email-box"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />

            <div className="newpassword">
              <label htmlFor="newpassword" className="form-newpassword">
                New Password
              </label>
            </div>

            <input
              type="password"
              className="newpassword-box"
              id="newpassword"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
            />

            <div className="confirmpassword">
              <label htmlFor="confirmpassword" className="confirmpassword-box">
                Confirm Password
              </label>
            </div>

            <input
              type="password"
              className="confirmpassword-box"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              maxLength="50"
            />

            <div className="newabout">
              <label htmlFor="about-me" className="form-about">
                New 'About Me' Text
              </label>
            </div>

            <textarea
              className="about-textarea"
              id="about-me"
              rows="5"
              value={aboutMe}
              onChange={(e) => setAboutMe(e.target.value)}
              title="Write at least something, but max 500 characters"
              required
            ></textarea>

            <div className="newprivacy">
              <label htmlFor="privacy" className="form-privacy">
                Privacy
              </label>
            </div>
            <select
              id="privacy"
              value={privacy}
              onChange={(e) => setPrivacy(e.target.value)}
            >
              <option value="0">Public</option>
              <option value="1">Private</option>
            </select>
          </div>
        </div>
        <div className="modal-footer">
          {validationMessage && (
            <div className="validation-message">{validationMessage}</div>
          )}
          {errorMessage && <div className="error-message">{errorMessage}</div>}
          <button type="button" className="btn-close" onClick={handleClose}>
            Close
          </button>
          <button type="button" className="btn-save" onClick={handleSubmit}>
            Save Changes
          </button>
        </div>
      </div>
    </div>
  );
}

export default EditProfileModal;
