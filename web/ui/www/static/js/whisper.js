window.onload = function() {
    var action = window.location.pathname.replace("/", "")
    if (action != "login" && action != "consent") {
        action = "error"
    }
    setupControl(action)
    setupConsentForm()
}

function setupConsentForm() {
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

function setupControl(action) {
    var params = new URLSearchParams(window.location.search)
    var control = getControl(action)
    control.element.hidden = false
    control.challengeElement.value = params.get(action+"_challenge")
}

function getControl(id) {
    return {
        matcher: new RegExp(id),
        element: document.getElementById(id+"-content"),
        challengeElement: document.getElementById(id+"-challenge")
    }
}