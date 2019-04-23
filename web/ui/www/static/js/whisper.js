window.onload = function() {
    var action = window.location.pathname.replace("/", "")
    setupConsentForm(action)
    setupLoginPage(action)
}

function setupLoginPage(action) {
    if (action == "login") {
        var params = new URLSearchParams(window.location.search)
        
        var firstLogin = params.get("first_login")
        if (firstLogin) {
            document.getElementById("credential-created-alert").hidden = false
            setTimeout(function(){
                document.getElementById("credential-created-alert").hidden = true
            }, 2000)
        }
        var username = params.get("username")
        document.getElementById("login-username").value = username
    }
}

function setupConsentForm(action) {
    if (action == "consent") {
        buttons = [{id: "consent-allow", value: "true"}, {id: "consent-deny", value: "false"}]
        for (i = 0; i < buttons.length; i++) {
            button = buttons[i]
            document.getElementById(button.id).addEventListener("click", function(value){
                return function(ev) {
                    event.preventDefault()
                    document.getElementById("accept-consent").value = value
                    document.getElementById("consent-form").submit();
                }
            }(button.value))
        }
    }
}