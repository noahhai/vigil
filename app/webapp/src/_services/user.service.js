import config from 'config';
import { authHeader, handleResponse } from '@/_helpers';

export const userService = {
    get,
    put,
};

function get(username) {
    const requestOptions = { method: 'GET', headers: authHeader() };
    return fetch(`${config.apiUrl}/users/${username}`, requestOptions).then(handleResponse);
}


function put(user) {
    const requestOptions = {
        method: 'PUT',
        headers: Object.assign(authHeader(), {'Content-Type': 'application/json'}),
        body: JSON.stringify(user)
    };

    return fetch(`${config.apiUrl}/users/`, requestOptions)
        .then(handleResponse);
}
