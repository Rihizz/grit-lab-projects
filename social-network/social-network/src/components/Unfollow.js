import React, { useEffect } from "react";
import { useAuth } from "../AuthContext";

function Unfollow({ userId }) {
    const { userID } = useAuth();
    
    const handleUnfollow = async () => {
        const requestOptions = {
        method: "POST",
    
        headers: {
            "Content-Type": "application/json",
        },
    
        body: JSON.stringify({
            id: userId,
        }),
    
        credentials: "include",
        };
    
        const response = await fetch(
        "http://localhost:6969/api/unfollow",
        requestOptions
        );
    
        if (response.ok) {
        // console.log("unfollowed");
        } else {
        // console.log("unfollow failed");
        }
    };
    
    useEffect(() => {
        handleUnfollow();
    }, [userId, userID]);
    
    return <></>;
}

export default Unfollow;