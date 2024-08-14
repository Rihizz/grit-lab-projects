class StateManager {
  constructor(initialState = {}) {
    const savedState = localStorage.getItem("appState");

    this.state = savedState ? JSON.parse(savedState) : initialState;

    localStorage.setItem("appState", JSON.stringify(this.state));
  }

  setState(newState) {
    this.state = { ...this.state, ...newState };

    localStorage.setItem("appState", JSON.stringify(this.state));
  }

  getState() {
    return this.state;
  }

  clearState() {
    this.state = {};

    localStorage.removeItem("appState");
  }
}

export default StateManager;
