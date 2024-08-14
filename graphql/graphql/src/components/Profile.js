import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import BarGraph from "./BarGraph";
import PieChart from "./PieChart";

function Profile() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [data, setData] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      const token = localStorage.getItem("token");
      try {
        const response = await fetch(
          "https://01.gritlab.ax/api/graphql-engine/v1/graphql",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              query: `query {
                user {
                  login
                  id
                  transactions(
                    where: {
                      _or: [
                        { attrs: { _eq: {} } }
                        { attrs: { _has_key: "group" } }
                      ]
                      _and: [
                        { type: { _eq: "xp" } }
                        { path: { _nlike: "%/piscine-js/%" } }
                        { path: { _nlike: "%/piscine-go/%" } }
                      ]
                    }
                  ) {
                    attrs
                    path
                    type
                    amount
                    createdAt
                  }
                  upAmount: transactions_aggregate(where: {type: {_eq: "up"}}) {
                    aggregate {
                      sum {
                        amount
                      }
                    }
                  }
                  downAmount: transactions_aggregate(where: {type: {_eq: "down"}}) {
                    aggregate {
                      sum {
                        amount
                      }
                    }
                  }
                  xpAmount: transactions_aggregate(
                    where: {
                      type: { _eq: "xp" }
                      _or: [
                        { attrs: { _eq: {} } }
                        { attrs: { _has_key: "group" } }
                      ]
                      _and: [
                        { path: { _nlike: "%/piscine-js/%" } }
                        { path: { _nlike: "%/piscine-go/%" } }
                      ]
                    }
                  ) {
                    aggregate {
                      sum {
                        amount
                      }
                    }
                  }
                  auditRatio 
                  up: transactions_aggregate(where: {type: {_eq: "up"}}) {
                    aggregate {
                      sum {
                        amount
                      }
                    }
                  }
                  down: transactions_aggregate(where: {type: {_eq: "down"}}) {
                    aggregate {
                      sum {
                        amount
                      }
                    }
                  }
                }
              }`,
            }),
          }
        );

        const result = await response.json();
        setData(result.data);
        console.log(result.data.user[0].transactions);
        console.log(result.data);
        setLoading(false);
      } catch (error) {
        setError(error);
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const transactions = data.user[0].transactions;
  const sortedData = transactions.sort(
    (a, b) => new Date(a.createdAt) - new Date(b.createdAt)
  );

  const checkpointsReduced = sortedData.reduce((acc, item) => {
    const pathParts = item.path.split("/");
    const label = pathParts.includes("checkpoint")
      ? "Checkpoint"
      : pathParts.pop();
    const existingItem = acc.find((i) => i.label === label);
    if (existingItem) {
      existingItem.value += item.amount;
    } else {
      acc.push({ label, value: item.amount });
    }
    return acc;
  }, []);

  const auditRatio = parseFloat(data.user[0].auditRatio.toFixed(2));

  console.log(checkpointsReduced);

  return (
    <div className="profile-container">
      <h1>Profile</h1>
      <p>ID: {data.user[0].id}</p>
      <p>Login: {data.user[0].login}</p>
      <p>XP: {data.user[0].xpAmount.aggregate.sum.amount}</p>

      <h2>XP/project</h2>

      <BarGraph data={checkpointsReduced} width={1000} height={400} />

      <h2>Audit ratio</h2>
      <p>Audit Ratio: {auditRatio}</p>

      <PieChart data={data.user[0]} width={400} height={400} />
      <br />
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
}

export default Profile;
