var params = new URLSearchParams(window.location.search);

window.onload = function() {
    var action = window.location.pathname.replace("/", "");
    setupConsentForm(action);
    setupLoginPage(action);
    setupUpdatePage(action);
    setupRegistrationPage(action);
};

function startSubmitting (obj) {
    obj.prop("disabled", true);

    obj.html(
        `<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>`
    );
}

function finishSubmitting (obj) {
    obj.prop("disabled", false);
    obj.html("Submit");
}

function notify (type, text) {
    var id = "#notification";
    var notificationTimeOut = 5000; // 5s

    $(id).html(text);
    $(id).attr("hidden", false);
    $(id).removeClass().addClass("alert alert-" + type);

    setTimeout(function(){
        $(id).attr("hidden", true)
    }, notificationTimeOut);
}

function notifyError (text) {
    notify("danger", text)
}

function notifySuccess (text) {
    notify("success", text)
}

function setupLoginPage(action) {
    if (action != "login") {
        return;
    }

    var username = params.get("username");
    var challenge = params.get("login_challenge");
    var firstLogin = params.get("first_login");

    if (firstLogin) {
        notifySuccess("Whisper credential created successfully!")
    }

    if (username) {
        document.getElementById("login-username").value = username;
    }

    $('#login-submit').on('click', function(event) {
        event.preventDefault();

        var $this = $(this);

        startSubmitting($this);

        $.ajax({
            url: "/login",
            type: "POST",
            data: JSON.stringify({
                username: $("#login-username").val(),
                password: $("#login-password").val(),
                remember: $("#login-remember").is(":checked"),
                challenge: challenge
            }),
            contentType: "application/json",
            success: function(data, status, xhr) {
                finishSubmitting($(this));
                window.location = data.redirect_to;
            },
            error: function(xhr, status, error) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupConsentForm(action) {
    if (action != "consent") {
        return;
    }

    buttons = [{id: "consent-allow", value: "true"}, {id: "consent-deny", value: "false"}];

    for (i = 0; i < buttons.length; i++) {
        button = buttons[i];
        document.getElementById(button.id).addEventListener("click", function(value){
            return function(ev) {
                ev.preventDefault();
                document.getElementById("accept-consent").value = value;
                document.getElementById("consent-form").submit();
            }
        }(button.value));
    }
}

function setupUpdatePage(action) {
    if (action != "secure/update") {
        return;
    }

    $('#update-submit').on('click', function(event) {
        event.preventDefault();

        var $this = $(this);

        startSubmitting($this);

        $.ajax({
            url: "/secure/update",
            type: "PUT",
            data: JSON.stringify({
                email: $("#update-email").val(),
                newPassword: $("#update-new-password").val(),
                newPasswordConfirmation: $("#update-new-password-confirmation").val(),
                oldPassword: $("#update-old-password").val()
            }),
            contentType: "application/json",
            headers: {
                "Authorization": "Bearer " + params.get("token")
            },
            success: function(data, status, xhr) {
                finishSubmitting($this);
                window.location = params.get("redirect_to");
            },
            error: function(xhr, status, error) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupRegistrationPage(action) {
    if (action != "registration") {
        return;
    }

    $('#registration-submit').on('click', function(event) {
        event.preventDefault();

        var $this = $(this);

        startSubmitting($this);

        $.ajax({
            url: "/registration",
            type: "POST",
            data: JSON.stringify({
                username: $("#registration-username").val(),
                email: $("#registration-email").val(),
                password: $("#registration-password").val(),
                passwordConfirmation: $("#registration-password-confirmation").val()
            }),
            contentType: "application/json",
            success: function(data, status, xhr) {
                finishSubmitting($(this));
                window.location = "/login?first_login=true&username="+$("#registration-username").val()+"&login_challenge="+$("#login-challenge").val();
            },
            error: function(xhr, status, error) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}