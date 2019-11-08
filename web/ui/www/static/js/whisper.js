$(document).ready(function(){
    $('[data-toggle="tooltip"]').tooltip();
});

var params = new URLSearchParams(window.location.search);

window.onload = function() {
    var action = window.location.pathname.replace("/", "");
    setupConsentForm(action);
    setupLoginPage(action);
    setupUpdatePage(action);
    setupRegistrationPage(action);
    setupEmailConfirmationPage(action);
    setupChangePasswordStep1Page(action);
    setupChangePasswordStep2Page(action);
};

function startSubmitting (obj) {
    obj.prop("disabled", true);

    obj.html(
        `<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>`
    );
}

function finishSubmitting (obj, text) {
    obj.prop("disabled", false);
    obj.html(text ? text : "Submit");
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

function isPasswordValid (password, username, email) {
    var minChar = $('#password-min-characters');
    var maxChar = $('#password-max-characters');
    var minUnique = $('#password-min-unique-characters');

    if (minChar == null || maxChar == null || minUnique == null) {
        return "Unable to load password policy";
    }

    minChar = parseInt(minChar.val());
    maxChar = parseInt(maxChar.val());
    minUnique = parseInt(minUnique.val());

    if (isNaN(minChar) || isNaN(maxChar) || isNaN(minUnique)) {
        return "Unable to load password policy";
    }

    if (!password || password.length < minChar) {
        return "Your password should have at least " + minChar + " characters";
    }

    if (!password || password.length > maxChar) {
        return "Your password should have at most " + maxChar + " characters";
    }

    var pass = password.toLowerCase();
    var user = username.toLowerCase();
    var mail = email.toLowerCase();

    if (pass.includes(user) || user.includes(pass)) {
        return "Your password is too similar to your username";
    }

    if (pass.includes(mail) || mail.includes(pass)) {
        return "Your password is too similar to your email";
    }

    var distinct = password.split('').filter(function (value, index, self) {
        return self.indexOf(value) === index;
    });

    if (!distinct || distinct.length < minUnique) {
        return "Your password should have at least " + minUnique + " unique characters";
    }

    return null;
}

function setupLoginPage(action) {
    if (action !== "login") {
        return;
    }

    var username = params.get("username");
    var firstLogin = params.get("first_login");

    if (username) {
        document.getElementById("login-username").value = username;
    }

    if (firstLogin) {
        notifySuccess("Whisper credential created successfully!")
    }

    $('#login-submit').on('click', function(event) {
        event.preventDefault();

        var $this = $(this);
        var request = {
            username: $("#login-username").val(),
            password: $("#login-password").val(),
            remember: $("#login-remember").is(":checked"),
            challenge: params.get("login_challenge")
        };

        if (!request.username) {
            notifyError("Username is missing");
            return;
        }

        if (!request.password) {
            notifyError("Password is missing");
            return;
        }

        if (!request.challenge) {
            notifyError("Challenge is missing");
            return;
        }

        startSubmitting($this);

        $.ajax({
            url: "/login",
            type: "POST",
            data: JSON.stringify(request),
            contentType: "application/json",
            success: function(data) {
                finishSubmitting($this);
                window.location = data.redirect_to;
            },
            error: function(xhr) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupConsentForm(action) {
    if (action !== "consent") {
        return;
    }

    function consent (answer) {
        return function (event) {
            event.preventDefault();

            var $this = $(this);
            var buttonText = answer ? "Allow" : "Deny";
            var request = {
                accept: answer,
                challenge: params.get("consent_challenge"),
                grantScope: $(".consent-grant-scope").toArray().map(function (item) { return item.value; }),
                remember: true
            };

            if (!request.challenge) {
                notifyError("Challenge is missing");
                return;
            }

            if (!request.grantScope) {
                notifyError("Grant Scopes are missing");
                return;
            }

            startSubmitting($this);

            $.ajax({
                url: "/consent",
                type: "POST",
                data: JSON.stringify(request),
                contentType: "application/json",
                success: function (data) {
                    finishSubmitting($this, buttonText);
                    window.location = data.redirect_to;
                },
                error: function (xhr) {
                    finishSubmitting($this, buttonText);
                    notifyError(xhr.responseText);
                }
            });
        }
    }

    $('#consent-allow').on('click', consent(true));
    $('#consent-deny').on('click', consent(false));
}

function setupUpdatePage(action) {
    if (action !== "secure/update") {
        return;
    }

    $('#update-submit').on('click', function(event) {
        event.preventDefault();

        var $this = $(this);
        var request = {
            email: $("#update-email").val(),
            newPassword: $("#update-new-password").val(),
            newPasswordConfirmation: $("#update-new-password-confirmation").val(),
            oldPassword: $("#update-old-password").val()
        };

        if (!request.email) {
            notifyError("Email is missing");
            return;
        }

        if (!request.newPassword) {
            notifyError("Invalid new password");
            return;
        }

        if (!request.oldPassword) {
            notifyError("Invalid old password");
            return;
        }

        if (request.oldPassword === request.newPassword) {
            notifyError("New password cannot be the same as the old");
            return;
        }

        if (request.newPassword !== request.newPasswordConfirmation) {
            notifyError("Invalid password confirmation");
            return;
        }

        var username = $("#update-username").val();
        var err = isPasswordValid(request.newPassword, username, request.email);

        if (err) {
            notifyError(err);
            return;
        }

        startSubmitting($this);

        $.ajax({
            url: "/secure/update",
            type: "PUT",
            data: JSON.stringify(request),
            contentType: "application/json",
            headers: {
                "Authorization": "Bearer " + params.get("token")
            },
            success: function() {
                finishSubmitting($this);
                window.location = params.get("redirect_to");
            },
            error: function(xhr) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupChangePasswordStep1Page(action) {
    if (action !== "change-password/step-1") {
        return;
    }

    var redirect_to = params.get("redirect_to");

    $('#submit').on('click', function (event) {
        event.preventDefault();

        var $this = $(this);
        var request = {
            redirect_to: redirect_to,
            email: $("#email").val(),
        };

        if (!request.email) {
            notifyError("Email is missing");
            return;
        }

        startSubmitting($this);

        $.ajax({
            url: "/change-password",
            type: "POST",
            data: JSON.stringify(request),
            contentType: "application/json",
            headers: {
                "Authorization": "Bearer " + params.get("token")
            },
            success: function() {
                finishSubmitting($this);
                notifySuccess("Check the inbox of your email.");
            },
            error: function(xhr) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupChangePasswordStep2Page(action) {
    if (action !== "change-password/step-2") {
        return;
    }

    var token = params.get("token");

    $('#submit').on('click', function (event) {
        event.preventDefault();

        var $this = $(this);
        var request = {
            token: token,
            newPassword: $("#new-password").val(),
            newPasswordConfirmation: $("#new-password-confirmation").val(),
        };

        if (!request.newPassword) {
            notifyError("Invalid new password");
            return;
        }

        if (request.newPassword !== request.newPasswordConfirmation) {
            notifyError("Invalid password confirmation");
            return;
        }

        var username = $("#username").val();
        var email = $("#email").val();
        var err = isPasswordValid(request.newPassword, username, email);

        if (err) {
            notifyError(err);
            return;
        }

        startSubmitting($this);

        $.ajax({
            url: "/change-password",
            type: "PUT",
            data: JSON.stringify(request),
            contentType: "application/json",
            success: function(data) {
                finishSubmitting($this);
                window.location = data.redirect_to;
            },
            error: function(xhr) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupRegistrationPage(action) {
    if (action !== "registration") {
        return;
    }

    $('#registration-submit').on('click', function(event) {
        event.preventDefault();

        var $this = $(this);
        var request = {
            username: $("#registration-username").val(),
            email: $("#registration-email").val(),
            password: $("#registration-password").val(),
            passwordConfirmation: $("#registration-password-confirmation").val(),
            challenge: params.get("login_challenge")
        };

        if (!request.username) {
            notifyError("Username is missing");
            return;
        }

        if (!request.email) {
            notifyError("Email is missing");
            return;
        }

        if (!request.password) {
            notifyError("Password is missing");
            return;
        }

        if (request.password !== request.passwordConfirmation) {
            notifyError("Invalid password confirmation");
            return;
        }

        var err = isPasswordValid(request.password, request.username, request.email);

        if (err) {
            notifyError(err);
            return;
        }

        startSubmitting($this);

        $.ajax({
            url: "/registration",
            type: "POST",
            data: JSON.stringify(request),
            contentType: "application/json",
            success: function() {
                finishSubmitting($(this));
                window.location = "/login?first_login=true&username="+$("#registration-username").val()+"&login_challenge="+$("#login-challenge").val();
            },
            error: function(xhr) {
                finishSubmitting($this);
                notifyError(xhr.responseText);
            }
        })
    })
}

function setupEmailConfirmationPage(action) {
    if (action !== "email-confirmation") {
        return;
    }

    var waitTime = 5000; // 5s
    var link = $("#redirect-to").val();
    var redirect = function () {
        window.location.href = link;
    };

    if (link) {
        setTimeout(redirect, waitTime)
    }
}