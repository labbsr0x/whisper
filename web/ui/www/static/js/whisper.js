window.onload = function() {
    var action = window.location.pathname.replace("/", "")
    setupConsentForm(action)
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