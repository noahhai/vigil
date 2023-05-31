import React from "react";
import {Router, Route, Link} from "react-router-dom";

import {history} from "@/_helpers";
import {authenticationService} from "@/_services";
import {PrivateRoute} from "@/_components";
import {HomePage} from "@/HomePage";
import {LoginPage} from "@/LoginPage";
import "./App.css";

class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      currentUser: null
    };
  }

  componentDidMount() {
    authenticationService.currentUser.subscribe(x =>
      this.setState({
        currentUser: x
      })
    );
  }

  logout() {
    authenticationService.logout();
    history.push("/login");
  }

  render() {
    const {currentUser} = this.state;
    return (
      <Router history={history}>
        <div>
          {currentUser && (
            <nav className="navbar navbar-expand navbar-dark bg-dark">
              <div className="navbar-nav mr-auto">
                <a
                  className="link-color nav-item nav-link"
                  style={{cursor: "auto"}}
                >
                  User: {currentUser.username}
                </a>
                {/*<Link to="/" className="nav-item nav-link">Home</Link>*/}
              </div>
              <div className="navbar-nav ml-auto">
                <a
                  onClick={this.logout}
                  className="nav-item nav-link link-color"
                >
                  Logout
                </a>
                <a
                  id="about"
                  href="https://github.com/noahhai"
                  target="_blank"
                  className="nav-item nav-link link-color"
                >
                  About Author
                </a>
              </div>
            </nav>
          )}
          <div className="jumbotron">
            <div className="container">
              <div className="row">
                <div className="col-md-6 offset-md-3">
                  <PrivateRoute exact path="/" component={HomePage} />
                  <Route path="/login" component={LoginPage} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </Router>
    );
  }
}

export {App};
