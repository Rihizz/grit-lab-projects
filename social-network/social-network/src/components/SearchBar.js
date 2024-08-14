import React, { useRef, useState } from "react";
import { Link } from "react-router-dom";

// send the search query to the server

const sendSearchQuery = async (searchQuery) => {
  if (searchQuery === "") {
    return;
  }
  const response = await fetch("http://localhost:6969/api/search-users", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      searchQuery: searchQuery,
    }),
    credentials: "include",
  });

  // console.log("search query sent");

  // handle response
  const data = await response.json();
  if (response.status === 200) {
    // console.log("search query received");
    // console.log(data);
  } else {
    alert("Error searching for users.");
  }

  let users = data;
  // console.log(users);

  // display the search results
  const searchResults = document.getElementById("search-results");
  searchResults.innerHTML = "";
  if (users) {
    users.forEach((user) => {
      // console.log(user.full_name);
      const userElement = document.createElement("search-user");
      userElement.innerHTML = `<a href="/profile/${user.id}">${user.full_name}</a>`;
      searchResults.appendChild(userElement);
    });
  }
};

const SearchBar = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [isSearchBarVisible, setSearchBarVisible] = useState(false);
  const serchInPutRef = useRef(null); // Create a reference to the search input element
  const handleInputChange = (event) => {
    setSearchQuery(event.target.value);
    sendSearchQuery(event.target.value);
    if (event.target.value === "") {
      const searchResults = document.getElementById("search-results");
      searchResults.innerHTML = "";
    }
  };

  const modalRef = React.useRef();

  // Close the search bar when the user clicks outside of it
  React.useEffect(() => {
    const handleOutsideClick = (event) => {
      if (modalRef.current && !modalRef.current.contains(event.target)) {
        setSearchBarVisible(false);
      }
    };

    document.addEventListener("mousedown", handleOutsideClick);

    // Clean up the event listener
    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
    };
  }, []);

  const toggleSearchBar = () => {
    setSearchBarVisible(!isSearchBarVisible);
    if (!isSearchBarVisible) {
      serchInPutRef.current.focus(); // Set focus on the search input when the bar is clicked
    } else {
      serchInPutRef.current.blur(); // Remove focus from the search input when the bar is clicked again
    }
  };

  return (
    <div className="searchbark" ref={modalRef}>
      <i
        className="fa-solid fa-magnifying-glass fa-lg"
        onClick={toggleSearchBar}
      ></i>

      <input
        type="text"
        placeholder="Search users..."
        value={searchQuery}
        onChange={handleInputChange}
        className={`search-input ${isSearchBarVisible ? "visible" : ""}`}
        ref={serchInPutRef}
      />

      <div
        className={`result ${isSearchBarVisible ? "visible" : ""}`}
        id="search-results"
      ></div>
    </div>
  );
};

export default SearchBar;
