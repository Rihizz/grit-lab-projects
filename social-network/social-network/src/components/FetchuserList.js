import { useState, useEffect } from "react";
import { useAuth } from "../AuthContext";

function useFetchUserList() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  //use AuthContext to get userID
  const { userID } = useAuth();

  useEffect(() => {
    //include userID in the request (FormData) in the body of the request
    const formData = new FormData();
    formData.append("userID", userID);
    // console.log(userID, "userID from useFetchUserList");
    async function fetchData() {
      const requestOption = {
        method: "POST",

        body: formData,
        credentials: "include", // send the cookie along with the request
      };

      try {
        const response = await fetch(
          "http://localhost:6969/api/userlist",
          requestOption
        );
        const data = await response.json();

        if (response.status !== 200) {
          throw Error(data.message);
        } else {
          setData(data);
          setLoading(false);
        }
      } catch (error) {
        setError(error);
        setLoading(false);
      }
    }

    fetchData();
  }, [userID]);

  return { data, loading, error };
}

export default useFetchUserList;
