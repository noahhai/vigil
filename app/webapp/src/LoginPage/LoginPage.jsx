import React from "react";
import {Formik, Field, Form, ErrorMessage} from "formik";
import * as Yup from "yup";
import "./LoginPage.css";

import {authenticationService} from "@/_services";

class LoginPage extends React.Component {
  constructor(props) {
    super(props);

    // redirect to home if already logged in
    if (authenticationService.currentUserValue) {
      this.props.history.push("/");
    }

    this.state = {
      signUp: false
    };
  }

  setSignUp(val) {
    this.setState({signUp: val});
  }

  render() {
    const {signUp} = this.state;
    if (!signUp) {
      return (
        <div>
          <div class="hero-image">
            <div class="hero-text">
              <h1 class="login">Vigil</h1>
            </div>
          </div>
          <br />
          <h3 class="login">Login</h3>
          <br />
          <Formik
            initialValues={{
              username: "",
              password: "",
              email: ""
            }}
            validationSchema={Yup.object().shape({
              username: Yup.string().required("Username/Email is required"),
              password: Yup.string().required("Password is required")
            })}
            onSubmit={({username, password}, {setStatus, setSubmitting}) => {
              setStatus();
              authenticationService
                .login(username, password)
                .then(
                  user => {
                    const {from} = this.props.location.state || {
                      from: {pathname: "/"}
                    };
                    this.props.history.push(from);
                  },
                  error => {
                    setSubmitting(false);
                    setStatus(error);
                  }
                )
                .catch(function(err) {
                  setSubmitting(false);
                  setStatus(err);
                });
            }}
            render={({errors, status, touched, isSubmitting}) => (
              <Form>
                <div className="form-group">
                  <label htmlFor="username">Username/Email</label>
                  <Field
                    name="username"
                    type="text"
                    className={
                      "form-control" +
                      (errors.username && touched.username ? " is-invalid" : "")
                    }
                  />
                  <ErrorMessage
                    name="username"
                    component="div"
                    className="invalid-feedback"
                  />
                </div>
                <br />
                <div className="form-group">
                  <label htmlFor="password">Password</label>
                  <Field
                    name="password"
                    type="password"
                    className={
                      "form-control" +
                      (errors.password && touched.password ? " is-invalid" : "")
                    }
                  />
                  <ErrorMessage
                    name="password"
                    component="div"
                    className="invalid-feedback"
                  />
                </div>
                <br />
                <div className="form-group">
                  <div style={{display: "flex", justifyContent: "flex-start"}}>
                    <button
                      type="submit"
                      className="btn btn-primary"
                      disabled={isSubmitting}
                    >
                      Login
                    </button>
                    <div style={{margin: "5px"}}></div>
                    <button
                      type="button"
                      onClick={() => this.setSignUp(true)}
                      className="btn btn-primary"
                      disabled={isSubmitting}
                    >
                      Register
                    </button>
                  </div>
                  {isSubmitting && (
                    <img src="data:image/gif;base64,R0lGODlhEAAQAPIAAP///wAAAMLCwkJCQgAAAGJiYoKCgpKSkiH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAEAAQAAADMwi63P4wyklrE2MIOggZnAdOmGYJRbExwroUmcG2LmDEwnHQLVsYOd2mBzkYDAdKa+dIAAAh+QQJCgAAACwAAAAAEAAQAAADNAi63P5OjCEgG4QMu7DmikRxQlFUYDEZIGBMRVsaqHwctXXf7WEYB4Ag1xjihkMZsiUkKhIAIfkECQoAAAAsAAAAABAAEAAAAzYIujIjK8pByJDMlFYvBoVjHA70GU7xSUJhmKtwHPAKzLO9HMaoKwJZ7Rf8AYPDDzKpZBqfvwQAIfkECQoAAAAsAAAAABAAEAAAAzMIumIlK8oyhpHsnFZfhYumCYUhDAQxRIdhHBGqRoKw0R8DYlJd8z0fMDgsGo/IpHI5TAAAIfkECQoAAAAsAAAAABAAEAAAAzIIunInK0rnZBTwGPNMgQwmdsNgXGJUlIWEuR5oWUIpz8pAEAMe6TwfwyYsGo/IpFKSAAAh+QQJCgAAACwAAAAAEAAQAAADMwi6IMKQORfjdOe82p4wGccc4CEuQradylesojEMBgsUc2G7sDX3lQGBMLAJibufbSlKAAAh+QQJCgAAACwAAAAAEAAQAAADMgi63P7wCRHZnFVdmgHu2nFwlWCI3WGc3TSWhUFGxTAUkGCbtgENBMJAEJsxgMLWzpEAACH5BAkKAAAALAAAAAAQABAAAAMyCLrc/jDKSatlQtScKdceCAjDII7HcQ4EMTCpyrCuUBjCYRgHVtqlAiB1YhiCnlsRkAAAOwAAAAAAAAAAAA==" />
                  )}
                </div>
                {status && <div className={"alert alert-danger"}>{status}</div>}
              </Form>
            )}
          />
        </div>
      );
    } else {
      return (
        <div>
          <div class="hero-image">
            <div class="hero-text">
              <h1 class="login">Vigil</h1>
            </div>
          </div>
          <br />
          <h3 class="login">Register</h3>
          <br />
          <Formik
            initialValues={{
              username: "",
              password: "",
              email: ""
            }}
            validationSchema={Yup.object().shape({
              username: Yup.string().required("Username is required"),
              password: Yup.string().required("Password is required")
            })}
            onSubmit={(
              {username, email, password},
              {setStatus, setSubmitting}
            ) => {
              setStatus();
              authenticationService
                .register(username, email, password)
                .then(
                  user => {
                    const {from} = this.props.location.state || {
                      from: {pathname: "/"}
                    };
                    this.props.history.push(from);
                  },
                  error => {
                    setSubmitting(false);
                    setStatus(error);
                  }
                )
                .catch(function(err) {
                  setSubmitting(false);
                  setStatus(err);
                });
            }}
            render={({errors, status, touched, isSubmitting}) => (
              <Form>
                <div className="form-group">
                  <label htmlFor="username">Username</label>
                  <Field
                    name="username"
                    type="text"
                    className={
                      "form-control" +
                      (errors.username && touched.username ? " is-invalid" : "")
                    }
                  />
                  <ErrorMessage
                    name="username"
                    component="div"
                    className="invalid-feedback"
                  />
                </div>
                <br />
                <div className="form-group">
                  <label htmlFor="email">Email (optional)</label>
                  <Field
                    name="email"
                    type="text"
                    className={
                      "form-control" +
                      (errors.email && touched.email ? " is-invalid" : "")
                    }
                  />
                  <ErrorMessage
                    name="email"
                    component="div"
                    className="invalid-feedback"
                  />
                </div>
                <br />
                <div className="form-group">
                  <label htmlFor="password">Password</label>
                  <Field
                    name="password"
                    type="password"
                    className={
                      "form-control" +
                      (errors.password && touched.password ? " is-invalid" : "")
                    }
                  />
                  <ErrorMessage
                    name="password"
                    component="div"
                    className="invalid-feedback"
                  />
                </div>
                <br />
                <div className="form-group">
                  <div style={{display: "flex", justifyContent: "flex-start"}}>
                    <button
                      type="submit"
                      className="btn btn-primary"
                      disabled={isSubmitting}
                    >
                      Register
                    </button>
                    <div style={{margin: "5px"}}></div>
                    <button
                      type="button"
                      onClick={() => this.setSignUp(false)}
                      className="btn btn-primary"
                      disabled={isSubmitting}
                    >
                      Cancel
                    </button>
                  </div>
                  {isSubmitting && (
                    <img src="data:image/gif;base64,R0lGODlhEAAQAPIAAP///wAAAMLCwkJCQgAAAGJiYoKCgpKSkiH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAEAAQAAADMwi63P4wyklrE2MIOggZnAdOmGYJRbExwroUmcG2LmDEwnHQLVsYOd2mBzkYDAdKa+dIAAAh+QQJCgAAACwAAAAAEAAQAAADNAi63P5OjCEgG4QMu7DmikRxQlFUYDEZIGBMRVsaqHwctXXf7WEYB4Ag1xjihkMZsiUkKhIAIfkECQoAAAAsAAAAABAAEAAAAzYIujIjK8pByJDMlFYvBoVjHA70GU7xSUJhmKtwHPAKzLO9HMaoKwJZ7Rf8AYPDDzKpZBqfvwQAIfkECQoAAAAsAAAAABAAEAAAAzMIumIlK8oyhpHsnFZfhYumCYUhDAQxRIdhHBGqRoKw0R8DYlJd8z0fMDgsGo/IpHI5TAAAIfkECQoAAAAsAAAAABAAEAAAAzIIunInK0rnZBTwGPNMgQwmdsNgXGJUlIWEuR5oWUIpz8pAEAMe6TwfwyYsGo/IpFKSAAAh+QQJCgAAACwAAAAAEAAQAAADMwi6IMKQORfjdOe82p4wGccc4CEuQradylesojEMBgsUc2G7sDX3lQGBMLAJibufbSlKAAAh+QQJCgAAACwAAAAAEAAQAAADMgi63P7wCRHZnFVdmgHu2nFwlWCI3WGc3TSWhUFGxTAUkGCbtgENBMJAEJsxgMLWzpEAACH5BAkKAAAALAAAAAAQABAAAAMyCLrc/jDKSatlQtScKdceCAjDII7HcQ4EMTCpyrCuUBjCYRgHVtqlAiB1YhiCnlsRkAAAOwAAAAAAAAAAAA==" />
                  )}
                </div>
                {status && <div className={"alert alert-danger"}>{status}</div>}
              </Form>
            )}
          />
        </div>
      );
    }
  }
}

export {LoginPage};
