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

    function initLoginForm(form) {
        function onLogin(res) {
            if (res.token) {
                http.setJWT(res.token);
                loggedIn = true;
                showForm('event');
            }
        }

        function onLoginFail(status, res) {
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

    function initUserForm(form) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

    function initEventForm(form) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

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

    init();
}

document.addEventListener('DOMContentLoaded', run);
