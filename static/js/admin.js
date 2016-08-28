'use strict'

function run() {
    var http,
        menu,
        loggedIn = false,
        pages = {};

    /**
     * Defines the basic functionality for the login form
     * @param {element} form - The DOM representation of the form
     */
    function initLoginForm(form) {
        /**
         * Callback for login, gets the JWT and sets it in the http lib,
         * also sets the logged In flag, activates the menu and shows the
         * next step in the UI
         */
        function onLogin(res) {
            if (res.token) {
                http.setJWT(res.token);
                loggedIn = true;
                showPage('event');
            }
        }

        /**
         * Callback for login failure, notifies the user that the login
         * attempt was unsuccesful and gives the reason (credentials, server error)
         */
        function onLoginFail(status, res) {
            // TODO(rdleon): show failure message
            form.querySelector('input[name="password"]').value = '';
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
     * for adding and modifying events.
     * @param {element} form - The DOM representation of the form
     */
    function initEventForm(form) {
        var today = new Date();

        // Set defaults
        form.querySelector('input[name="duration"]').value = 120;
        form.querySelector('input[name="startTime"]').value = '18:00';
        form.querySelector('input[name="startDate"]').value = today.getFullYear() + '-01-01';

        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

    /**
     * Used to select which page to show at any given time in the UI
     * @param {string} pageName - The name of the page to show
     */
    function showPage(pageId) {
        var id;
        // Toggle hidden in the selected form
        // hide all others.
        for (id in pages) {
           if (id == pageId) {
               removeClass(pages[id], "hidden");
           } else if (pages[id]) {
               addClass(pages[id], "hidden")
           }
       }
    }

    /**
     * This is the function that bootstraps the forms and the
     * bulk of the web UI.
     */
    function init() {
        http = new HTTP();

        menu = document.getElementById('menu');

        pages.login = document.getElementById('loginForm');
        pages.user = document.getElementById('userForm');
        pages.event = document.getElementById('eventForm');
        pages.userList = document.getElementById('userList');

        initLoginForm(pages.login);
        initUserForm(pages.user);
        initEventForm(pages.event);

        if (!loggedIn) {
            showPage('login');
        } else {
            showPage('event');
        }
    }

    init();
}

document.addEventListener('DOMContentLoaded', run);
