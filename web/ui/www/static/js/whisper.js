var params = new URLSearchParams(window.location.search);

window.onload = function() {
    var action = window.location.pathname.replace("/", "");
    setupConsentForm(action);
    setupLoginPage(action);
    setupUpdatePage(action);
    setupRegistrationPage(action);
}

function setupLoginPage(action) {
    if (action == "login") {
        
        var firstLogin = params.get("first_login");
        if (firstLogin) {
            document.getElementById("credential-created-alert").hidden = false;
            setTimeout(function(){
                document.getElementById("credential-created-alert").hidden = true
            }, 2000);
        }
        var username = params.get("username");
        document.getElementById("login-username").value = username;
    }
}

function setupConsentForm(action) {
    if (action == "consent") {
        buttons = [{id: "consent-allow", value: "true"}, {id: "consent-deny", value: "false"}];
        for (i = 0; i < buttons.length; i++) {
            button = buttons[i];
            document.getElementById(button.id).addEventListener("click", function(value){
                return function(ev) {
                    ev.preventDefault()
                    document.getElementById("accept-consent").value = value;
                    document.getElementById("consent-form").submit();
                }
            }(button.value));
        }
    }
}

function setupUpdatePage(action) {
    if (action == "secure/update") {
        
        $('#update-submit').on('click', function(event) {
            event.preventDefault()
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
                    window.location = params.get("redirect_to")
                },
                error: function(xhr, status, error) {
                    $("#credential-updated-alert")[0].innerHTML = xhr.responseText
                    $("#credential-updated-alert")[0].hidden = false
                    setTimeout(function(){
                        $("#credential-updated-alert")[0].hidden = true
                    }, 2000);
                }
            })
        })
    }
}

function setupRegistrationPage(action) {
    if (action == "registration") {
        $('#registration-submit').on('click', function(event) {
            event.preventDefault()

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
                    window.location = "/login?first_login=true&username="+$("#registration-username").val()+"&login_challenge="+$("#login-challenge").val()
                },
                error: function(xhr, status, error) {
                    $("#credential-registration-alert")[0].innerHTML = xhr.responseText
                    $("#credential-registration-alert")[0].hidden = false
                    setTimeout(function(){
                        $("#credential-registration-alert")[0].hidden = true
                    }, 2000);
                }
            })
        })
    }
}