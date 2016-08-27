'use strict'

function run() {
    var container,
        menu,
        loggedIn = false,
        forms = {};

    menu = document.getElementById('menu');
    container = document.getElementById('container');

    function initLoginForm(form) {
    }

    function initUserForm(form) {
    }

    function initEventForm(form) {
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
