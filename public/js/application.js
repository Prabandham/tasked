$(document).on("turbolinks:load", function() {
    // Hiding Links
    var path = window.location.pathname;
    
    if(path == "/") {
        $("#login-link").removeClass("hidden")
    } else {
        $("#login-link").addClass("hidden")
    }
});
