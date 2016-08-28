'use strict'

function run() {
    var container,
        http,
        menu,
        loggedIn = false,
        forms = {};

    menu = document.getElementById('menu');
    container = document.getElementById('container');
    http = new HTTP();

    /**
     * Defines the basic functionality for
     * the login form
     * @param {element} form - The DOM representation of the form
     */
    function initLoginForm(form) {
        function onLogin(res) {
            if (res.token) {
                http.setJWT(res.token);
                loggedIn = true;
                showForm('event');
            }
        }

        function onLoginFail(status, res) {
            // TODO(rdleon): show failure message
            console.log(status, res);
        }

        form.addEventListener('submit', function (event) {
            event.preventDefault();
            var usrInfo = {
                'email': '',
                'password': ''
                };

            usrInfo.email = form.querySelector('input[name="email"]').value;
            usrInfo.password = form.querySelector('input[name="password"]').value;

            http.post('/api/login', usrInfo, onLogin, onLoginFail);
        });
    }

    /**
     * Defines the basic functionality for the user form
     * for adding and modifying users.
     * @param {element} form - The DOM representation of the form
     */
    function initUserForm(form) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

    /**
     * Defines the basic functionality for the user form
     * for adding and modifying users.
     * @param {element} form - The DOM representation of the form
     */
    function initEventForm(form) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

    /**
     * Used to select which form to show at any given time in the UI
     * @param {string} formName - The name of the form to show
     */
    function showForm(formName) {
        // Toggle hidden in the selected form
        // hide all others.
        for (name in forms) {
           if (name == formName) {
               removeClass(forms[name], "hidden");
           } else {
               addClass(forms[name], "hidden")
           }
       }
    }

    /**
     * This is the function that bootstraps the forms and the
     * bulk of the web UI.
     */
    function init() {
        forms.login = document.getElementById('loginForm');
        forms.user = document.getElementById('userForm');
        forms.event = document.getElementById('eventForm');

        initLoginForm(forms.login);
        initUserForm(forms.user);
        initEventForm(forms.event);

        if (!loggedIn) {
            showForm('login');
        } else {
            showForm('event');
        }
    }

    init();
}

document.addEventListener('DOMContentLoaded', run);
