import React from "react";
import "../styles/EventContainer.css";

const EventContainer = ({ groupId, userID, eventsData }) => {
  const [events, setEvents] = React.useState(eventsData || []);
  const [eventResponses, setEventResponses] = React.useState({});
  const [userResponses, setUserResponses] = React.useState({});

  const fetchEventResponses = async () => {
    const responses = {};
    // console.log("print event length", events.length);
    for (let i = 0; i < events.length; i++) {
      // console.log("event id", events[i].id);
      const response = await fetch(
        "http://localhost:6969/api/serve-event-responses",
        {
          method: "POST",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ EventId: events[i].id }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        responses[events[i].id] = data;
        // console.log(`Fetched event response for eventId: ${events[i].id}`);
      } else {
        console.error(
          `Failed to fetch event response for eventId: ${events[i].id}`
        );
      }
    }

    setEventResponses(responses);
  };

  //print event responses
  // console.log(eventResponses);
  //console.log(userResponses);
  // console.log("printing event responses for event 1");
  // console.log(eventResponses[1]);
  // console.log("printing event responses for event 2");
  // console.log(eventResponses[2]);
  // console.log("printing event responses for event 3");
  // console.log(eventResponses[3]);
  //print user names for event 2
  // console.log("printing user names for event 2");
  if (eventResponses[2]) {
    for (let i = 0; i < eventResponses[2].length; i++) {
      // console.log(eventResponses[2][i].full_name);
    }
  }

  React.useEffect(() => {
    fetchEventResponses();
  }, [events]);

  const handleResponseChange = (eventId, e) => {
    setUserResponses({
      ...userResponses,
      [eventId]: e.target.value,
    });
  };

  const handleResponseSubmit = (eventId, e) => {
    e.preventDefault();
    postEventResponse(eventId, userResponses[eventId], userID);
  };

  const postEventResponse = async (eventId, response, userID) => {
    // console.log(eventId, response, userID);
    const responseObj = await fetch(
      "http://localhost:6969/api/event-response",
      {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          eventId: eventId,
          response: response,
          userID: userID,
        }),
      }
    );
    if (responseObj.ok) {
      // console.log("Event response posted");
      fetchEventResponses();
    } else {
      console.error("Failed to post event response");
    }
  };

  if (eventsData && eventsData !== events) {
    setEvents(eventsData);
  }
  //gigachad way of just deducting 3 hours from the time and displaying the shit out of it
  const formatDateTime = (dateTime) => {
    // Parse the date-time string
    const date = new Date(dateTime);

    // Subtract 3 hours in milliseconds
    date.setTime(date.getTime() - 3 * 60 * 60 * 1000);

    // Format the adjusted time
    const options = {
      weekday: "long",
      day: "numeric",
      month: "long",
      hour: "numeric",
      minute: "numeric",
    };
    return date.toLocaleString(undefined, options);
  };
  return (
    <div>
      {events.length === 0 ? (
        <p>No events to display, Create one!</p>
      ) : (
        events.map((event) => {
          const responses = eventResponses[event.id] || [];
          // console.log("printing responses in return lol");
          // console.log(responses);
          const attendees = responses.filter(
            (response) => response.response === "1"
          );
          // console.log("printing attendees");
          // console.log(attendees);
          const cowards = responses.filter(
            (response) => response.response === "0"
          );

          return (
            <div key={event.id}>
              <div className="event-container" key={event.id}>
                <b>
                  <h2>{event.title}</h2>
                </b>
                <h4>{formatDateTime(event.dateTime)}</h4>
                <p>{event.description}</p>
                <h4>Attendees: </h4>
                {attendees.map((attendee) => (
                  <span key={attendee.id}>{attendee.full_name}</span>
                ))}
                <h4>Cowards: </h4>
                {cowards.map((coward) => (
                  <span key={coward.id}>{coward.full_name}</span>
                ))}
                <form onSubmit={(e) => handleResponseSubmit(event.id, e)}>
                  <label>
                    Going:
                    <input
                      className="event-going-radio-button"
                      type="radio"
                      name={`response-${event.id}`}
                      value="going"
                      onChange={(e) => handleResponseChange(event.id, e)}
                    />
                  </label>
                  <label>
                    Not Going:
                    <input
                      className="event-going-radio-button"
                      type="radio"
                      name={`response-${event.id}`}
                      value="not going"
                      onChange={(e) => handleResponseChange(event.id, e)}
                    />
                  </label>
                  <button className="event-submitt-button" type="submit">
                    Submit
                  </button>
                </form>
              </div>
            </div>
          );
        })
      )}
    </div>
  );
};

export default EventContainer;
