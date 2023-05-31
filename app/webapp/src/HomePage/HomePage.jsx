import React from "react";
import {Checkbox} from "pretty-checkbox-react";
import Accordion from "./Accordion";
import "./HomePage.css";
import "./CheckboxStyle.scss";

import {userService, authenticationService} from "@/_services";
import {EditableContainer} from "@/Inputs/EditableContainer";

class HomePage extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      currentUser: authenticationService.currentUserValue,
      userFull: null,
      selectedIndex: 0
    };
  }

  componentDidMount() {
    const {username} = this.state.currentUser;
    userService.get(username).then(resp => {
      if (resp && resp.Data && resp.Data.length > 0) {
        this.setState({userFull: resp.Data[0]});
      }
    });
  }

  getUpdateUser(k) {
    return v => {
      const {userFull: u} = this.state;
      u[k] = v;
      userService.put(u);
      this.setState(state => ({userFull: u}));
    };
  }

  getHandleEventCheckboxUser(k) {
    const f = this.getUpdateUser(k);
    return e => {
      const val = e.nativeEvent.target.checked;
      return f(val);
    };
  }

  render() {
    const {currentUser, userFull} = this.state;
    return (
      <div>
        <h3>Alert Configuration</h3>

        {userFull && (
          <div>
            <EditableContainer
              label="Email"
              initialValue={userFull["email"]}
              onBlur={this.getUpdateUser("email")}
            ></EditableContainer>
            <EditableContainer
              label="Phone Number"
              initialValue={userFull["phone_number"]}
              onBlur={this.getUpdateUser("phone_number")}
            ></EditableContainer>

            <div className="flex-row">
              <div className="flex-items">SMS Notification</div>
              <div className="flex-items">
                <Checkbox
                  checked={userFull["notification_phone"]}
                  onChange={this.getHandleEventCheckboxUser(
                    "notification_phone"
                  )}
                  animation="smooth"
                ></Checkbox>
              </div>
            </div>

            <div className="flex-row">
              <div className="flex-items">Email Notification</div>
              <div className="flex-items">
                <Checkbox
                  checked={userFull["notification_email"]}
                  onChange={this.getHandleEventCheckboxUser(
                    "notification_email"
                  )}
                  animation="smooth"
                ></Checkbox>
              </div>
            </div>

            <div className="flex-row">
              <div className="flex-items">Agent Auth Token</div>
              <div className="flex-items">
                <div>{userFull["token"]} </div>
              </div>
            </div>
          </div>
        )}

        <br />
        <Accordion
          className="accordion"
          selectedIndex={this.state.selectedIndex}
        >
          <div data-header="About" className="accordion-item">
            <br />
            <p>
              Vigil is a simple command line utility that wraps other commands
              and notifies you when your command has finished running (via a sms
              or email). Example usage include long running scripts such as
              training ML models or data ETL.
            </p>
            <br />
          </div>
          <div data-header=" Installation And Usage" className="accordion-item">
            <div>
              <br />
              Create your configuration file at ~/.vigil/config.yaml and set:
            </div>
            <div style={{paddingLeft: "2em"}}>
              <div>username: &lt;USERNAME&gt;</div>
              <div>token: &lt;TOKEN&gt;</div>
            </div>
            <br />
            <p>
              Download the binary for your OS and add its location to your path.
              You may want to rename it to 'vigil'.
            </p>

            <div style={{paddingLeft: "2em"}}>
              <p>
                <a
                  target="_blank"
                  href="https://getvigil.io/agent/vigil-linux-x64"
                >
                  Linux
                </a>
              </p>
              <p>
                <a
                  target="_blank"
                  href="https://getvigil.io/agent/vigil-windows-x64.exe"
                >
                  Windows
                </a>
              </p>
              <p>
                <a
                  target="_blank"
                  href="https://getvigil.io/agent/vigil-darwin-x64"
                >
                  Mac
                </a>
              </p>
            </div>

            <p>
              To run, simply wrap your command with <em>vigil</em>. e.g.:{" "}
              <em>vigil echo "test"</em>
            </p>
          </div>
        </Accordion>
      </div>
    );
  }
}

export {HomePage};
