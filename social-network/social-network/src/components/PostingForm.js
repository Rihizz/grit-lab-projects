import React, { useState } from "react";
import { useAuth } from "../AuthContext";

const fetchFollowers = async (activeUser) => {
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };
  const response = await fetch(
    "http://localhost:6969/api/followers?id=" + activeUser,
    requestOptions
  );
  const data = await response.json();
  if (response.status === 200) {
    // console.log("followers fetched");
    // console.log(data.followers);
    return data.followers;
  } else {
    alert("Error fetching followers.");
  }
};

const PostingForm = ({ fetchPosts, setPosts }) => {
  const [privacy, setPrivacy] = useState("0");
  const [selectedFriends, setSelectedFriends] = useState([]);
  const [followers, setFollowers] = useState([]);
  const { userID } = useAuth();

  React.useEffect(() => {
    const getFollowers = async () => {
      const followers = await fetchFollowers(userID);
      setFollowers(followers);
    };
    getFollowers();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();

    // send post to database after validation
    const post = document.getElementById("post").value;
    if (post.length < 1 || post.length > 500) {
      alert("Post must be 1-500 characters long.");
      return;
    }
    //check that post does not consist of only spaces
    if (post.trim().length === 0) {
      alert("Post must contain at least one non-space character.");
      return;
    }
    const privacyInput = privacy;
    const picture = document.getElementById("picture");

    const formData = new FormData();
    formData.append("content", post);
    formData.append("privacy", privacyInput);
    formData.append("users_who_see", JSON.stringify(selectedFriends)); // Convert to JSON string

    if (picture.files[0]) {
      formData.append("picture", picture.files[0]);
    }

    if (privacy === "onlyme" && selectedFriends.length > 0) {
      formData.append("friends", selectedFriends.join(","));
    }

    // console.log("selectedFriends", selectedFriends);

    if (privacy === "onlyme" && selectedFriends.length > 0) {
      formData.append("friends", selectedFriends.join(","));
    }

    const requestOptions = {
      method: "POST",
      // headers object is removed since the browser will set the correct content type and boundary for FormData
      body: formData,
      credentials: "include",
    };
    const response = await fetch(
      "http://localhost:6969/api/posting",
      requestOptions
    );
    const data = await response.json();
    // console.log(data);
    if (data.status === 200) {
      // clear form after posting successfully
      document.getElementById("post").value = "";
      document.getElementById("file-name").textContent = "";
      document.getElementById("picture").value = null;
      // Set privacy to "Public" (value = "0")
      setPrivacy("0");
      // fetch posts again
      setPosts(await fetchPosts());
    } else {
      alert("Error posting, check that the length is at least 10 characters.");
    }
  };

  const handlePrivacyChange = (e) => {
    setPrivacy(e.target.value);
    setSelectedFriends([]);
  };

  const handleFriendChange = (e) => {
    const selectedOptions = Array.from(e.target.selectedOptions).map(
      (option) => option.value
    );
    setSelectedFriends(selectedOptions);
  };

  const handleFileChange = (event) => {
    const fileInput = event.target;
    const fileNameSpan = document.getElementById("file-name");
  
    if (fileInput.files.length > 0) {
      const fileName = fileInput.files[0].name;
      fileNameSpan.textContent = fileName;
    } else {
      fileNameSpan.textContent = "";
    }
  };
  

  return (
    <div className="posting-form">
      <textarea
        className="post-box"
        type="text"
        rows="10"
        placeholder="What's on your mind?"
        id="post"
        required
        maxLength="500"
        minLength="1"
        title="Post should be 10-500 characters."
      />
      <div className="postoptions">
        <label className="upload" htmlFor="picture">
          Upload image/gif: <i className="fa-solid fa-folder-open fa-xl"></i>
          <input
          type="file"
          id="picture"
          accept="image/*"
          className="file-input"
          onChange={(e) => handleFileChange(e)}
          />
  <span id="file-name"></span>
        </label>
        <input
          type="file"
          id="picture"
          accept="image/*"
          style={{ display: "none" }}
        />

        <select id="privacy" value={privacy} onChange={handlePrivacyChange}>
          <option value="0">Public</option>
          <option value="1">Followers</option>
          <option value="2">Choose friends</option>
        </select>
        {privacy === "2" && followers != null && followers.length > 0 && (
          <select
            id="friends-dropdown"
            multiple
            value={selectedFriends}
            onChange={handleFriendChange}
          >
            {/* map of the followers */}
            {followers.map((follower) => (
              <option key={follower.id} value={follower.id}>
                {follower.full_name}
              </option>
            ))}
          </select>
        )}
        <button className="submit-post" onClick={handleSubmit}>
          Post
        </button>
      </div>
    </div>
  );
};

export default PostingForm;
